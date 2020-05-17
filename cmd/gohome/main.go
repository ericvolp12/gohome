package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ericvolp12/go_home/internal/lights"
	"github.com/ericvolp12/go_home/internal/outlets"
	"github.com/gin-gonic/gin"
)

// Application stores the mid-level interfaces for go_home
type Application struct {
	Hue    *lights.Hue
	Wemo   *outlets.Wemo
	APIKey string
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

	hue, err := lights.NewHue(ctx, hueGatewayIP, hueUsername)
	if err != nil {
		log.Panic("failed to initialize HUE bridge: ", err)
	}

	wemo, err := outlets.NewWemo(ctx)
	if err != nil {
		log.Panic("failed to initialize Wemo connections: ", err)
	}

	app := Application{Hue: hue, Wemo: wemo, APIKey: apiKey}

	r := gin.Default()

	r.POST("/on", app.OnHandler)

	r.POST("/off", app.OffHandler)

	r.Run() // listen and serve on 0.0.0.0:8080
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
	log.Printf("%+v", errors)

	if len(errors) != 0 {
		c.JSON(500, gin.H{
			"errors": errors,
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
	log.Printf("%+v", errors)

	if len(errors) != 0 {
		c.JSON(500, gin.H{
			"errors": errors,
		})

		return
	}

	c.JSON(200, gin.H{
		"message": "all lights turned off successfully",
	})
}
