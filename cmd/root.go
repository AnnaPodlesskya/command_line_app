package cmd

import (
	"command_line_programs/pkg/repo_manager"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"strings"
)

var configFilename string
var ignoreErrors bool
var rootCmd = &cobra.Command{
	Use:   "multi-git",
	Short: "Runs git commands over multiple repos",
	Long: `Runs git commands over multiple repos.

	Requires the following environment variables defined:
	MG_ROOT: root directory of target git repositories
	MG_REPOS: list of repository names to operate on`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := viper.GetString("root")
		if root[len(root)-1] != '/' {
			root += "/"
		}
		repoNames := []string{}
		if len(viper.GetString("repos")) > 0 {
			repoNames = strings.Split(viper.GetString("repos"), ",")
		}

		repoManger, err := repo_manager.NewRepoManager(root, repoNames, viper.GetBool("ignore-errors"))
		if err != nil {
			log.Fatal(err)
		}
		command := strings.Join(args, " ")
		output, err := repoManger.Exec(command)
		if err != nil {
			fmt.Fprintf(os.Stderr, ">>> Error returned from Exec: %v\n", err)
		}
		for repo, out := range output {
			fmt.Printf("[%s]: git %s\n", path.Base(repo), command)
			fmt.Println(out)
		}
		fmt.Println("Done.")
	},
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	home, err := homedir.Dir()
	check(err)

	defaultConfigFilename := path.Join(home, ".config/multi-git.toml")
	rootCmd.Flags().StringVar(&configFilename,
		"config",
		defaultConfigFilename,
		"config file path (default is $HOME/multi-git.toml)")
	rootCmd.Flags().Bool(
		"ignore-errors",
		false,
		`will continue executing the command for all repos if ignore-errors
                 is true otherwise it will stop execution when an error occurs`)
	err = viper.BindPFlag("ignore-errors", rootCmd.Flags().Lookup("ignore-errors"))
	if err != nil {
		panic("Unable to bind flag")
	}
}
func initConfig() {
	_, err := os.Stat(configFilename)
	if os.IsNotExist(err) {
		check(err)
	}
	viper.SetConfigFile(configFilename)
	err = viper.ReadInConfig()
	check(err)

	viper.SetEnvPrefix("MG")
	err = viper.BindEnv("root")
	check(err)

	err = viper.BindEnv("repos")
	check(err)
}
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
