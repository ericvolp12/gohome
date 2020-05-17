package outlets

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bcurren/go-ssdp"
	"github.com/pkg/errors"
)

// Wemo is a wemo interface to talk to wemo devices
type Wemo struct {
	Devices []*Device
}

// Device tracks a wemo Device
type Device struct {
	Host string
}

// NewWemo returns an instance of a Wemo interface
func NewWemo(ctx context.Context) (*Wemo, error) {

	// From https://github.com/savaki/go.wemo/blob/faafcb35be48b0c210213a2483a07fb15041df69/discover.go#L29
	urns := []string{
		"urn:Belkin:device:controllee:1",
	}

	devices := []*Device{}

	for _, urn := range urns {
		responses, err := ssdp.Search(urn, 3*time.Second)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to search for SSDP clients")
		}

		for _, response := range responses {
			devices = append(devices, &Device{response.Location.String()})
		}
	}

	fmt.Printf("Found %v devices\n", len(devices))

	return &Wemo{Devices: devices}, nil
}

// Blob from here - https://github.com/savaki/go.wemo/blob/faafcb35be48b0c210213a2483a07fb15041df69/messages.go#L52
func newSetBinaryStateMessage(on bool) string {
	value := 0
	if on {
		value = 1
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
						<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
						<s:Body>
							<u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1">
							<BinaryState>%v</BinaryState>
							</u:SetBinaryState>
						</s:Body>
						</s:Envelope>`, value)
}

// TurnOff turns off a device
func (d *Device) TurnOff(ctx context.Context) error {

	url := fmt.Sprintf("http://%v/upnp/control/basicevent1", d.Host)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader([]byte(newSetBinaryStateMessage(false))))
	if err != nil {
		return errors.WithMessage(err, "failed to initialize request")
	}

	req.Header.Add("SOAPACTION", "urn:Belkin:service:basicevent:1#SetBinaryState")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, "failed to execute POST request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed to POST to wemo device, status:%+v", resp.StatusCode))
	}

	return nil
}

// TurnOffEverything turns off all wemo devices (potentially dangerous?)
func (w *Wemo) TurnOffEverything(ctx context.Context) []error {
	offErrors := []error{}

	for _, device := range w.Devices {
		err := device.TurnOff(ctx)
		if err != nil {
			offErrors = append(offErrors, errors.WithMessage(err, fmt.Sprintf("failed to turn off device (%+v)", device.Host)))
		}
	}

	return offErrors
}
