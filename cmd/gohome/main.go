package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ericvolp12/gohome/internal/lights"
	"github.com/ericvolp12/gohome/internal/outlets"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Application stores the mid-level interfaces for go_home
type Application struct {
	Hue     *lights.Hue
	Wemo    *outlets.Wemo
	Tasmota *outlets.TasmotaMQTT
	APIKey  string
}

// AuthenticatedRequest is a valid, authenticated POST request to this API
type AuthenticatedRequest struct {
	APIKey string `json:"apiKey"`
}

// PowerStateRequest is a valid POST request to this API to set a device state
type PowerStateRequest struct {
	Device     string `json:"device"`
	PowerState string `json:"powerState"`
}

func main() {
	fmt.Println("Go HOME")

	ctx := context.Background()

	hueGatewayIP := os.Getenv("HUE_GATEWAY_IP")
	hueUsername := os.Getenv("HUE_USERNAME")
	apiKey := os.Getenv("GOHOME_API_KEY")
	mqttServer := "tcp://" + os.Getenv("GOHOME_MQTT_SERVER")
	mqttTopic := os.Getenv("GOHOME_MQTT_TOPIC")

	hue, err := lights.NewHue(ctx, hueGatewayIP, hueUsername)
	if err != nil {
		log.Panic("failed to initialize HUE bridge: ", err)
	}
	log.Printf("Hue Initialized successfully")

	wemo, err := outlets.NewWemo(ctx)
	if err != nil {
		log.Panic("failed to initialize Wemo connections: ", err)
	}

	log.Printf("Wemo initialized with %v device(s): \n", len(wemo.Devices))
	for _, device := range wemo.Devices {
		log.Printf("\tDevice: (%+v)\n", device.Host)
	}

	tasmota, err := outlets.NewTasmotaMQTT(ctx, mqttServer, mqttTopic)
	if err != nil {
		log.Panic("failed to initialize TasmotaMQTT: ", err)
	}
	log.Printf("TasmotaMQTT Initialized successfully")

	app := Application{Hue: hue, Wemo: wemo, APIKey: apiKey, Tasmota: tasmota}

	r := gin.Default()

	r.Use(app.Auth())

	r.POST("/on", app.OnHandler)

	r.POST("/off", app.OffHandler)

	r.POST("/setState", app.SetStateHandler)

	r.Run() // listen and serve on 0.0.0.0:8053
}

// Auth is a simple authentication middleware for gohome
func (app *Application) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthenticatedRequest
		err := c.ShouldBindBodyWith(&req, binding.JSON)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{
				"errors": "bad request, missing apiKey",
			})
			return
		}

		if req.APIKey != app.APIKey {
			c.AbortWithStatusJSON(401, gin.H{
				"errors": "invalid apiKey",
			})
			return
		}
		c.Next()
		return
	}
}

// OnHandler handles turning things on
func (app *Application) OnHandler(c *gin.Context) {
	errors := app.Hue.TurnOnEverything(c.Request.Context())
	if len(errors) != 0 {
		c.JSON(500, gin.H{
			"errors": errors,
		})

		return
	}

	errors = app.Wemo.TurnOnEverything(c.Request.Context())
	if len(errors) != 0 {
		log.Printf("%+v", errors)
		c.JSON(500, gin.H{
			"errors": errors,
		})

		return
	}

	err := app.Tasmota.TurnOnEverything(c.Request.Context())
	if err != nil {
		log.Printf("%+v", err)
		c.JSON(500, gin.H{
			"error": err,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "all lights turned on successfully",
	})
}

// SetStateHandler sets the state of a specific MQTT device
func (app *Application) SetStateHandler(c *gin.Context) {
	var req PowerStateRequest
	err := c.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil || req.Device == "" || req.PowerState == "" {
		c.JSON(400, gin.H{
			"errorMessage": "bad request, make sure to send device, and powerState in the body",
		})
		return
	}

	err = app.Tasmota.SetDevicePowerState(c.Request.Context(), req.Device, req.PowerState)
	if err != nil {
		log.Printf("%+v", err)
		c.JSON(500, gin.H{
			"error": err,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("successfully set device (%s) power state to (%s)", req.Device, req.PowerState),
	})
}

// OffHandler handles turning things off
func (app *Application) OffHandler(c *gin.Context) {
	errors := app.Hue.TurnOffEverything(c.Request.Context())
	if len(errors) != 0 {
		c.JSON(500, gin.H{
			"errors": errors,
		})

		return
	}

	errors = app.Wemo.TurnOffEverything(c.Request.Context())
	if len(errors) != 0 {
		log.Printf("%+v", errors)
		c.JSON(500, gin.H{
			"errors": errors,
		})

		return
	}

	err := app.Tasmota.TurnOffEverything(c.Request.Context())
	if err != nil {
		log.Printf("%+v", err)
		c.JSON(500, gin.H{
			"error": err,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "all lights turned off successfully",
	})
}
