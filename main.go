package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/bluemir/sentry/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var persistentConfig struct {
	printVersion bool
	verbose      bool
	cfgFile      string
}

var config *core.Config

func _main(cmd *cobra.Command, args []string) {
	if persistentConfig.printVersion {
		fmt.Printf("sentry version %s-%d\n", __VERSION__, __BUILD_NUM__)
		return
	}
	if persistentConfig.verbose {
		log.SetLevel(log.DebugLevel)
	}

	core.NewSentry(config).Run()
}

var RootCmd = &cobra.Command{
	Use:   "sentry",
	Short: "file system watcher and running command",
	Long: `Sentry is command line tool that execute command when file is changed.
Sentry watch file system event, and everytime you change/create on file it will re-execute command

github : http://github.com/bluemir/sentry
`,
	Run: _main,
}

func init() {
	config = &core.Config{}

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&persistentConfig.cfgFile, "config", "", "config file (default is $HOME/.tick.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&persistentConfig.verbose, "verbose", "v", false, "show detail log")
	RootCmd.PersistentFlags().BoolVar(&persistentConfig.printVersion, "version", false, "show version")

	RootCmd.Flags().StringVarP(&config.Command, "command", "c", "", "command to run")
	RootCmd.Flags().Int32VarP(&config.Delay, "delay", "d", 500, "delay that wait events")
	RootCmd.Flags().StringSliceVarP(&config.WatchPaths, "watch", "w", []string{"./"}, "paths to watch")
	RootCmd.Flags().StringSliceVarP(&config.Exclude, "exclude", "x", []string{}, "exclude pattern(See https://golang.org/pkg/path/filepath/#Match)")
	RootCmd.Flags().StringVarP(&config.Shell, "shell", "s", os.Getenv("SHELL"), "shell to execute command")
	RootCmd.Flags().BoolVarP(&config.KillOnRestart, "kill-on-restart", "k", true, "kill on restart")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if persistentConfig.cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(persistentConfig.cfgFile)
	}

	viper.SetConfigName(".sentry") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
