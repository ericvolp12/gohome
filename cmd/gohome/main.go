package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ericvolp12/gohome/internal/lights"
	"github.com/ericvolp12/gohome/internal/outlets"
	"github.com/gin-gonic/gin"
)

// Application stores the mid-level interfaces for go_home
type Application struct {
	Hue     *lights.Hue
	Wemo    *outlets.Wemo
	Tasmota *outlets.TasmotaMQTT
	APIKey  string
}

// Request is a valid POST request to this API
type Request struct {
	APIKey string `json:"apiKey"`
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

	r.POST("/on", app.OnHandler)

	r.POST("/off", app.OffHandler)

	r.Run() // listen and serve on 0.0.0.0:8053
}

// OnHandler handles turning things on
func (app *Application) OnHandler(c *gin.Context) {
	var req Request

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"errors": "bad request, missing API key",
		})
		return
	}

	if req.APIKey != app.APIKey {
		c.JSON(401, gin.H{
			"errors": "invalid API key",
		})
		return
	}

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

	err = app.Tasmota.TurnOnEverything(c.Request.Context())
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

// OffHandler handles turning things off
func (app *Application) OffHandler(c *gin.Context) {
	var req Request

	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"errors": "bad request, missing API key",
		})
		return
	}

	if req.APIKey != app.APIKey {
		c.JSON(401, gin.H{
			"errors": "invalid API key",
		})
		return
	}

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

	err = app.Tasmota.TurnOffEverything(c.Request.Context())
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
