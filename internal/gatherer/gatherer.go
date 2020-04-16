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

// Hardware stores the Mac hardware information
type Hardware struct {
	Hostname string `json:"hostname"`
	UUID     string `json:"uuid"`
	Usb      []Usb  `json:"usb_devices"`
}

// Usb stores a usb devices
type Usb struct {
	vendor string
	model  string
}

// System stores system info
type System struct {
	Uptime     uint     `json:"uptime"`
	LastReboot []string `json:"last_reboot"`
}

// GetInfo retrieves all the information of the client
func GetInfo(c *utils.OsQuery) {
	var hardware Hardware
	hardware.getSystemInfo(c)
	hardware.getUsbDevices(c)

	var system System
	system.getSysInfo(c)
	system.getLastReboot(c)
	// fmt.Printf("%+v", hardware)
	fmt.Printf("%+v", system)
	for _, p := range system.LastReboot {
		fmt.Println(p)
	}

}

func shellExec(cmd *exec.Cmd) string {
	var b bytes.Buffer
	cmd.Stdout = &b
	if err := cmd.Run(); err != nil {
		log.Fatalf("error while gathering last reboots: %v\n", err)
	}
	return b.String()
}

// TODO: refacto
func (s *System) getLastReboot(c *utils.OsQuery) {
	var ret []string

	out := shellExec(exec.Command("last", "reboot"))
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

// TODO: (X1) voir le nom
func (s *System) getSysInfo(c *utils.OsQuery) {
	resp, err := c.Client.Query("select total_seconds from uptime")
	if err != nil {
		log.Fatalf("Error while gathering SysInfo: %v\n", err)
	}
	conv, err := strconv.Atoi(resp.Response[0]["total_seconds"])
	if err != nil {
		log.Printf("error while converting uptime.total_seconds: %v", err)
		s.Uptime = 0
	} else {
		s.Uptime = uint(conv)
	}
}

// TODO: (X1) voir le nom
func (h *Hardware) getSystemInfo(c *utils.OsQuery) {
	resp, err := c.Client.Query("select uuid, hostname from system_info;")
	if err != nil {
		log.Fatalf("error while trying to get the system_info: %v\n", err)
	}
	h.Hostname = resp.Response[0]["hostname"]
	h.UUID = resp.Response[0]["uuid"]
}

func (h *Hardware) getUsbDevices(c *utils.OsQuery) {
	resp, err := c.Client.Query("select vendor, model from usb_devices;")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, usb := range resp.Response {
		h.Usb = append(h.Usb, Usb{vendor: usb["vendor"], model: usb["model"]})
	}
}
