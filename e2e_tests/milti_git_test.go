package e2e_tests

import (
	. "command_line_programs/pkg/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"strings"
)

const baseDir = "/tmp/test-multi-git"

var repoList string

var _ = Describe("multi-git e2e tests", func() {
	var err error
	var absBaseDir string

	removeAll := func() {
		err = os.RemoveAll(absBaseDir)
		Ω(err).Should(BeNil())
	}
	BeforeEach(func() {
		absBaseDir, err = filepath.Abs(baseDir)
		Ω(err).Should(BeNil())
		removeAll()
		err = CreateDir(baseDir, "", false)
		Ω(err).Should(BeNil())

	})
	AfterSuite(removeAll)

	Context("Tests for empty/undefined environment failure cases", func() {
		It("Should fail with invalid base dir", func() {
			output, err := RunMultiGit("status", false, "/no-such-dit", repoList)
			Ω(err).ShouldNot(BeNil())
			suffix := "base dir: '/no-such-dir/' doesnt exist \n"
			Ω(output).Should(HaveSuffix(suffix))
		})

		It("Should fail with empty repo list", func() {
			output, err := RunMultiGit("status", false, absBaseDir, repoList)
			Ω(err).ShouldNot(BeNil())
			Ω(output).Should(ContainSubstring("repo list cant be empty"))
		})
	})
	Context("Tests for success cases", func() {
		It("Should do git init successfully", func() {
			err = CreateDir(baseDir, "dir-1", false)
			Ω(err).Should(BeNil())
			err = CreateDir(baseDir, "dir-2", false)
			Ω(err).Should(BeNil())
			repoList = "dir-1,dir-2"

			output, err := RunMultiGit("init", false, absBaseDir, repoList)
			Ω(err).Should(BeNil())
			count := strings.Count(output, "Initialized empty Git repository")
			Ω(count).Should(Equal(2))

		})
	})
})
