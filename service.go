package shogoa

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/dimfeld/httptreemux"
)

// Service is the data structure supporting shogoa services.
// It provides methods for configuring a service and running it.
// At the basic level a service consists of a set of controllers, each implementing a given
// resource actions. shogoagen generates global functions - one per resource - that make it
// possible to mount the corresponding controller onto a service. A service contains the
// middleware, not found handler, encoders and muxes shared by all its controllers.
type Service struct {
	// Name of service used for logging, tracing etc.
	Name string
	// Mux is the service request mux
	Mux ServeMux
	// Server is the service HTTP server.
	Server *http.Server
	// Context is the root context from which all request contexts are derived.
	// Set values in the root context prior to starting the server to make these values
	// available to all request handlers.
	Context context.Context
	// Request body decoder
	Decoder *HTTPDecoder
	// Response body encoder
	Encoder *HTTPEncoder

	middleware []Middleware       // Middleware chain
	cancel     context.CancelFunc // Service context cancel signal trigger
}

// Controller defines the common fields and behavior of generated controllers.
type Controller struct {
	// Controller resource name
	Name string

	// Service that exposes the controller
	Service *Service

	// BaseContext is a function that returns the base context for each request.
	BaseContext func(req *http.Request) context.Context

	// MaxRequestBodyLength is the maximum length read from request bodies.
	// Set to 0 to remove the limit altogether. Defaults to 1GB.
	MaxRequestBodyLength int64

	// FileSystem is used in FileHandler to open files. By default it returns
	// http.Dir but you can override it with another one that implements http.FileSystem.
	// For example using github.com/elazarl/go-bindata-assetfs is like below.
	//
	//	ctrl.FileSystem = func(dir string) http.FileSystem {
	//		return &assetfs.AssetFS{
	//			Asset: Asset,
	//			AssetDir: AssetDir,
	//			AssetInfo: AssetInfo,
	//			Prefix: dir,
	//		}
	//	}
	FileSystem func(string) http.FileSystem

	middleware []Middleware // Controller specific middleware if any
}

// FileServer is the interface implemented by controllers that can serve static files.
type FileServer interface {
	// FileHandler returns a handler that serves files under the given request path.
	FileHandler(path, filename string) Handler
}

// Handler defines the request handler signatures.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// Unmarshaler defines the request payload unmarshaler signatures.
type Unmarshaler func(context.Context, *Service, *http.Request) error

// DecodeFunc is the function that initialize the unmarshaled payload from the request body.
type DecodeFunc func(context.Context, io.ReadCloser, any) error

// New instantiates a service with the given name.
func New(name string) *Service {
	logHandler := slog.NewJSONHandler(os.Stderr, nil)
	ctx := WithLogger(context.Background(), NewLogger(logHandler))
	cctx, cancel := context.WithCancel(ctx)
	mux := NewMux()
	service := &Service{
		Name:    name,
		Context: cctx,
		Mux:     mux,
		Server: &http.Server{
			Handler: mux,
		},
		Decoder: NewHTTPDecoder(),
		Encoder: NewHTTPEncoder(),

		cancel: cancel,
	}
	service.Server.BaseContext = func(_ net.Listener) context.Context {
		return service.Context
	}
	var notFoundHandler Handler
	var methodNotAllowedHandler Handler

	// Setup default NotFound handler
	mux.HandleNotFound(func(rw http.ResponseWriter, req *http.Request, params url.Values) {
		if resp := ContextResponse(ctx); resp != nil && resp.Written() {
			return
		}
		// Use closure to do lazy computation of middleware chain so all middlewares are
		// registered.
		if notFoundHandler == nil {
			notFoundHandler = func(_ context.Context, _ http.ResponseWriter, req *http.Request) error {
				return ErrNotFound(req.URL.Path)
			}
			chain := service.middleware
			ml := len(chain)
			for i := range chain {
				notFoundHandler = chain[ml-i-1](notFoundHandler)
			}
		}
		ctx := NewContext(rw, req, params)
		err := notFoundHandler(ctx, ContextResponse(ctx), req)
		if !ContextResponse(ctx).Written() {
			if err := service.Send(ctx, 404, err); err != nil {
				service.LogError("failed to send 404 response", "error", err)
			}
		}
	})

	// Setup default MethodNotAllowed handler
	mux.HandleMethodNotAllowed(func(rw http.ResponseWriter, req *http.Request, params url.Values, methods map[string]httptreemux.HandlerFunc) {
		if resp := ContextResponse(ctx); resp != nil && resp.Written() {
			return
		}
		// Use closure to do lazy computation of middleware chain so all middlewares are
		// registered.
		if methodNotAllowedHandler == nil {
			methodNotAllowedHandler = func(_ context.Context, rw http.ResponseWriter, req *http.Request) error {
				allowedMethods := make([]string, len(methods))
				i := 0
				for k := range methods {
					allowedMethods[i] = k
					i++
				}
				slices.Sort(allowedMethods)
				rw.Header().Set("Allow", strings.Join(allowedMethods, ", "))
				return MethodNotAllowedError(req.Method, allowedMethods)
			}
			chain := service.middleware
			ml := len(chain)
			for i := range chain {
				methodNotAllowedHandler = chain[ml-i-1](methodNotAllowedHandler)
			}
		}
		ctx := NewContext(rw, req, params)
		err := methodNotAllowedHandler(ctx, ContextResponse(ctx), req)
		if !ContextResponse(ctx).Written() {
			if err := service.Send(ctx, 405, err); err != nil {
				service.LogError("failed to send 405 response", "error", err)
			}
		}
	})

	return service
}

// CancelAll sends a cancel signals to all request handlers via the context.
// See https://golang.org/pkg/context/ for details on how to handle the signal.
func (service *Service) CancelAll() {
	service.cancel()
}

// Use adds a middleware to the service wide middleware chain.
// shogoa comes with a set of commonly used middleware, see the middleware package.
// Controller specific middleware should be mounted using the Controller struct Use method instead.
func (service *Service) Use(m Middleware) {
	service.middleware = append(service.middleware, m)
}

// WithLogger sets the logger used internally by the service and by Log.
func (service *Service) WithLogger(logger LogAdapter) {
	service.Context = WithLogger(service.Context, logger)
}

// LogInfo logs the message and values at odd indexes using the keys at even indexes of the keyvals slice.
func (service *Service) LogInfo(msg string, keyvals ...any) {
	ctx := service.Context

	// this block should be synced with LogInfo.
	// we want not to call LogInfo because it changes the call stack
	// and makes the log adapter more complex to implement.
	if l := ctx.Value(logKey); l != nil {
		switch logger := l.(type) {
		case LogAdapter:
			logger.InfoContext(ctx, msg, keyvals...)
		}
	}
}

// LogError logs the error and values at odd indexes using the keys at even indexes of the keyvals slice.
func (service *Service) LogError(msg string, keyvals ...any) {
	ctx := service.Context

	// this block should be synced with LogError.
	// we want not to call LogError because it changes the call stack
	// and makes the log adapter more complex to implement.
	if l := ctx.Value(logKey); l != nil {
		switch logger := l.(type) {
		case LogAdapter:
			logger.ErrorContext(ctx, msg, keyvals...)
		}
	}
}

// ListenAndServe starts a HTTP server and sets up a listener on the given host/port.
func (service *Service) ListenAndServe(addr string) error {
	service.LogInfo("listen", "transport", "http", "addr", addr)
	service.Server.Addr = addr
	return service.Server.ListenAndServe()
}

// ListenAndServeTLS starts a HTTPS server and sets up a listener on the given host/port.
func (service *Service) ListenAndServeTLS(addr, certFile, keyFile string) error {
	service.LogInfo("listen", "transport", "https", "addr", addr)
	service.Server.Addr = addr
	return service.Server.ListenAndServeTLS(certFile, keyFile)
}

// Serve accepts incoming HTTP connections on the listener l, invoking the service mux handler for each.
func (service *Service) Serve(l net.Listener) error {
	return service.Server.Serve(l)
}

// NewController returns a controller for the given resource. This method is mainly intended for
// use by the generated code. User code shouldn't have to call it directly.
func (service *Service) NewController(name string) *Controller {
	return &Controller{
		Name:    name,
		Service: service,
		BaseContext: func(req *http.Request) context.Context {
			ctx := req.Context()
			ctx = context.WithValue(ctx, ctrlKey, name)
			return ctx
		},
		MaxRequestBodyLength: 1073741824, // 1 GB
		FileSystem: func(dir string) http.FileSystem {
			return http.Dir(dir)
		},
	}
}

// Send serializes the given body matching the request Accept header against the service
// encoders. It uses the default service encoder if no match is found.
func (service *Service) Send(ctx context.Context, code int, body any) error {
	r := ContextResponse(ctx)
	if r == nil {
		return fmt.Errorf("no response data in context")
	}
	r.WriteHeader(code)
	return service.EncodeResponse(ctx, body)
}

// ServeFiles create a "FileServer" controller and calls ServerFiles on it.
func (service *Service) ServeFiles(path, filename string) error {
	ctrl := service.NewController("FileServer")
	return ctrl.ServeFiles(path, filename)
}

// DecodeRequest uses the HTTP decoder to unmarshal the request body into the provided value based
// on the request Content-Type header.
func (service *Service) DecodeRequest(req *http.Request, v any) error {
	body, contentType := req.Body, req.Header.Get("Content-Type")
	defer body.Close()

	if err := service.Decoder.Decode(v, body, contentType); err != nil {
		return fmt.Errorf("failed to decode request body with content type %#v: %s", contentType, err)
	}

	return nil
}

// EncodeResponse uses the HTTP encoder to marshal and write the response body based on the request
// Accept header.
func (service *Service) EncodeResponse(ctx context.Context, v any) error {
	accept := ContextRequest(ctx).Header.Get("Accept")
	return service.Encoder.Encode(v, ContextResponse(ctx), accept)
}

// ServeFiles replies to the request with the contents of the named file or directory. See
// FileHandler for details.
func (ctrl *Controller) ServeFiles(path, filename string) error {
	ctx := ctrl.Service.Context
	if strings.Contains(path, ":") {
		return fmt.Errorf("path may only include wildcards that match the entire end of the URL (e.g. *filepath)")
	}
	LogInfo(ctx, "mount file", "name", filename, "route", fmt.Sprintf("GET %s", path))
	handler := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		if !ContextResponse(ctx).Written() {
			return ctrl.FileHandler(path, filename)(ctx, rw, req)
		}
		return nil
	}
	ctrl.Service.Mux.Handle("GET", path, ctrl.MuxHandler("serve", handler, nil))
	return nil
}

// Use adds a middleware to the controller.
// Service-wide middleware should be added via the Service Use method instead.
func (ctrl *Controller) Use(m Middleware) {
	ctrl.middleware = append(ctrl.middleware, m)
}

// MuxHandler wraps a request handler into a MuxHandler. The MuxHandler initializes the request
// context by loading the request state, invokes the handler and in case of error invokes the
// controller (if there is one) or Service error handler.
// This function is intended for the controller generated code. User code should not need to call
// it directly.
func (ctrl *Controller) MuxHandler(name string, hdlr Handler, unm Unmarshaler) MuxHandler {
	// Use closure to enable late computation of handlers to ensure all middleware has been
	// registered.
	var handler Handler
	var initHandler sync.Once

	return func(rw http.ResponseWriter, req *http.Request, params url.Values) {
		// Build handler middleware chains on first invocation
		initHandler.Do(func() {
			handler = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
				if !ContextResponse(ctx).Written() {
					return hdlr(ctx, rw, req)
				}
				return nil
			}
			mwLen := len(ctrl.Service.middleware)
			chain := append(ctrl.Service.middleware[:mwLen:mwLen], ctrl.middleware...)
			ml := len(chain)
			for i := range chain {
				handler = chain[ml-i-1](handler)
			}
		})

		// Build context
		ctx := NewContext(rw, req, params)
		req = req.WithContext(ctx)
		ctx = ctrl.BaseContext(req)
		ctx = WithAction(ctx, name)

		// Protect against request bodies with unreasonable length
		if ctrl.MaxRequestBodyLength > 0 {
			req.Body = http.MaxBytesReader(rw, req.Body, ctrl.MaxRequestBodyLength)
		}

		// Load body if any
		if req.ContentLength > 0 && unm != nil {
			if err := unm(ctx, ctrl.Service, req); err != nil {
				if err.Error() == "http: request body too large" {
					msg := fmt.Sprintf("request body length exceeds %d bytes", ctrl.MaxRequestBodyLength)
					err = ErrRequestBodyTooLarge(msg)
				} else {
					err = ErrBadRequest(err)
				}
				ctx = WithError(ctx, err)
			}
		}

		// Invoke handler
		if err := handler(ctx, ContextResponse(ctx), req); err != nil {
			LogError(ctx, "uncaught error", "err", err)
			respBody := fmt.Sprintf("Internal error: %s", err) // Sprintf catches panics
			_ = ctrl.Service.Send(ctx, 500, respBody)
		}
	}
}

// FileHandler returns a handler that serves files under the given filename for the given route path.
// The logic for what to do when the filename points to a file vs. a directory is the same as the
// standard http package ServeFile function. The path may end with a wildcard that matches the rest
// of the URL (e.g. *filepath). If it does the matching path is appended to filename to form the
// full file path, so:
//
//	c.FileHandler("/index.html", "/www/data/index.html")
//
// Returns the content of the file "/www/data/index.html" when requests are sent to "/index.html"
// and:
//
//	c.FileHandler("/assets/*filepath", "/www/data/assets")
//
// returns the content of the file "/www/data/assets/x/y/z" when requests are sent to
// "/assets/x/y/z".
func (ctrl *Controller) FileHandler(path, filename string) Handler {
	var wc string
	if idx := strings.LastIndex(path, "/*"); idx > -1 && idx < len(path)-1 {
		wc = path[idx+2:]
		if strings.Contains(wc, "/") {
			wc = ""
		}
	}
	return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// prevent path traversal
		if attemptsPathTraversal(req.URL.Path, path) {
			return ErrNotFound(req.URL.Path)
		}
		fname := filename
		if len(wc) > 0 {
			if m, ok := ContextRequest(ctx).Params[wc]; ok {
				fname = filepath.Join(filename, m[0])
			}
		}
		LogInfo(ctx, "serve file", "name", fname, "route", req.URL.Path)
		dir, name := filepath.Split(fname)
		fs := ctrl.FileSystem(dir)
		f, err := fs.Open(name)
		if err != nil {
			return ErrInvalidFile(err)
		}
		defer f.Close()
		d, err := f.Stat()
		if err != nil {
			return ErrInvalidFile(err)
		}
		// use contents of index.html for directory, if present
		if d.IsDir() {
			index := strings.TrimSuffix(name, "/") + "/index.html"
			ff, err := fs.Open(index)
			if err == nil {
				defer ff.Close()
				dd, err := ff.Stat()
				if err == nil {
					name = index
					d = dd
					f = ff
				}
			}
		}
		_ = name

		// serveContent will check modification time
		// Still a directory? (we didn't find an index.html file)
		if d.IsDir() {
			return dirList(rw, f)
		}
		http.ServeContent(rw, req, d.Name(), d.ModTime(), f)
		return nil
	}
}

func attemptsPathTraversal(req string, path string) bool {
	if !strings.Contains(req, "..") {
		return false
	}

	currentPathIdx := 0
	if idx := strings.LastIndex(path, "/*"); idx > -1 && idx < len(path)-1 {
		req = req[idx+1:]
	}
	for _, runeValue := range strings.FieldsFunc(req, isSlashRune) {
		if runeValue == ".." {
			currentPathIdx--
			if currentPathIdx < 0 {
				return true
			}
		} else {
			currentPathIdx++
		}
	}
	return false
}

func isSlashRune(r rune) bool {
	return os.IsPathSeparator(uint8(r))
}

var replacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)

func dirList(w http.ResponseWriter, f http.File) error {
	dirs, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Sort(byName(dirs))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", url.String(), replacer.Replace(name))
	}
	fmt.Fprintf(w, "</pre>\n")
	return nil
}

type byName []os.FileInfo

func (s byName) Len() int           { return len(s) }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
