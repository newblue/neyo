package neyo

// 简化的Log调用

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	ERROR = 0
	WARN  = iota + 1
	INFO
	DEBUG
)

var (
	debug = flag.Bool("debug", false, "Enable Debug, verbose output")
)

func Log(level int, format string, v ...interface{}) {
	var prefix string
	switch level {
	case ERROR:
		prefix = "ERROR"
	case INFO:
		prefix = "INFO"
	case WARN:
		prefix = "WARN"
	case DEBUG:
		prefix = "DEBUG"
	default:
		prefix = "UNKOWN"
	}

	if *debug == true && level == DEBUG {
		_, file, line, ok := runtime.Caller(1)

		if ok {
			file = filepath.Base(file)
			fmt.Printf("[%6s] %s %s(%d): %s\n", prefix, time.Now().Format("15:04:05"), file, line, fmt.Sprintf(format, v...))
		} else {
			fmt.Printf("[%6s] %s >> %s\n", prefix, time.Now().String(), fmt.Sprintf(format, v...))
		}
	} else if level != DEBUG {
		fmt.Printf("[%6s] %s\n", prefix, fmt.Sprintf(format, v...))
	}
	if level == ERROR {
		os.Exit(-1)
	}
}
