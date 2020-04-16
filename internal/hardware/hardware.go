package hardware

import (
	"fmt"
	"log"

	"github.com/jgengo/macinfo-client/internal/utils"
)

// Hardware is storing the Mac hardware information
type Hardware struct {
	Hostname string `json:"hostname"`
	Usb      []Usb  `json:"usb_devices"`
}

// Usb is storing a usb devices
type Usb struct {
	vendor string
	model  string
}

// GetInfo retrieves all the information of the client
func GetInfo(c *utils.OsQuery) {
	var h Hardware
	h.getHostname(c)
	h.getUsbDevices(c)

	fmt.Printf("%+v", h)
}

func (h *Hardware) getHostname(c *utils.OsQuery) {
	resp, err := c.Client.Query("select hostname from system_info;")
	if err != nil {
		log.Fatalf("error while trying to get the hostname")
	}
	h.Hostname = resp.Response[0]["hostname"]
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
