package logs

import (
	"fmt"
	"log"
	"runtime"
)

func Error(format string, a ...any) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d [ERROR] %s", file, line, fmt.Sprintf(format, a...))
}

func Debug(format string, a ...any) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d [DEBUG] %s", file, line, fmt.Sprintf(format, a...))
}

func Info(format string, a ...any) {
	_, file, line, _ := runtime.Caller(1)
	log.Printf("%s:%d [INFO] %s", file, line, fmt.Sprintf(format, a...))
}
