package neyo

import (
	"flag"
	"fmt"
	//"github.com/wsxiaoys/terminal/color"
	"os"
	"path/filepath"
	"runtime"
)

const (
	ERROR = 0
	WARN  = iota + 1
	INFO
	DEBUG
)

var (
	debug = flag.Bool("debug", false, "Enable Debug, verbose output")
	color = flag.Bool("disable-color", true, "Disable color output.")
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

	if *debug == true {
		_, file, line, ok := runtime.Caller(1)

		if ok {
			file = filepath.Base(file)
			fmt.Printf("[%6s] %s(%d): %s\n", prefix, file, line, fmt.Sprintf(format, v...))
		}
	} else if level < DEBUG {
		fmt.Printf("[%6s] %s\n", prefix, fmt.Sprintf(format, v...))
	}
	if level == ERROR {
		os.Exit(-1)
	}
}
