package main

import (
	"archive/zip"
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"manager/fingerprint"
	"manager/profile"
	userProxy "manager/proxy"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/puzpuzpuz/xsync/v3"
	socks5 "github.com/things-go/go-socks5"
	"golang.org/x/net/proxy"
)

// App struct
type App struct {
	ctx             context.Context
	config          *Config
	proxies         *xsync.MapOf[int, *ProxyInstance]
	profilesFactory *profile.Factory
}

type Config struct {
	AuthCode int `json:"auth_code"`
}

func (a *App) Logout() {
	a.config.AuthCode = 0
	a.SaveConfig(a.config)
}
func browserExecutable() (p string, err error) {
	execPath := filepath.Join(browserPath, "sybilfox")
	switch os := runtime.GOOS; os {
	case "windows":
		execPath += ".exe"
	case "darwin":
		execPath = filepath.Join(browserPath, "Sybilfox.app", "Contents", "Resources", "sybilfox")
	case "linux":
		execPath += "-bin"
	}
	_, err = os.Stat(execPath)
	if err != nil {
		return
	}
	return execPath, nil
}
func (a *App) IsBrowserInstalled() bool {
	execPath := filepath.Join(browserPath, "sybilfox")
	switch os := runtime.GOOS; os {
	case "windows":
		execPath += ".exe"
	case "darwin":
		break
	case "linux":
		execPath += "-bin"
	}
	_, err := os.Stat(execPath)
	return err == nil
}
func unzipFromReader(rc io.ReadCloser, dest string) error {
	defer rc.Close()
	tmpFile, err := os.CreateTemp("", "temp-zip-*.zip")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Copy the HTTP response to the temporary file
	if _, err := io.Copy(tmpFile, rc); err != nil {
		defer tmpFile.Close()
		return err
	}
	defer tmpFile.Close()

	// Open ZIP archive
	zipReader, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to open ZIP archive: %w", err)
	}

	// Extract files
	for _, file := range zipReader.File {
		filePath := filepath.Join(dest, file.Name)

		// Check if it's a directory
		if file.FileInfo().IsDir() {
			// Ensure directory exists
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Ensure the parent directories exist
		err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create parent directories: %w", err)
		}

		// Remove the existing file if it exists to allow overwriting
		if _, err := os.Stat(filePath); err == nil {
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("failed to overwrite file: %w", err)
			}
		}

		// Open ZIP file
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in ZIP: %w", err)
		}
		defer srcFile.Close()

		// Create destination file (overwriting)
		destFile, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer destFile.Close()

		// Copy contents
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}
	}

	return nil
}

func (a *App) InstallBrowser() (err error) {
	archiveUrl := "https://github.com/reueeueeu/sybilfox-manager/releases/download/135/sybilfox-135-"
	switch runtime.GOOS {
	case "darwin":
		archiveUrl += "mac."
		if runtime.GOARCH == "amd64" {
			archiveUrl += "x86_64.zip"
		} else if runtime.GOARCH == "arm64" {
			archiveUrl += "arm64.zip"
		}
	case "linux":
		archiveUrl += "lin.x86_64.zip"
	case "windows":
		archiveUrl += "win.x86_64.zip"
	}
	fmt.Println(archiveUrl)
	resp, err := http.Get(archiveUrl)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	err = unzipFromReader(resp.Body, browserPath)
	if err != nil {
		return
	}
	return
}
func (a *App) IsAuthed() bool {
	return a.Auth(a.config.AuthCode) == nil
}
func (a *App) Auth(accessCode int) (err error) {
	fmt.Println(accessCode)
	authed, err := fingerprint.Auth(accessCode)
	if err != nil {
		return
	}
	if !authed {
		return errors.New("can't auth")
	}
	fmt.Println(accessCode)
	a.config.AuthCode = accessCode
	err = a.SaveConfig(a.config)
	return
}
func getOSValue() string {
	switch runtime.GOOS {
	case "windows":
		return "win"
	case "linux":
		return "lin"
	case "darwin":
		return "macos"
	default:
		return runtime.GOOS
	}
}

func updatePrefs(filePath string, port int, enabled bool) error {
	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		// File not exist creating predefined profile
		os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if enabled {
			err = os.WriteFile(filePath, []byte(fmt.Sprintf(`user_pref("webgl.force-enabled", true);
user_pref("network.proxy.socks", "127.0.0.1");
user_pref("network.proxy.socks_port", %d);
user_pref("network.proxy.socks_remote_dns", true);
user_pref("network.proxy.type", 1);
user_pref("browser.taskbar.grouping.useProfile", true);
user_pref("taskbar.grouping.useProfile", true);
user_pref("browser.taskbar.grouping.useprofile", true);
user_pref("taskbar.grouping.useprofile", true);
`, port)), os.ModePerm)
		} else {
			err = os.WriteFile(filePath, []byte(fmt.Sprintf(`user_pref("webgl.force-enabled", true);
		user_pref("network.proxy.socks", "127.0.0.1");
		user_pref("network.proxy.socks_port", %d);
		user_pref("network.proxy.socks_remote_dns", true);
		user_pref("network.proxy.type", 0);
		user_pref("browser.taskbar.grouping.useProfile", true);
		user_pref("taskbar.grouping.useProfile", true);
		user_pref("browser.taskbar.grouping.useprofile", true);
		user_pref("taskbar.grouping.useprofile", true);
		`, port)), os.ModePerm)
		}
		if err != nil {
			return err
		}
		return nil
	}
	defer file.Close()

	// Read the file line by line
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Regular expression to match the network.proxy.socks_port line
	portRegex := regexp.MustCompile(`(?m)^user_pref\("network.proxy.socks_port",\s*\d+\);`)

	// Replace the value of network.proxy.socks_port
	found := false
	for i, line := range lines {
		if portRegex.MatchString(line) {
			lines[i] = fmt.Sprintf(`user_pref("network.proxy.socks_port", %d);`, port)
			found = true
			break
		}
	}
	keysToDelete := []string{
		"browser.newtabpage.activity-stream.newtabWallpapers.v2.enabled",
		"browser.contentblocking.introCount",
		"browser.newtabpage.activity-stream.discoverystream.region-weather-config",
		"browser.newtabpage.enabled",
		"browser.search.update",
		//"browser.sessionstore.resume_from_crash",
		//"browser.startup.couldRestoreSession.count",
		"browser.tabs.unloadOnLowMemory",
		"browser.tabs.remote.unloadDelayMs",
		"browser.topsites.contile.enabled",
		"browser.urlbar.merino.endpointURL",
		"browser.urlbar.suggest.searches",
		"datareporting.healthreport.documentServerURI",
		"dom.navigation.navigationRateLimit.count",
		"dom.push.connection.enabled",
		"dom.screenorientation.allow-lock",
		"dom.successive_dialog_time_limit",
		"extensions.blocklist.detailsURL",
		"extensions.blocklist.itemURL",
		"extensions.getAddons.discovery.api_url",
		"extensions.getAddons.get.url",
		"extensions.getAddons.search.browseURL",
		"extensions.hotfix.url",
		"extensions.installDistroAddons",
		"extensions.update.background.url",
		"extensions.update.url",
		"geo.provider.network.url",
		"identity.fxaccounts.auth.uri",
		"remote.prefs.recommended.applied",
		"security.remote_settings.intermediates.enabled",
		"services.settings.server",
		"signon.autofillForms",
		"signon.rememberSignons",
		"startup.homepage_welcome_url",
		"toolkit.telemetry.server",
		"widget.windows.window_occlusion_tracking.enabled",
	}
	if !found {
		if enabled {
			proxyConf := strings.Split(fmt.Sprintf(`user_pref("network.proxy.socks", "127.0.0.1");
user_pref("network.proxy.socks_port", %d);
user_pref("network.proxy.socks_remote_dns", true);
user_pref("network.proxy.type", 1);
user_pref("browser.taskbar.grouping.useProfile", true);
user_pref("taskbar.grouping.useProfile", true);
user_pref("browser.taskbar.grouping.useprofile", true);
user_pref("taskbar.grouping.useprofile", true);
`, port), "\n")
			lines = append(lines, proxyConf...)
		} else {

			proxyConf := strings.Split(fmt.Sprintf(`user_pref("network.proxy.socks", "127.0.0.1");
		user_pref("network.proxy.socks_port", %d);
		user_pref("network.proxy.socks_remote_dns", true);
		user_pref("network.proxy.type", 0);
		user_pref("browser.taskbar.grouping.useProfile", true);
		user_pref("taskbar.grouping.useProfile", true);
		user_pref("browser.taskbar.grouping.useprofile", true);
		user_pref("taskbar.grouping.useprofile", true);
		`, port), "\n")
			lines = append(lines, proxyConf...)
		}

	}

	// Open the file for writing
	file, err = os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %w", err)
	}
	defer file.Close()

	// Write the updated lines back to the file
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		blacklisted := false
		for _, banned := range keysToDelete {
			if strings.Contains(line, banned) {
				blacklisted = true
				break
			}
		}
		if blacklisted {
			continue
		}
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing data to file: %w", err)
	}

	return nil
}

type ProxyInstance struct {
	Server       *socks5.Server
	ProfileProxy *ProfileProxy
	Addr         *net.TCPAddr
}

type ProfileProxy struct {
	lock sync.RWMutex
	proxy.Dialer
}

func (pp *ProfileProxy) updateProxy(pc userProxy.Config) {
	pp.lock.Lock()
	defer pp.lock.Unlock()
	var proxyAuth *proxy.Auth
	if pc.AuthEnabled() {
		proxyAuth = &proxy.Auth{
			User:     pc.User,
			Password: pc.Password,
		}
	}
	profileProxy, err := proxy.SOCKS5("tcp", pc.Host, proxyAuth, nil)
	if err != nil {
		return
	}
	pp.Dialer = profileProxy
}

func dailerFromProxy(pc userProxy.Config) *ProfileProxy {
	pp := &ProfileProxy{}
	pp.updateProxy(pc)
	return pp
}
func makeVars(configData []byte) []string {
	var envVars []string

	// Determine chunk size based on OS
	var chunkSize int
	if os.Getenv("OS") == "Windows_NT" {
		chunkSize = 2047
	} else {
		chunkSize = 32767
	}

	configStr := string(configData)
	for i := 0; i < len(configStr); i += chunkSize {
		end := i + chunkSize
		if end > len(configStr) {
			end = len(configStr)
		}
		chunk := configStr[i:end]
		envName := fmt.Sprintf("SYBIL_CONFIG_%d", (i/chunkSize)+1)
		envVars = append(envVars, envName+"="+chunk)
	}
	return envVars
}
func (a *App) RunConfig(id int) (err error) {
	c, err := a.profilesFactory.GetConfig(id)
	if err != nil {
		return
	}
	profilePath := filepath.Join(profilesPath, strconv.Itoa(c.ID))
	if c.Request.Proxy.IsEnabled() {
		proxyInstance, loaded := a.proxies.LoadOrCompute(c.ID, func() *ProxyInstance {
			pp := dailerFromProxy(c.Request.Proxy)
			browserProxy := socks5.NewServer(socks5.WithDial(func(ctx context.Context, network, addr string) (net.Conn, error) {
				pp.lock.RLock()
				defer pp.lock.RUnlock()
				return pp.Dial(network, addr)
			}))
			l, err := net.Listen("tcp", ":0")
			if err != nil {
				return nil
			}
			addr := l.Addr().(*net.TCPAddr)
			go browserProxy.Serve(l)
			return &ProxyInstance{
				ProfileProxy: pp,
				Server:       browserProxy,
				Addr:         addr,
			}
		})
		if !loaded || proxyInstance == nil {
			return errors.New("proxy not correctly loaded")
		}
		err = updatePrefs(filepath.Join(profilePath, "prefs.js"), proxyInstance.Addr.Port, true)
	} else {
		err = updatePrefs(filepath.Join(profilePath, "prefs.js"), 0, false)
	}
	if err != nil {
		return err
	}
	hc := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := hc.Get(fingerprint.BackendUrl + "fingerprint/get/" + c.Fingerprint.Fingerprint.ID)
	if err != nil {
		return errors.New("can't fetch fingerprint")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("can't fetch fingerprint")
	}
	p, err := browserExecutable()
	if err != nil {
		return
	}
	cmd := exec.Command(p, "-profile", profilePath, "-no-remote", "-wait-for-browser", "-foreground")
	cmd.Env = append(os.Environ(), makeVars(body)...)
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		return
	}
	return
}
func (a *App) HideConfig(id int) (err error) {
	c, err := a.profilesFactory.GetConfig(id)
	if err != nil {
		return
	}
	c.Hidden = true
	err = a.profilesFactory.SaveConfig(c)
	return
}
func (a *App) ListConfigs() []profile.Config {
	return a.profilesFactory.ListConfigs(a.config.AuthCode)
}
func (a *App) UpdateConfig(id int, request profile.Request) (updated profile.Config, err error) {
	config, err := a.profilesFactory.GetConfig(id)
	if err != nil {
		return
	}
	updated, err = a.profilesFactory.UpdateConfig(config, request)
	if err != nil {
		return
	}
	updated.Request.Imported = false
	err = a.profilesFactory.SaveConfig(updated)
	if err != nil {
		return
	}
	i, exist := a.proxies.Load(id)
	if exist {
		i.ProfileProxy.updateProxy(updated.Request.Proxy)
	}
	return updated, nil
}
func (a *App) NewConfig(name string, pc userProxy.Config) (c profile.Config, err error) {
	c, err = a.profilesFactory.NewConfig(profile.Request{
		AccessCode: a.config.AuthCode,
		Proxy:      pc,
		Name:       name,
	})
	if err != nil {
		return
	}
	err = a.profilesFactory.SaveConfig(c)
	return
}
func (a *App) GetConfig() (c *Config) {
	return a.config
}
func (a *App) SaveConfig(c *Config) (err error) {
	data, err := json.Marshal(c)
	if err != nil {
		log.Fatalln("can't write config file at ./config.json")
	}
	err = os.WriteFile(filepath.Join(appDir, "./config.json"), data, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

var userDir, userDirErr = os.UserConfigDir()
var appDir = filepath.Join(userDir, "SybilfoxManager")
var browserPath = filepath.Join(appDir, "./browser")
var configsPath = filepath.Join(appDir, "./configs")
var profilesPath = filepath.Join(appDir, "./profiles")

// NewApp creates a new App application struct
func NewApp() *App {
	if userDirErr != nil {
		log.Fatalln(userDirErr)
	}
	c := new(Config)
	data, err := os.ReadFile(filepath.Join(appDir, "./config.json"))
	if err == nil {
		fmt.Println("config exist. reading it")
		json.Unmarshal(data, c)
	}
	os.MkdirAll(appDir, os.ModePerm)
	os.MkdirAll(browserPath, os.ModePerm)
	os.MkdirAll(configsPath, os.ModePerm)
	os.MkdirAll(profilesPath, os.ModePerm)
	a := &App{
		config:          c,
		profilesFactory: profile.NewFactory(profilesPath, configsPath, browserPath),
		proxies:         xsync.NewMapOf[int, *ProxyInstance](),
	}
	err = a.SaveConfig(c)
	if err != nil {
		log.Fatalln(err)
	}
	return a
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
