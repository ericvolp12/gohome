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

func main() {
	fmt.Println("Go HOME")

	ctx := context.Background()

	hueGatewayIP := os.Getenv("HUE_GATEWAY_IP")

	hueUsername := os.Getenv("HUE_USERNAME")

	hue, err := lights.NewHue(ctx, hueGatewayIP, hueUsername)
	if err != nil {
		log.Panic("failed to initialize HUE bridge: ", err)
	}

	wemo, err := outlets.NewWemo(ctx)
	if err != nil {
		log.Panic("failed to initialize Wemo connections: ", err)
	}

	r := gin.Default()
	r.GET("/on", func(c *gin.Context) {
		errors := hue.TurnOnEverything(c.Request.Context())
		if len(errors) != 0 {
			c.JSON(500, gin.H{
				"errors": errors,
			})

			return
		}
		c.JSON(200, gin.H{
			"message": "all lights turned on successfully",
		})
	})

	r.GET("/off", func(c *gin.Context) {
		errors := hue.TurnOffEverything(c.Request.Context())
		if len(errors) != 0 {
			c.JSON(500, gin.H{
				"errors": errors,
			})

			return
		}

		errors = wemo.TurnOffEverything(ctx)
		if len(errors) != 0 {
			c.JSON(500, gin.H{
				"errors": errors,
			})

			return
		}

		c.JSON(200, gin.H{
			"message": "all lights turned off successfully",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
