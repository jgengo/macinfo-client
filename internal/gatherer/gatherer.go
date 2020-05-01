package gatherer

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/jgengo/macinfo-client/internal/utils"
)

// System stores the Mac system information
type System struct {
	Token      string    `json:"token"`
	Hostname   string    `json:"hostname"`
	ActiveUser string    `json:"active_user"`
	UUID       string    `json:"uuid"`
	Usb        []Usb     `json:"usb_devices"`
	Uptime     uint      `json:"uptime"`
	LastReboot []string  `json:"last_reboot"`
	Sensors    []Sensor  `json:"sensors"`
	OsVersion  OsVersion `json:"os_version"`
}

// OsVersion stores the os version information
type OsVersion struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

// Sensor stores a temperature sensor information
type Sensor struct {
	Name    string  `json:"name"`
	Celsius float64 `json:"celsius"`
}

// Usb stores a usb device information
type Usb struct {
	Vendor string `json:"vendor"`
	Model  string `json:"model"`
}

// GetInfo retrieves all the information of the client
func GetInfo() *System {
	var system System

	system.Token = utils.Cfg.APIToken

	system.getActiveUser()
	system.getLastReboot()
	system.getSystemInfo()
	system.getUsbDevices()
	system.getUptime()
	system.getTemperatureSensors()
	system.getOsVersion()

	return &system
}

func execToString(cmd *exec.Cmd) (string, error) {
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (s *System) getActiveUser() {
	resp, err := utils.OsQ.Client.Query("select user from logged_in_users where type='user' and tty='console' limit 1;")
	if err != nil {
		log.Printf("failed to get ActiveUser: %v\n", err)
		return
	}
	if len(resp.Response) > 0 {
		s.ActiveUser = resp.Response[0]["user"]
	} else {
		s.ActiveUser = ""
	}

}

// You may know that the command last reboot will get slower
// if you don't sometimes clean your asl log files.
func (s *System) getLastReboot() {
	out, err := execToString(exec.Command("last", "reboot"))
	if err != nil {
		log.Printf("failed to get the last reboots: %v\n", err)
		return
	}
	r, _ := regexp.Compile("[a-zA-Z]{3}\\s{1}[a-zA-Z]{3}\\s{1,2}\\d{1,2}\\s{1}\\d{2}\\:\\d{2}") // catches "Sun Apr 19 11:41"
	ret := r.FindAllString(out, -1)
	s.LastReboot = ret
}

func (s *System) getOsVersion() {
	resp, err := utils.OsQ.Client.Query("select version, build from os_version;")
	if err != nil {
		log.Printf("failed to get os info: %v\n", err)
		return
	}
	s.OsVersion = OsVersion{Version: resp.Response[0]["version"], Build: resp.Response[0]["build"]}
}

func (s *System) getTemperatureSensors() {
	resp, err := utils.OsQ.Client.Query("select name, celsius from temperature_sensors;")
	if err != nil {
		log.Printf("failed to get Sensors: %v\n", err)
		return
	}

	for _, sensor := range resp.Response {
		celsius, _ := strconv.ParseFloat(sensor["celsius"], 1)
		s.Sensors = append(s.Sensors, Sensor{Name: sensor["name"], Celsius: celsius})
	}
}

func (s *System) getUptime() {
	resp, err := utils.OsQ.Client.Query("select total_seconds from uptime")
	if err != nil {
		log.Printf("failed to get uptime: %v\n", err)
		return
	}
	conv, _ := strconv.Atoi(resp.Response[0]["total_seconds"])
	s.Uptime = uint(conv)
}

func (s *System) getSystemInfo() {
	resp, err := utils.OsQ.Client.Query("select uuid, hostname from system_info;")
	if err != nil {
		log.Printf("failed to get the system_info: %v\n", err)
		return
	}
	s.Hostname = resp.Response[0]["hostname"]
	s.UUID = resp.Response[0]["uuid"]
}

func (s *System) getUsbDevices() {
	resp, err := utils.OsQ.Client.Query("select vendor, model from usb_devices;")
	if err != nil {
		log.Printf("failed to get usbDevices: %v\n", err)
		return
	}

	for _, usb := range resp.Response {
		s.Usb = append(s.Usb, Usb{
			Vendor: strings.TrimSpace(usb["vendor"]),
			Model:  strings.TrimSpace(usb["model"]),
		})
	}
}
