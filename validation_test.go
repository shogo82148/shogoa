package shogoa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shogo82148/shogoa"
)

var _ = Describe("ValidateFormat", func() {
	var f shogoa.Format
	var val string
	var valErr error

	BeforeEach(func() {
		val = ""
	})

	JustBeforeEach(func() {
		valErr = shogoa.ValidateFormat(f, val)
	})

	Context("Date", func() {
		BeforeEach(func() {
			f = shogoa.FormatDate
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "201510-26"
			})

			It("does not validates", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "2015-10-26"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("DateTime", func() {
		BeforeEach(func() {
			f = shogoa.FormatDateTime
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "201510-26T08:31:23Z"
			})

			It("does not validates", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "2015-10-26T08:31:23Z"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("UUID", func() {
		BeforeEach(func() {
			f = shogoa.FormatUUID
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "96054a62-a9e45ed26688389b"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with an valid value", func() {
			BeforeEach(func() {
				val = "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("Email", func() {
		BeforeEach(func() {
			f = shogoa.FormatEmail
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "foo"
			})

			It("does not validates", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "raphael@shogoa.design"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("Hostname", func() {
		BeforeEach(func() {
			f = shogoa.FormatHostname
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "_hi_"
			})

			It("does not validates", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "shogoa.design"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("IPv4", func() {
		BeforeEach(func() {
			f = shogoa.FormatIPv4
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "192-168.0.1"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid IPv6 value", func() {
			BeforeEach(func() {
				val = "::1"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "192.168.0.1"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("IPv6", func() {
		BeforeEach(func() {
			f = shogoa.FormatIPv6
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "foo"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid IPv4 value", func() {
			BeforeEach(func() {
				val = "10.10.10.10"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "0:0:0:0:0:0:0:1"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("IP", func() {
		BeforeEach(func() {
			f = shogoa.FormatIP
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "::1.1"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid IPv4 value", func() {
			BeforeEach(func() {
				val = "127.0.0.1"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

		Context("with a valid IPv6 value", func() {
			BeforeEach(func() {
				val = "::1"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("URI", func() {
		BeforeEach(func() {
			f = shogoa.FormatURI
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "foo_"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "hhp://shogoa.design/contact"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("MAC", func() {
		BeforeEach(func() {
			f = shogoa.FormatMAC
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "bar"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "06-00-00-00-00-00"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("CIDR", func() {
		BeforeEach(func() {
			f = shogoa.FormatCIDR
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "foo"
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "10.0.0.0/8"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("Regexp", func() {
		BeforeEach(func() {
			f = shogoa.FormatRegexp
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "foo["
			})

			It("does not validate", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "^shogoa$"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})

	})

	Context("RFC1123", func() {
		BeforeEach(func() {
			f = shogoa.FormatRFC1123
		})

		Context("with an invalid value", func() {
			BeforeEach(func() {
				val = "Mon 04 Jun 2017 23:52:05 MST"

			})

			It("does not validates", func() {
				Ω(valErr).Should(HaveOccurred())
			})
		})

		Context("with a valid value", func() {
			BeforeEach(func() {
				val = "Mon, 04 Jun 2017 23:52:05 MST"
			})

			It("validates", func() {
				Ω(valErr).ShouldNot(HaveOccurred())
			})
		})
	})
})
