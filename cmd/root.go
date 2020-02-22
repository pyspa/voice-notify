package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/pyspa/voice-notify/log"
	"github.com/pyspa/voice-notify/service/dbus"
	"github.com/pyspa/voice-notify/service/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile = ""
)

var rootCmd = &cobra.Command{
	Use:     "voice-notify",
	Version: "0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func start() {

	log.Debug("debug mode", nil)
	if viper.GetString("pushbullet.access_token") != "" {
		s := pb.NewService()
		go func() {
			s.Start()
		}()
	}

	if runtime.GOOS == "linux" {
		s := dbus.NewService()
		go func() {
			s.Start()
		}()
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLogger)
	home, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	p := filepath.Join(home, "voice-notify.toml")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", p, "config file")
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("speech.lang", "ja-JP")
	viper.SetDefault("speech.speaking_rate", 1.5)
	viper.SetDefault("speech.pitch", 1.5)
	viper.SetDefault("speech.text_max", 256)

	viper.SetDefault("log.debug", false)
	viper.SetDefault("log.log", "stderr")
	viper.SetDefault("log.err_log", "stderr")

	if fileExists(cfgFile) {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
	}
}

func initLogger() {
	if err := log.Init(); err != nil {
		panic(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
