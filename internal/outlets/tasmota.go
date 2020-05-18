package outlets

import (
	"context"
	"crypto/tls"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

// TasmotaMQTT is a struct for a TasmotaMQTT client
type TasmotaMQTT struct {
	Client MQTT.Client
	Topic  string
}

// NewTasmotaMQTT is a factory for a TasmotaMQTT client
func NewTasmotaMQTT(ctx context.Context, server, topic string) (*TasmotaMQTT, error) {
	client, err := acquireMQTTClient(ctx, server)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to acquire TasmotaMQTT client")
	}
	return &TasmotaMQTT{Client: client, Topic: topic}, nil
}

func acquireMQTTClient(ctx context.Context, server string) (MQTT.Client, error) {
	clientid := "gohome_" + strconv.Itoa(time.Now().Second())

	connOpts := MQTT.NewClientOptions().AddBroker(server).SetClientID(clientid).SetCleanSession(true)

	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)

	attempts := 1
	var err error

	for attempts < 10 {
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			err = token.Error()
		} else {
			return client, nil
		}
		time.Sleep(time.Second * time.Duration(attempts))
		attempts++
	}

	return nil, errors.WithMessage(err, "failed to initialize MQTT client")
}

// TurnOnEverything turns on all devices the Tasmota knows about
func (t *TasmotaMQTT) TurnOnEverything(ctx context.Context) error {
	if token := t.Client.Publish(t.Topic, 0, false, "on"); token.Wait() && token.Error() != nil {
		return errors.WithMessage(token.Error(), "failed to send power on MQTT message")
	}
	return nil
}

// TurnOffEverything turns off all devices the Tasmota knows about
func (t *TasmotaMQTT) TurnOffEverything(ctx context.Context) error {
	if token := t.Client.Publish(t.Topic, 0, false, "off"); token.Wait() && token.Error() != nil {
		return errors.WithMessage(token.Error(), "failed to send power off MQTT message")
	}
	return nil
}
