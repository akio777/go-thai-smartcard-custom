package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"go-thai-smartcard-custom/pkg/logger"
	"go-thai-smartcard-custom/pkg/model"
	"go-thai-smartcard-custom/pkg/server"
	"go-thai-smartcard-custom/pkg/smc"
	"go-thai-smartcard-custom/pkg/util"
)

func main() {
	// load env
	logger.InitLogger()
	port := util.GetEnv("SMC_AGENT_PORT", "9898")
	showImage := util.GetEnvBool("SMC_SHOW_IMAGE", true)
	showLaser := util.GetEnvBool("SMC_SHOW_LASER", true)
	showNhso := util.GetEnvBool("SMC_SHOW_NHSO", false)

	broadcast := make(chan model.Message)

	serverCfg := server.ServerConfig{
		Broadcast: broadcast,
		Port:      port,
	}
	go server.Serve(serverCfg)

	opts := &smc.Options{
		ShowFaceImage: showImage,
		ShowNhsoData:  showNhso,
		ShowLaserData: showLaser,
	}

	go func() {
		smc := smc.NewSmartCard()
		for {
			err := smc.StartDaemon(broadcast, opts)
			if err != nil {
				logger.LOGGER().Printf("Error occurred in daemon process (%v), wait 2 seconds to retry or press Ctrl+C to exit.", err.Error())

				message := model.Message{
					Event: "smc-error",
					Payload: map[string]string{
						"message": fmt.Sprintf("Error occurred in daemon process, %v.", err.Error()),
					},
				}
				broadcast <- message

				time.Sleep(2 * time.Second)
			}
		}
	}()

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	s := <-sig
	logger.LOGGER().Printf("Received %v signal to shutdown.", s)

}
