package monitor

import (
	"fmt"
	"net/http"
	"strings"
)

const webHookToken = ""

func MonitorLogStream() error {
	logs := newLogTail()
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered ", r)
				}
			}()

			if err := logs.StartGathering(); err != nil {
				fmt.Println("Log gathering failed with error: %w", err)
				return
			}

			for log := range logs.Channel {
				if strings.Contains(log.EventMessage, "Cameras changed") {
					switch log.EventMessage {
					case "Cameras changed to []":
						fmt.Printf("%s: Camera disconnected\n", log.Timestamp)
						http.Get(fmt.Sprintf("https://mkzense.com/webhook/alexa/%s/CallEnd", webHookToken))
					default:
						fmt.Printf("%s: Camera connected\n", log.Timestamp)
						http.Get(fmt.Sprintf("https://mkzense.com/webhook/alexa/%s/CallStart", webHookToken))
					}
				}
			}
		}()
	}
}
