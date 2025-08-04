package repo_manager

import (
	. "command_line_programs/pkg/helpers"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const baseDir = "/tmp/test-multi-git"

var repoList = []string{}

var _ = Describe("Repo manager tests", func() {
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
		err = CreateDir(baseDir, "dir-1", true)
		Ω(err).Should(BeNil())
		repoList = []string{"dir-1"}
	})
	AfterEach(removeAll)

	It("Should fail with invalid base dir", func() {
		_, err := NewRepoManager("/no-such-dir", repoList, true)
		Ω(err).ShouldNot(BeNil())
	})
	It("should fail with empty repo list", func() {
		_, err := NewRepoManager(baseDir, []string{}, true)
		Ω(err).ShouldNot(BeNil())
	})
	It("Should commit files successfully", func() {
		rm, err := NewRepoManager(baseDir, repoList, true)
		Ω(err).Should(BeNil())

		output, err := rm.Exec("checkout -b test-branch")
		Ω(err).Should(BeNil())

		for _, out := range output {
			Ω(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
		}
		AddFiles(baseDir, repoList[0], true, "file_1.txt", "file_2.txt")

		wd, _ := os.Getwd()
		defer os.Chdir(wd)

		dir := path.Join(absBaseDir, repoList[0])
		err = os.Chdir(dir)
		Ω(err).Should(BeNil())

		output, err = rm.Exec("log --oneline")
		Ω(err).Should(BeNil())

		ok := strings.HasSuffix(output[dir], "added some files...\n")
		Ω(ok).Should(BeTrue())
	})
	It("Should create brunch successfully", func() {
		repoList = append(repoList, "dir-2")
		CreateDir(baseDir, repoList[1], true)
		rm, err := NewRepoManager(baseDir, repoList, true)
		Ω(err).Should(BeNil())

		output, err := rm.Exec("checkout -b test-branch")
		Ω(err).Should(BeNil())

		for _, out := range output {
			Ω(out).Should(Equal("Switched to a new branch 'test-branch'\n"))
		}
	})
	It("Should get repo list successfully with non-git directories", func() {
		repoList = append(repoList, "dir-2")
		repoList = append(repoList, "not-a-git-repo")
		CreateDir(baseDir, repoList[1], true)
		CreateDir(baseDir, repoList[2], false)

		rm, err := NewRepoManager(baseDir, repoList, true)
		Ω(err).Should(BeNil())

		repos := rm.GetRepos()
		Ω(repos).Should(HaveLen(3))

		Ω(repos).Should(ConsistOf(
			path.Join(absBaseDir, "dir-1"),
			path.Join(absBaseDir, "dir-2"),
			path.Join(absBaseDir, "not-a-git-repo"),
		))

		outputs, err := rm.Exec("status")
		Ω(err).Should(BeNil())

		Ω(outputs[path.Join(absBaseDir, "dir-1")]).Should(ContainSubstring("On branch"))
		Ω(outputs[path.Join(absBaseDir, "dir-2")]).Should(ContainSubstring("On branch"))
		fmt.Println("OUTPUT for not-a-git-repo:")
		fmt.Println(outputs[path.Join(absBaseDir, "not-a-git-repo")])
		Ω(outputs[path.Join(absBaseDir, "not-a-git-repo")]).Should(ContainSubstring("fatal: not a git repository"))
	})
})
