package lights

import (
	"context"
	"fmt"

	"github.com/amimof/huego"
	"github.com/pkg/errors"
)

// Hue is a struct that maintains hue connections and info
type Hue struct {
	Bridge *huego.Bridge
}

// NewHue creates an instance of a hue controller
func NewHue(ctx context.Context, ip, username string) (*Hue, error) {
	bridge := huego.New(ip, username)

	return &Hue{bridge}, nil
}

// TurnOffEverything finds all lights on the bridge and shuts them off
func (hue *Hue) TurnOffEverything(ctx context.Context) []error {
	lights, err := hue.Bridge.GetLightsContext(ctx)
	if err != nil {
		return []error{errors.WithMessage(err, "failed to get lights")}
	}

	offErrors := []error{}

	for _, light := range lights {
		err := light.OffContext(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("encountered error while shutting off light (%+v)", light.Name)
			offErrors = append(offErrors, errors.WithMessage(err, errMsg))
		}
	}

	return offErrors
}

// TurnOnEverything finds all lights on the bridge and turns them on
func (hue *Hue) TurnOnEverything(ctx context.Context) []error {
	lights, err := hue.Bridge.GetLightsContext(ctx)
	if err != nil {
		return []error{errors.WithMessage(err, "failed to get lights")}
	}

	onErrors := []error{}

	for _, light := range lights {
		err := light.OnContext(ctx)
		if err != nil {
			errMsg := fmt.Sprintf("encountered error while turning on light (%+v)", light.Name)
			onErrors = append(onErrors, errors.WithMessage(err, errMsg))
		}
	}

	return onErrors
}
