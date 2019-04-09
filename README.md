# Unifi exporter

Export metrics from your Unifi controller.
This application expose or push Sensision metrics into a [Warp10 platform](https://github.com/senx/warp10-platform).

- No dependancies to other libs
- multi-platform build

## Targets

OS:

- linux
- darwin

Arch:

- arm(5-6-7-64)
- amd64

## Metrics

|Metric| labels| attributs|Desc|
|------|-------|----------|----|
|unifi.system.status|siteId, system||System status, can be OK or warning|
|unifi.system.rx.bytes|siteId, system||Bytes per seconds received by system|
|unifi.system.tx.bytes|siteId, system||Bytes per seconds sent by system|
