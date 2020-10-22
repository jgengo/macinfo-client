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
func GetInfo() (*System, error) {
	var system System

	system.Token = utils.Cfg.APIToken

	if err := system.getActiveUser(); err != nil {
		return nil, err
	}
	if err := system.getLastReboot(); err != nil {
		return nil, err
	}
	if err := system.getSystemInfo(); err != nil {
		return nil, err
	}
	if err := system.getUsbDevices(); err != nil {
		return nil, err
	}
	if err := system.getUptime(); err != nil {
		return nil, err
	}
	if err := system.getTemperatureSensors(); err != nil {
		return nil, err
	}
	if err := system.getOsVersion(); err != nil {
		return nil, err
	}

	return &system, nil
}

func execToString(cmd *exec.Cmd) (string, error) {
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return b.String(), nil
}

func (s *System) getActiveUser() error {
	resp, err := utils.OsQ.Client.Query("select user from logged_in_users where type='user' and tty='console' limit 1;")
	if err != nil {
		log.Printf("failed to get ActiveUser: %v\n", err)
		return err
	}
	if len(resp.Response) > 0 {
		s.ActiveUser = resp.Response[0]["user"]
	} else {
		s.ActiveUser = ""
	}
	return nil
}

// You may know that the command last reboot will get slower
// if you don't sometimes clean your asl log files.
func (s *System) getLastReboot() error {
	out, err := execToString(exec.Command("last", "reboot"))
	if err != nil {
		log.Printf("failed to get the last reboots: %v\n", err)
		return err
	}
	r, _ := regexp.Compile("[a-zA-Z]{3}\\s{1}[a-zA-Z]{3}\\s{1,2}\\d{1,2}\\s{1}\\d{2}\\:\\d{2}") // catches "Sun Apr 19 11:41"
	ret := r.FindAllString(out, -1)
	s.LastReboot = ret
	return nil
}

func (s *System) getOsVersion() error {
	resp, err := utils.OsQ.Client.Query("select version, build from os_version;")
	if err != nil {
		log.Printf("failed to get os info: %v\n", err)
		return err
	}
	s.OsVersion = OsVersion{Version: resp.Response[0]["version"], Build: resp.Response[0]["build"]}
	return nil
}

func (s *System) getTemperatureSensors() error {
	resp, err := utils.OsQ.Client.Query("select name, celsius from temperature_sensors;")
	if err != nil {
		log.Printf("failed to get Sensors: %v\n", err)
		return err
	}

	for _, sensor := range resp.Response {
		celsius, err := strconv.ParseFloat(sensor["celsius"], 1)
		if err != nil {
			log.Printf("failed to convert celsius: %s\n", sensor["celsius"])
			return err
		}
		s.Sensors = append(s.Sensors, Sensor{Name: sensor["name"], Celsius: celsius})
	}
	return nil
}

func (s *System) getUptime() error {
	resp, err := utils.OsQ.Client.Query("select total_seconds from uptime")
	if err != nil {
		log.Printf("failed to get uptime: %v\n", err)
		return err
	}
	conv, err := strconv.Atoi(resp.Response[0]["total_seconds"])
	if err != nil {
		log.Printf("failed to convert uptime: %s\n", resp.Response[0]["total_seconds"])
		return err
	}
	s.Uptime = uint(conv)
	return nil
}

func (s *System) getSystemInfo() error {
	resp, err := utils.OsQ.Client.Query("select uuid, hostname from system_info;")
	if err != nil {
		log.Printf("failed to get the system_info: %v\n", err)
		return err
	}
	s.Hostname = resp.Response[0]["hostname"]
	s.UUID = resp.Response[0]["uuid"]
	return nil
}

func (s *System) getUsbDevices() error {
	resp, err := utils.OsQ.Client.Query("select vendor, model from usb_devices;")
	if err != nil {
		log.Printf("failed to get usbDevices: %v\n", err)
		return err
	}

	for _, usb := range resp.Response {
		s.Usb = append(s.Usb, Usb{
			Vendor: strings.TrimSpace(usb["vendor"]),
			Model:  strings.TrimSpace(usb["model"]),
		})
	}

	return nil
}
