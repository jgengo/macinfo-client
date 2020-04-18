package gatherer

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/jgengo/macinfo-client/internal/utils"
)

// System stores the Mac system information
type System struct {
	Hostname   string    `json:"hostname"`
	UUID       string    `json:"uuid"`
	Usb        []Usb     `json:"usb_devices"`
	Uptime     uint      `json:"uptime"`
	LastReboot []string  `json:"last_reboot"`
	Sensors    []Sensor  `json:"sensors"`
	OsVersion  OsVersion `json:"os_version"`
}

// OsVersion stores the os version information
type OsVersion struct {
	version string
	build   string
}

// Sensor stores a temperature sensor information
type Sensor struct {
	name    string
	celsius float64
}

// Usb stores a usb device information
type Usb struct {
	vendor string
	model  string
}

// GetInfo retrieves all the information of the client
func GetInfo(c *utils.OsQuery) {
	var system System
	system.getSystemInfo(c)
	system.getUsbDevices(c)
	system.getUptime(c)
	system.getLastReboot(c)
	system.getTemperatureSensors(c)
	system.getOsVersion(c)
}

func execToString(cmd *exec.Cmd) string {
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		log.Fatalf("error while gathering last reboots: %v\n", err)
	}
	return b.String()
}

func (s *System) getOsVersion(c *utils.OsQuery) {
	resp, err := c.Client.Query("select version, build from os_version;")
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.OsVersion = OsVersion{version: resp.Response[0]["version"], build: resp.Response[0]["build"]}
}

func (s *System) getTemperatureSensors(c *utils.OsQuery) {
	resp, err := c.Client.Query("select name, celsius from temperature_sensors;")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, sensor := range resp.Response {
		celsius, _ := strconv.ParseFloat(sensor["celsius"], 1)
		s.Sensors = append(s.Sensors, Sensor{name: sensor["name"], celsius: celsius})
	}
}

// TODO: refacto
func (s *System) getLastReboot(c *utils.OsQuery) {
	var ret []string

	out := execToString(exec.Command("last", "reboot"))
	lines := strings.Split(out, "\n")

	r, _ := regexp.Compile("[a-zA-Z]{3}\\s{1}[a-zA-Z]{3}\\s{1,2}\\d{1,2}\\s{1}\\d{2}\\:\\d{2}")
	for _, line := range lines {
		findStr := r.FindString(line)
		fmt.Println(findStr)
		if findStr != "" {
			ret = append(ret, findStr)
		}
	}
	s.LastReboot = ret
}

func (s *System) getUptime(c *utils.OsQuery) {
	resp, err := c.Client.Query("select total_seconds from uptime")
	if err != nil {
		log.Fatalf("Error while gathering SysInfo: %v\n", err)
	}
	conv, _ := strconv.Atoi(resp.Response[0]["total_seconds"])
	s.Uptime = uint(conv)
}

func (s *System) getSystemInfo(c *utils.OsQuery) {
	resp, err := c.Client.Query("select uuid, hostname from system_info;")
	if err != nil {
		log.Fatalf("error while trying to get the system_info: %v\n", err)
	}
	s.Hostname = resp.Response[0]["hostname"]
	s.UUID = resp.Response[0]["uuid"]
}

func (s *System) getUsbDevices(c *utils.OsQuery) {
	resp, err := c.Client.Query("select vendor, model from usb_devices;")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, usb := range resp.Response {
		s.Usb = append(s.Usb, Usb{vendor: usb["vendor"], model: usb["model"]})
	}
}
