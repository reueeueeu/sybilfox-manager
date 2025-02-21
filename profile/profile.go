package profile

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"manager/fingerprint"
	"manager/proxy"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type Request struct {
	Name       string       `json:"name"`
	AccessCode int          `json:"access_code"`
	Proxy      proxy.Config `json:"proxy"`
	Imported   bool         `json:"imported"`
}

type Config struct {
	ID          int                `json:"id"`
	Fingerprint fingerprint.Record `json:"fingerprint"`
	Request     Request            `json:"request"`
	Hidden      bool               `json:"hidden"`
}

type Factory struct {
	//profilesPath string
	configsPath string
	browserPath string
}

func NewFactory(profilesPath, configsPath, browserPath string) *Factory {
	return &Factory{
		configsPath: configsPath,
		//profilesPath: profilesPath,
		browserPath: browserPath,
	}
}

func (f *Factory) ListConfigs(accessKey int) (configs []Config) {
	configs = make([]Config, 0)
	filepath.Walk(f.configsPath, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			fmt.Println(filepath.Ext(path))
			return nil
		}
		var profile Config
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		err = json.Unmarshal(data, &profile)
		if err != nil {
			return nil
		}
		if profile.Hidden {
			return nil
		}
		configs = append(configs, profile)
		return nil
	})
	if len(configs) == 0 && accessKey != -1 {
		hc := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		res, err := hc.Post(fingerprint.BackendUrl+"fingerprints/"+strconv.Itoa(accessKey), "application/json", nil)
		if err != nil {
			return
		}
		var fingerprints []fingerprint.Fingerprint
		d := json.NewDecoder(res.Body)
		err = d.Decode(&fingerprints)
		if err != nil {
			return
		}
		for i, fp := range fingerprints {
			c := Config{
				ID: i,
				Fingerprint: fingerprint.Record{
					Fingerprint:   fp,
					Country:       "",
					ProxyExitHost: fp.Ipv4,
				},
				Request: Request{
					Name:       fmt.Sprintf("Imported %d", i),
					AccessCode: accessKey,
					Proxy:      proxy.Config{},
					Imported:   true,
				},
				Hidden: false,
			}
			err = f.SaveConfig(c)
			if err != nil {
				continue
			}
			configs = append(configs, c)
		}
	}
	return
}

func (f *Factory) LastConfig() (config Config, err error) {
	configs := f.ListConfigs(-1)
	if len(configs) == 0 {
		return Config{}, errors.New("no profiles exist")
	}
	return configs[len(configs)-1], nil
}

func (f *Factory) GetConfig(id int) (config Config, err error) {
	configs := f.ListConfigs(-1)
	if len(configs) == 0 {
		return Config{}, errors.New("no profiles exist")
	}
	for _, config := range configs {
		if config.ID == id {
			return config, nil
		}
	}
	return Config{}, errors.New("no config found with index " + strconv.Itoa(id))
}

func (f *Factory) NewConfig(req Request) (config Config, err error) {
	id := 0
	if c, err := f.LastConfig(); err == nil {
		id = c.ID + 1
	}
	return newConfig(id, req)
}

func (f *Factory) UpdateConfig(config Config, updateReq Request) (updated Config, err error) {
	return updateConfig(config, updateReq)
}
func (f *Factory) SaveConfig(config Config) (err error) {
	data, err := json.Marshal(config)
	if err != nil {
		return
	}
	err = os.WriteFile(filepath.Join(f.configsPath, fmt.Sprintf("%d.json", config.ID)), data, os.ModePerm)
	return
}
func newConfig(id int, req Request) (config Config, err error) {
	fp, err := fingerprint.New(req.AccessCode, req.Proxy)
	if err != nil {
		return
	}
	return Config{
		ID:          id,
		Fingerprint: fp,
		Request:     req,
	}, nil
}

func updateConfig(c Config, req Request) (config Config, err error) {
	fp, err := fingerprint.Update(c.Fingerprint.Fingerprint.ID, req.AccessCode, req.Proxy)
	if err != nil {
		return
	}
	c.Fingerprint = fp
	c.Request = req
	return c, nil
}
