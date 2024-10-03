package utils

import (
	"fmt"
)

type Logger struct{}

// TODO: Later on update the implementation to generate log file
func (l *Logger) Info(msg ...any) {
	fmt.Println(msg)
}
func (l *Logger) Error(msg ...any) {
	fmt.Println(msg)
}
