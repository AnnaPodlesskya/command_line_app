package main

import (
	"command_line_programs/pkg/repo_manager"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	command := flag.String("command", "", "The git command")
	ignoreErros := flag.Bool(
		"ignore-errors",
		false, "Keep running after error if true")
	flag.Parse()

	root := os.Getenv("MG_ROOT")
	if root[len(root)-1] != '/' {
		root += "/"
	}
	repoNames := strings.Split(os.Getenv("MG_REPOS"), ",")

	repoManger, err := repo_manager.NewRepoManager(root, repoNames, *ignoreErros)
	if err != nil {
		log.Fatal(err)
	}
	output, _ := repoManger.Exec(*command)
	for repo, out := range output {
		fmt.Printf("[%s]: git %s\n", path.Base(repo), *command)
		fmt.Println(out)
	}
	fmt.Println("Done.")
}
