package main

import (
	"fmt"
	"os"

	"github.com/bluemir/tick/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var printVersion bool
var cfgFile string
var config *core.Config

func _main(cmd *cobra.Command, args []string) {
	if printVersion {
		fmt.Printf("tick version %s.%d", __VERSION__, __BUILD_NUM__)
	} else {
		core.NewTick(config).Run()
	}

}

var RootCmd = &cobra.Command{
	Use:   "tick",
	Short: "file system watcher",
	Long:  `file System watcher for varius command`,
	Run:   _main,
}

func init() {
	config = &core.Config{}

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tick.yaml)")

	RootCmd.Flags().StringVarP(&config.Command, "command", "c", "", "command to run")
	RootCmd.Flags().StringVarP(&config.Exclude, "exclude", "x", "", "exculde pattern regexp")
	RootCmd.Flags().Int32VarP(&config.Delay, "delay", "d", 500, "delay that wait events")
	RootCmd.Flags().StringVarP(&config.Path, "path", "p", "./", "path")
	RootCmd.Flags().StringVarP(&config.Shell, "shell", "s", os.Getenv("SHELL"), "shell")
	RootCmd.Flags().BoolVarP(&config.KillOnRestart, "kill-on-restart", "k", true, "kill on restart")
	RootCmd.Flags().BoolVar(&printVersion, "version", false, "show version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".tick") // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

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
