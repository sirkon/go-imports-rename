package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func newLogger() *zerolog.Logger {
	writer := zerolog.NewConsoleWriter()
	writer.TimeFormat = time.RFC3339
	writer.FormatCaller = func(i interface{}) string {
		if i == nil {
			return ""
		}
		value := fmt.Sprintf("%v", i)

		pos := strings.Index(value, "github.com")
		if pos < 0 {
			pos = strings.Index(value, "gitlab.")
			if pos >= 0 {
				value = value[pos:]
			}
		} else {
			value = value[pos:]
		}
		return "(" + value + ")"
	}
	writer.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\033[1m%v\033[0m", i)
	}
	writer.FormatTimestamp = func(i interface{}) string {
		if i == nil {
			return ""
		}
		return fmt.Sprintf("\033[33;1m%v\033[0m", i)
	}
	writer.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\033[35m%s\033[0m", i)
	}
	writer.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("[%s]", i)
	}
	writer.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("\033[31m%s\033[0m", i)
	}
	writer.FormatErrFieldValue =
		func(i interface{}) string {
			return fmt.Sprintf("\033[31m[%v]\033[0m", i)
		}
	log := zerolog.New(writer).Level(zerolog.DebugLevel)
	return &log
}
