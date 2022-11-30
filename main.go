package main

import (
	"fmt"

	"github.com/owenrumney/webcam-watcher/pkg/monitor"
)

func main() {

	fmt.Println("Starting webcam watcher")

	if err := monitor.MonitorLogStream(); err != nil {
		fmt.Println(err)
	}

}
