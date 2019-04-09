package main

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initCobra() {
	viper.AddConfigPath("/etc/unifi_exporter")
	viper.AddConfigPath("$HOME/.unifi_exporter")
	viper.AddConfigPath("/etc/default/unifi_exporter")
	viper.AddConfigPath(".")

	viper.SetEnvKeyReplacer(strings.NewReplacer("_", "."))
	viper.SetEnvPrefix("unifi")

	viper.SetDefault("controller.host", "localhost")
	viper.SetDefault("controller.port", "8080")
	viper.SetDefault("controller.version", 5)
	viper.SetDefault("refresh.period", 30*time.Second)
	viper.SetDefault("site.whitelist", []string{})
	viper.SetDefault("listen", "localhost:9101")

	err := viper.BindPFlags(RootCmd.PersistentFlags())
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindPFlags(RootCmd.Flags())
	if err != nil {
		log.Fatal(err)
	}

	// Load user defined config
	cfgFile := viper.GetString("config")
	if cfgFile != "" {
		log.Infof("Using configuration file '%s'", cfgFile)
		viper.SetConfigFile(cfgFile)
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatal("cannot load configuration file")
	}

	err = viper.MergeInConfig()
	if err != nil {
		log.WithError(err).Fatal("cannot merge configurations")
	}
}
