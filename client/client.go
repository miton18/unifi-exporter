package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"strconv"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type (
	Client struct {
		httpClient *http.Client
		host       string
		user       string
		pass       string
	}
)

// New Unifi client
func New(user, pass, host string, port int, httpClient *http.Client, secure bool) *Client {
	log.Info(user + pass + host)

	jar, _ := cookiejar.New(nil)

	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient.Jar = jar

	proto := "http://"
	if secure {
		proto = "https://"
	}

	return &Client{
		httpClient: httpClient,
		host:       proto + host + ":" + strconv.Itoa(port),
		user:       user,
		pass:       pass,
	}
}

// Connect Unifi client
func (c *Client) Connect() error {
	b, err := json.Marshal(map[string]string{
		"username": c.user,
		"password": c.pass,
	})
	if err != nil {
		return err
	}

	br := bytes.NewReader(b)

	req, err := http.NewRequest("POST", c.host+"/api/login", br)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		b, err = httputil.DumpResponse(res, true)
		if err != nil {
			return errors.Wrapf(err, "Failed to call Unifi controller: %s", res.Status)
		}
		return fmt.Errorf("Invalid response from controller: %+v", string(b))
	}

	return nil
}

func (c *Client) request(method, path string, out interface{}) error {
	req, err := http.NewRequest(method, c.host+path, nil)
	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	resB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.WithError(err).Error("cannot close unifi response body")
		}
	}()

	var uRes UnifiResponse
	err = json.Unmarshal(resB, &uRes)
	if err != nil {
		return err
	}

	err = json.Unmarshal(uRes.Data, out)
	if err != nil {
		return err
	}

	return nil
}

// Sites List unifi sites
func (c *Client) Sites() ([]Site, error) {
	var siteList []Site
	err := c.request("GET", "/api/self/sites", &siteList)
	if err != nil {
		return nil, err
	}

	return siteList, nil
}

// Health call
func (c *Client) Health(site string) ([]Health, error) {
	var health []Health
	url := fmt.Sprintf("/api/s/%s/stat/health", site)

	err := c.request("GET", url, &health)
	if err != nil {
		return nil, err
	}

	return health, nil
}

// Sta for each client
func (c *Client) Sta(site string) ([]Sta, error) {
	var stas []Sta
	url := fmt.Sprintf("/api/s/%s/stat/sta", site)

	err := c.request("GET", url, &stas)
	if err != nil {
		return nil, err
	}

	return stas, nil
}
