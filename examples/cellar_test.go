/*
Package example_test validates that the code generated by goagen by the "main" and "app"
generators from the cellar example design package produce valid Go code that compiles and runs.
*/
package examples_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("example cellar", func() {
	var tempdir string

	var files = []string{
		filepath.FromSlash("app"),
		filepath.FromSlash("app/contexts.go"),
		filepath.FromSlash("app/controllers.go"),
		filepath.FromSlash("app/hrefs.go"),
		filepath.FromSlash("app/media_types.go"),
		filepath.FromSlash("app/user_types.go"),
		filepath.FromSlash("main.go"),
		filepath.FromSlash("account.go"),
		filepath.FromSlash("bottle.go"),
		filepath.FromSlash("client"),
		filepath.FromSlash("client/cellar-cli"),
		filepath.FromSlash("client/cellar-cli/main.go"),
		filepath.FromSlash("client/cellar-cli/commands.go"),
		filepath.FromSlash("client/client.go"),
		filepath.FromSlash("client/account.go"),
		filepath.FromSlash("client/bottle.go"),
		filepath.FromSlash("swagger"),
		filepath.FromSlash("swagger/swagger.json"),
		filepath.FromSlash("swagger/swagger.go"),
		"",
	}

	BeforeEach(func() {
		var err error
		gopath := filepath.SplitList(os.Getenv("GOPATH"))[0]
		tempdir, err = ioutil.TempDir(filepath.Join(gopath, "src"), "cellar-test-tmpdir-")
		Ω(err).ShouldNot(HaveOccurred())
	})

	JustBeforeEach(func() {
		cmd := exec.Command("goagen", "bootstrap", "-d", "github.com/raphael/goa/examples/cellar/design")
		cmd.Dir = tempdir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("\n==========%s\n==========\n", out)
		}
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(out)).Should(Equal(strings.Join(files, "\n")))
	})

	It("goagen generated valid Go code", func() {
		bin := "cellar"
		if runtime.GOOS == "windows" {
			bin += ".exe"
		}
		cmd := exec.Command("go", "build", "-o", bin)
		cmd.Dir = tempdir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("\n==========%s\n==========\n", out)
		}
		Ω(err).ShouldNot(HaveOccurred())
		cmd = exec.Command(fmt.Sprintf(".%c%s", filepath.Separator, bin))
		cmd.Dir = tempdir
		b := &bytes.Buffer{}
		cmd.Stdout = b
		err = cmd.Start()
		Ω(err).ShouldNot(HaveOccurred())
		done := make(chan error, 1)
		go func() {
			done <- cmd.Wait()
		}()
		select {
		case <-time.After(200 * time.Millisecond):
			cmd.Process.Kill()
			<-done
		case err := <-done:
			Ω(err).ShouldNot(HaveOccurred())
		}
		Ω(err).ShouldNot(HaveOccurred())
		Ω(b.String()).Should(ContainSubstring("file=swagger/swagger.json"))
	})

	AfterEach(func() {
		os.RemoveAll(tempdir)
	})
})
