package logger

import (
	"fmt"
	"strconv"
	"time"
)

type Logger struct {
	Channel   chan string
	Responses int
	State     int
	Ticker    *time.Ticker
	Running   bool
	Counter   int
}

func (log *Logger) Initialize(size int) {
	log.Channel = make(chan string, size)
	log.Responses = 0
	log.Ticker = time.NewTicker(time.Second / 2.0)
	log.Running = true
	log.Counter = 0
}

func (log *Logger) Increment() {
	log.Responses += 1
}

func (log *Logger) Queue(s string) {
	log.Channel <- s
}

func (log *Logger) Log() {
	d := [4]string{"/", "-", "\\", "|"}

	for range log.Ticker.C {
		if len(log.Channel) > 0 {
			fmt.Println(<-log.Channel)
		} else if log.Running {
			out := d[log.Counter%4] + " " + strconv.Itoa(log.Responses) + " responses in " + strconv.FormatFloat(float64(log.Counter)*0.5, 'f', 1, 64) + " seconds.\r"
			fmt.Print(out)
			log.Counter += 1
		}
	}
}

func (log *Logger) Kill() {
	fmt.Print("\n")
	log.Running = false
}
