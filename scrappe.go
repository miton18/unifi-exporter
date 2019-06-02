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

			log.Infof("Drops %+v", systemHealth.Drops)
			if systemHealth.Drops > 0 {
				drops := instrumentation.NewGauge("unifi.system.drops", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "Bytes dropped")
				drops.Set(uint64(systemHealth.Drops))
				metrics = append(metrics, drops)
			}

			log.Infof("Latency %+v", systemHealth.Latency)
			if systemHealth.Latency > 0 {
				lat := instrumentation.NewGauge("unifi.system.latency", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "Latency")
				lat.Set(uint64(systemHealth.Latency))
				metrics = append(metrics, lat)
			}

			log.Infof("Speed test ping %+v", systemHealth.SpeedTestPing)
			if systemHealth.SpeedTestPing > 0 {
				sPing := instrumentation.NewGauge("unifi.system.speedtest.latency", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "Speed test latency")
				sPing.Set(uint64(systemHealth.SpeedTestPing))
				metrics = append(metrics, sPing)
			}

			log.Infof("Speed test Status %+v", systemHealth.SpeedTestStatus)
			if systemHealth.SpeedTestStatus != "" {
				sStatus := instrumentation.NewState("unifi.system.speedtest.status", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "Speed test status")
				sStatus.Set(systemHealth.SpeedTestStatus)
				metrics = append(metrics, sStatus)
			}

			log.Infof("clients %+v", systemHealth.LanUsers)
			if systemHealth.LanUsers > 0 {
				clients := instrumentation.NewGauge("unifi.system.clients", base.Labels{
					"siteId": site.ID,
					"system": systemHealth.SubSytem,
				}, "connected clients")
				clients.Set(uint64(systemHealth.LanUsers))
				metrics = append(metrics, clients)
			}
		}

		stas, err := client.Sta(site.Name)
		if err != nil {
			return nil, errors.Wrap(err, "cannot get clients stats")
		}

		for _, sta := range stas {
			clientBytesRx := instrumentation.NewGauge("unifi.client.bytes.rx", base.Labels{
				"siteId":   site.ID,
				"mac":      sta.MAC,
				"hostname": sta.Hostname,
				"name":     sta.Name,
				"oui":      sta.OUI,
			}, "Client bytes received")
			clientBytesRx.Set(uint64(sta.BytesRxLan))
			metrics = append(metrics, clientBytesRx)

			clientBytesS := instrumentation.NewGauge("unifi.client.bytes.tx", base.Labels{
				"siteId":   site.ID,
				"mac":      sta.MAC,
				"hostname": sta.Hostname,
				"name":     sta.Name,
				"oui":      sta.OUI,
			}, "Client bytes sent")
			clientBytesS.Set(uint64(sta.BytesTxLan))
			metrics = append(metrics, clientBytesS)

			clientBytesR := instrumentation.NewGauge("unifi.client.bytes.retry", base.Labels{
				"siteId":   site.ID,
				"mac":      sta.MAC,
				"hostname": sta.Hostname,
				"name":     sta.Name,
				"oui":      sta.OUI,
			}, "Client bytes retry")
			clientBytesR.Set(uint64(sta.BytesRetry))
			metrics = append(metrics, clientBytesR)

			if !sta.IsWired {
				clientBytesSat := instrumentation.NewGauge("unifi.client.satisfaction", base.Labels{
					"siteId":   site.ID,
					"mac":      sta.MAC,
					"hostname": sta.Hostname,
					"name":     sta.Name,
					"oui":      sta.OUI,
				}, "Client satisfaction")
				clientBytesSat.Set(uint64(sta.Satisfaction))
				metrics = append(metrics, clientBytesSat)
			}
		}
	}

	return metrics, nil
}
