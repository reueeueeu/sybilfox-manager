package fingerprint

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"manager/proxy"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	netProxy "golang.org/x/net/proxy"
)

type response struct {
	Status      bool        `json:"status"`
	Fingerprint Fingerprint `json:"fingerprint"`
	Error       string      `json:"error"`
}
type Config struct {
	Timezone      string `json:"timezone"`
	WebGlRenderer string `json:"webGl:renderer"`
	WebGlVendor   string `json:"webGl:vendor"`
}
type Fingerprint struct {
	Version    int    `json:"version"`
	Ipv4       string `json:"ipv4"`
	Config     Config `json:"config"`
	Ipv6       string `json:"ipv6"`
	Os         string `json:"os"`
	Key        int    `json:"key"`
	ID         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	AccessedAt any    `json:"accessed_at"`
}

type Record struct {
	Fingerprint   Fingerprint `json:"fingerprint"`
	Country       string      `json:"country"`
	ProxyExitHost string      `json:"proxy_exit_host"`
}

// const BackendUrl = "https://192.168.140.200:14443/"
const BackendUrl = "https://164.90.199.203:14443/"

var u, err = url.Parse(BackendUrl)

func init() {
	if err != nil {
		log.Fatalln(err)
	}
}

type fingerprintRequest struct {
	Key     int    `json:"key"`
	Os      string `json:"os"`
	Version int    `json:"version"`
	IPv4    string `json:"ipv4"`
	IPv6    string `json:"ipv6"`
}

func ipGet(what string, proxy proxy.Config) (string, error) {
	// Set up SOCKS5 proxy authentication
	auth := netProxy.Auth{
		User:     proxy.User,
		Password: proxy.Password,
	}
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if proxy.IsEnabled() {
		var err error
		var dialer netProxy.Dialer
		if proxy.AuthEnabled() {
			// Create a SOCKS5 dialer
			dialer, err = netProxy.SOCKS5("tcp", proxy.Host, &auth, netProxy.Direct)
			if err != nil {
				return "", fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
			}

		} else {
			// Create a SOCKS5 dialer
			dialer, err = netProxy.SOCKS5("tcp", proxy.Host, nil, netProxy.Direct)
			if err != nil {
				return "", fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
			}
		}
		httpTransport.Dial = func(network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}

	// Create HTTP client with transport and timeout
	client := &http.Client{
		Transport: httpTransport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get(BackendUrl + what)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return strings.ReplaceAll(string(body), "\"", ""), nil
}

func New(accessKey int, pc proxy.Config) (fp Record, err error) {
	exitHost, err := ipGet("ip", pc)
	if err != nil {
		return
	}
	country, err := ipGet("country", pc)
	if err != nil {
		return
	}
	c, err := tls.Dial("tcp", u.Host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	defer c.Close()
	fpReq := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(fpReq)
	err = encoder.Encode(fingerprintRequest{
		Key:     accessKey,
		Os:      "windows",
		Version: 135,
		IPv4:    exitHost,
		IPv6:    "",
	})
	if err != nil {
		return
	}
	fpReqByte := fpReq.Bytes()
	req, err := http.NewRequest("POST", BackendUrl+"fingerprint/create", bytes.NewBuffer(fpReqByte))
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json")
	err = req.Write(c)
	if err != nil {
		return
	}
	res, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(res.Body)
	var resp response
	err = decoder.Decode(&resp)
	if err != nil {
		return
	}
	if !resp.Status {
		return Record{}, errors.New(resp.Error)
	}
	fp.Country = country
	fp.ProxyExitHost = exitHost
	fp.Fingerprint = resp.Fingerprint
	return
}
func Update(id string, accessKey int, pc proxy.Config) (fp Record, err error) {
	exitHost, err := ipGet("ip", pc)
	if err != nil {
		return
	}
	country, err := ipGet("country", pc)
	if err != nil {
		return
	}
	c, err := tls.Dial("tcp", u.Host, &tls.Config{
		InsecureSkipVerify: true,
	})
	fmt.Println(country)
	fmt.Println(exitHost)
	if err != nil {
		return
	}
	defer c.Close()
	fpReq := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(fpReq)
	err = encoder.Encode(fingerprintRequest{
		Key:     accessKey,
		Os:      "windows",
		Version: 135,
		IPv4:    exitHost,
		IPv6:    "",
	})
	if err != nil {
		return
	}
	fpReqByte := fpReq.Bytes()
	req, err := http.NewRequest("POST", BackendUrl+"fingerprint/"+id+"/update", bytes.NewBuffer(fpReqByte))
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json")
	err = req.Write(c)
	if err != nil {
		return
	}
	res, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(res.Body)
	var resp response
	err = decoder.Decode(&resp)
	if err != nil {
		return
	}
	if !resp.Status {
		return Record{}, errors.New(resp.Error)
	}
	fp.Country = country
	fp.ProxyExitHost = exitHost
	fp.Fingerprint = resp.Fingerprint
	return
}
func Auth(accessKey int) (authed bool, err error) {
	c, err := tls.Dial("tcp", u.Host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	defer c.Close()
	fpReq := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(fpReq)
	err = encoder.Encode(fingerprintRequest{
		Key: accessKey,
	})
	if err != nil {
		return
	}
	fpReqByte := fpReq.Bytes()
	req, err := http.NewRequest("POST", BackendUrl+"auth", bytes.NewBuffer(fpReqByte))
	if err != nil {
		return
	}
	req.Header.Add("content-type", "application/json")
	err = req.Write(c)
	if err != nil {
		return
	}
	res, err := http.ReadResponse(bufio.NewReader(c), req)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		return false, errors.New("can't auth")
	}
	decoder := json.NewDecoder(res.Body)
	var resp response
	err = decoder.Decode(&resp)
	if err != nil {
		return
	}
	if !resp.Status {
		return false, errors.New(resp.Error)
	}
	return true, nil
}
