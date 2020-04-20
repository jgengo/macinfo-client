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
	Vendor string
	Model  string
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

func execToString(cmd *exec.Cmd) string {
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		log.Fatalf("error while gathering last reboots: %v\n", err)
	}
	return b.String()
}

func (s *System) getActiveUser() {
	resp, err := utils.OsQ.Client.Query("select user from logged_in_users where tty='console' limit 1;")
	if err != nil {
		log.Fatalf("error while gathering active user: %v", err)
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
	var ret []string

	out := execToString(exec.Command("last", "reboot"))
	lines := strings.Split(out, "\n")

	r, _ := regexp.Compile("[a-zA-Z]{3}\\s{1}[a-zA-Z]{3}\\s{1,2}\\d{1,2}\\s{1}\\d{2}\\:\\d{2}")
	for _, line := range lines {
		findStr := r.FindString(line)
		if findStr != "" {
			ret = append(ret, findStr)
		}
	}
	s.LastReboot = ret
}

func (s *System) getOsVersion() {
	resp, err := utils.OsQ.Client.Query("select version, build from os_version;")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.OsVersion = OsVersion{Version: resp.Response[0]["version"], Build: resp.Response[0]["build"]}
}

func (s *System) getTemperatureSensors() {
	resp, err := utils.OsQ.Client.Query("select name, celsius from temperature_sensors;")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, sensor := range resp.Response {
		celsius, _ := strconv.ParseFloat(sensor["celsius"], 1)
		s.Sensors = append(s.Sensors, Sensor{Name: sensor["name"], Celsius: celsius})
	}
}

func (s *System) getUptime() {
	resp, err := utils.OsQ.Client.Query("select total_seconds from uptime")
	if err != nil {
		log.Fatalf("Error while gathering SysInfo: %v\n", err)
	}
	conv, _ := strconv.Atoi(resp.Response[0]["total_seconds"])
	s.Uptime = uint(conv)
}

func (s *System) getSystemInfo() {
	resp, err := utils.OsQ.Client.Query("select uuid, hostname from system_info;")
	if err != nil {
		log.Fatalf("error while trying to get the system_info: %v\n", err)
	}
	s.Hostname = resp.Response[0]["hostname"]
	s.UUID = resp.Response[0]["uuid"]
}

func (s *System) getUsbDevices() {
	resp, err := utils.OsQ.Client.Query("select vendor, model from usb_devices;")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, usb := range resp.Response {
		s.Usb = append(s.Usb, Usb{Vendor: usb["vendor"], Model: usb["model"]})
	}
}
