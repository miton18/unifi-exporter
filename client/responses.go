package client

import "encoding/json"

type (
	// UnifiResponse is a Unifi standard API response
	UnifiResponse struct {
		Data json.RawMessage `json:"data"`
		Meta struct {
			RC string `json:"rc"`
		} `json:"meta"`
	}

	// SiteList is a collection of sites
	SiteList []Site

	Site struct {
		ID          string `json:"_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// Health is a collection of system health
	Health []SytemHealth

	// SytemHealth is a subsystem health status
	SytemHealth struct {
		SubSytem string `json:"subsystem"`
		Status   string `json:"status"`
		RxBytes  int64  `json:"rx_bytes-r"`
		TxBytes  int64  `json:"tx_bytes-r"`
	}
)
