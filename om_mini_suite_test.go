package main_test

import (
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOmMini(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "om-mini Suite")
}

var _ = Describe("the binary", func() {
	It("compiles successfully", func() {
		path, err := gexec.Build("github.com/jtarchie/om-mini")
		Expect(err).NotTo(HaveOccurred())

		command := exec.Command(path, "--help")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
	})
})
