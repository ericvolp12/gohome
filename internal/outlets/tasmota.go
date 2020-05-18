package outlets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// TasmotaDevice contains information about a single Tasmota hardware device
type TasmotaDevice struct {
	Host string
	Name string
}

// Tasmota contains an instance of the Tasmota devices on the network
type Tasmota struct {
	Devices []*TasmotaDevice
}

// NewTasmota is a factory for a Tasmota interface
func NewTasmota(ctx context.Context, hosts []string, names []string) (*Tasmota, error) {
	devices := []*TasmotaDevice{}
	for i, host := range hosts {
		if len(names)-1 >= i {
			devices = append(devices, &TasmotaDevice{Host: host, Name: names[i]})
		}
	}

	return &Tasmota{Devices: devices}, nil
}

// TurnOnEverything turns on all devices the Tasmota knows about
func (t *Tasmota) TurnOnEverything(ctx context.Context) []error {
	var errs []error

	for _, device := range t.Devices {
		url := fmt.Sprintf("http://%v/cm?cmnd=Power%%20On", device.Host)
		resp, err := http.Get(url)
		if err != nil {
			errs = append(errs, errors.WithMessage(err, fmt.Sprintf("failed to turn on device (%v)", device.Name)))
		}

		if resp.StatusCode != http.StatusOK {
			errs = append(errs, fmt.Errorf("failed to turn on device (%v), status: (%v)", device.Name, resp.StatusCode))
		}
	}

	return errs
}

// TurnOffEverything turns off all devices the Tasmota knows about
func (t *Tasmota) TurnOffEverything(ctx context.Context) []error {
	var errs []error

	for _, device := range t.Devices {
		url := fmt.Sprintf("http://%v/cm?cmnd=Power%%20Off", device.Host)
		resp, err := http.Get(url)
		if err != nil {
			errs = append(errs, errors.WithMessage(err, fmt.Sprintf("failed to turn off device (%v)", device.Name)))
		}

		if resp.StatusCode != http.StatusOK {
			errs = append(errs, fmt.Errorf("failed to turn off device (%v), status: (%v)", device.Name, resp.StatusCode))
		}
	}

	return errs
}
