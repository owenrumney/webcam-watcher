package monitor

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"sync"
)

// implementation generally lifted from https://github.com/MaxSchaefer/macos-log-stream/blob/main/pkg/mls/logs.go

type Log struct {
	TraceID      int64  `json:"traceID"`
	EventMessage string `json:"eventMessage"`
	Timestamp    string `json:"timestamp"`
}

type Logs struct {
	m       sync.Mutex
	Channel chan Log
	exit    chan bool
}

func newLogTail() *Logs {
	return &Logs{
		Channel: make(chan Log),
		exit:    make(chan bool),
	}
}

func (logs *Logs) StartGathering() error {
	args := []string{
		"stream",
		"--color=none",
		"--style=ndjson",
	}

	cmd := exec.Command("log", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go func() {
		logs.m.Lock()
		defer logs.m.Unlock()

		cmd.Start()
		defer cmd.Process.Kill()

		// drop first message
		bufio.NewReader(stdout).ReadLine()

		dec := json.NewDecoder(stdout)

		for {
			select {
			case <-logs.exit:
				return
			default:
				log := Log{}
				if err := dec.Decode(&log); err != nil {
					continue
				} else {
					logs.Channel <- log
				}
			}
		}
	}()

	go func() {
		stderrBuf := bufio.NewReader(stderr)
		for {
			line, _, _ := stderrBuf.ReadLine()
			if len(line) > 0 {
				logs.StopGathering()
				continue
			}
		}
	}()

	return nil
}

func (logs *Logs) StopGathering() {
	logs.exit <- true
}
