package main

import (
	"github.com/miton18/go-warp10/base"
	"github.com/miton18/go-warp10/instrumentation"
	"github.com/miton18/unifi-exporter/client"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func scrappe(client *client.Client, siteWhitelist []string) (instrumentation.Metrics, error) {
	metrics := instrumentation.Metrics{}

	sites, err := client.Sites()
	if err != nil {
		return nil, errors.Wrap(err, "cannot list sites")
	}

	for _, site := range sites {
		logrus.Debug("Site: " + site.Name)

		if len(siteWhitelist) > 0 && !in(site.Name, siteWhitelist) {
			log.Debugf("Skip filtered site '%s'", site.Name)
			continue
		}

		health, err := client.Health(site.Name)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get site health")
		}
		log.Infof("%+v", health)

		for _, systemHealth := range health {
			healthState := instrumentation.NewState("unifi.system.status", base.Labels{
				"siteId": site.ID,
				"system": systemHealth.SubSytem,
			}, "System state")
			healthState.Set(systemHealth.Status)
			metrics = append(metrics, healthState)

			if systemHealth.TxBytes > 0 {
				tx := instrumentation.NewGauge("unifi.system.tx.bytes", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "Bytes send")
				tx.Set(uint64(systemHealth.TxBytes))
				metrics = append(metrics, tx)
			}

			log.Infof("RX %+v", systemHealth.RxBytes)
			if systemHealth.RxBytes > 0 {
				rx := instrumentation.NewGauge("unifi.system.rx.bytes", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "Bytes Received")
				rx.Set(uint64(systemHealth.RxBytes))
				metrics = append(metrics, rx)
			}
		}
	}

	return metrics, nil
}
