package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
)

func GetLogger(outPath string) (*log.Logger, error) {
	file, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	prefix := fmt.Sprintf("[PID=%v] ", strconv.Itoa(os.Getpid()))

	logger := log.New(file, prefix, log.Ldate|log.Lmicroseconds)

	return logger, nil
}
