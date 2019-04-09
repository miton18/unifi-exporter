package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/miton18/go-warp10/base"

	"github.com/miton18/go-warp10/instrumentation"
	"github.com/miton18/unifi-exporter/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd launch the nothing.
var RootCmd = &cobra.Command{
	Use: "unifi_exporter",
	Run: rootFn,
}

func init() {
	RootCmd.PersistentFlags().String("config", "", "config file to use")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	cobra.OnInitialize(initCobra)

	log.SetLevel(log.InfoLevel)
	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	err := RootCmd.Execute()
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func rootFn(cmd *cobra.Command, args []string) {
	host := viper.GetString("controller.host")
	port := viper.GetInt("controller.port")
	user := viper.GetString("controller.user")
	pass := viper.GetString("controller.password")
	refresh := viper.GetDuration("refresh.period")
	whitelist := viper.GetStringSlice("site.whitelist")
	warpHost := viper.GetString("warp10.host")
	warpToken := viper.GetString("warp10.token")
	listen := viper.GetString("listen")

	if user == "" {
		log.Fatal("cannot connect to Unifi controller without user name")
	}
	if pass == "" {
		log.Fatal("cannot connect to Unifi controller without user password")
	}

	gracefulStop := make(chan os.Signal, 2)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	data := make(chan instrumentation.Metrics)

	log.Infof("Starting HTTP server on: %s", listen)
	startServer(listen, data)

	client := client.New(user, pass, host, port, true)

	err := client.Connect()
	if err != nil {
		log.WithError(err).Fatal("cannot join Unifi controller")
	}

	log.Infof("Start scrapping Unifi controller every %s", refresh)
	t := time.NewTicker(refresh)

	// Initial scrappe
	metrics, err := scrappe(client, whitelist)
	if err != nil {
		log.WithError(err).Fatal("cannot scrape initial unifi controller metrics")
	}
	flush(warpHost, warpToken, metrics)
	log.Infof("metrics %+v", metrics.Get().Sensision())

	for {
		select {
		case <-t.C:
			metrics, err := scrappe(client, whitelist)
			if err != nil {
				log.WithError(err).Error("cannot scrape unifi controller metrics")
				continue
			}

			fmt.Println(metrics.Get().Sensision())
			flush(warpHost, warpToken, metrics)
			//data <- metrics
		case <-gracefulStop:
			goto exit
		}
	}

exit:
	log.Info("Stopping...")
	err = stopServer()
	if err != nil {
		log.WithError(err).Error("cannot close HTTP server")
	}
}

func flush(warp, token string, metrics instrumentation.Metrics) {
	if warp == "" || token == "" {
		return
	}

	wClient := base.NewClient(warp)
	wClient.WriteToken = token

	err := wClient.Update(metrics.Get())
	if err != nil {
		log.WithError(err).Error("cannot send metrics to Warp10")
	}
}