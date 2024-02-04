package logger

import (
	"fmt"
	"log"
	"runtime"
)

type sLogLine struct {
	lvl int
	msg string
}

const (
	LVL_ERROR = 1
	LVL_WARN  = 2
	LVL_INFO  = 3
	LVL_DEBUG = 4
)

var (
	Vbs  = LVL_ERROR
	done = make(chan struct{})
	logs = make(chan sLogLine, 1000)
)

func Start(v int) {
	Vbs = v
	go monitorLoop()
}

func Stop() {
	close(logs)
	<-done
}

func Log(l int, m string) {
	var c string
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		details := runtime.FuncForPC(pc)
		c = details.Name()
	} else {
		c = "unknown"
	}
	if l <= Vbs {
		if Vbs == LVL_DEBUG {
			m = fmt.Sprintf("[%s] %s", c, m)
		}
		logs <- sLogLine{lvl: l, msg: m}
	}
}

func monitorLoop() {
	for f := range logs {
		log.Print(f.msg)
	}
	close(done)
}
