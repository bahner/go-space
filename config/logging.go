package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultLogLevel = "info"
	defaultLogFile  = "." + NAME + ".log" // Make default log file hidden
)

func init() {

	// Logging
	pflag.String("loglevel", defaultLogLevel, "Loglevel to use for application")
	viper.BindPFlag("log.level", pflag.Lookup("loglevel"))
	viper.SetDefault("log.level", defaultLogLevel)

	pflag.String("logfile", defaultLogFile, "Logfile to write logs to")
	viper.BindPFlag("log.file", pflag.Lookup("logfile"))
	viper.SetDefault("log.file", defaultLogFile)
}
func initLogging() {

	// Init logger
	level, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	file, err := os.OpenFile(viper.GetString("log.file"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(73) // EX_CANTCREAT
	}
	log.SetOutput(file)
	log.Info("Logger initialized")

}
