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

	Site struct {
		ID          string `json:"_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	// Health is a subsystem health status
	Health struct {
		SubSytem        string `json:"subsystem"`
		Status          string `json:"status"`
		RxBytes         int64  `json:"rx_bytes-r"`
		TxBytes         int64  `json:"tx_bytes-r"`
		Drops           int64  `json:"drops"`
		Latency         int64  `json:"latency"`
		SpeedTestPing   int64  `json:"speedtest_ping"`
		SpeedTestStatus string `json:"speedtest_status"`
		LanUsers        int64  `json:"num_user"`
	}

	Sta struct {
		BytesRxLan      int64  `json:"rx_bytes"`
		BytesTxLan      int64  `json:"tx_bytes"`
		BytesRxRetryLan int64  `json:"rx_bytes-r"`
		BytesTxRetryLan int64  `json:"tx_bytes-r"`
		BytesRxWan      int64  `json:"wired-rx_bytes"`
		BytesTxWan      int64  `json:"wired-tx_bytes"`
		BytesTxRetryWan int64  `json:"wired-tx_bytes-r"`
		BytesRxRetryWan int64  `json:"wired-rx_bytes-r"`
		MAC             string `json:"mac"`
		Hostname        string `json:"hostname"`
		Name            string `json:"name"`
		OUI             string `json:"oui"`
	}
)
