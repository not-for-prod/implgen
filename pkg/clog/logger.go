package clog

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Fatal(msg ...interface{}) {
	fmt.Println(append([]interface{}{color.RedString("Fatal error:")}, msg...)...)
	os.Exit(1)
}

func Error(msg ...interface{}) {
	fmt.Println(append([]interface{}{color.RedString("Error:")}, msg...)...)
}

func Warn(msg ...interface{}) {
	fmt.Println(append([]interface{}{color.YellowString("Warning:")}, msg...)...)
}

func Info(msg ...interface{}) {
	fmt.Println(append([]interface{}{color.GreenString("Info:")}, msg...)...)
}

func Fatalf(f string, v ...interface{}) {
	fmt.Printf(color.RedString("Fatal error: ")+f+"\n", v...)
	os.Exit(1)
}

func Errorf(f string, v ...interface{}) {
	fmt.Printf(color.RedString("Error: ")+f+"\n", v...)
}

func Warnf(f string, v ...interface{}) {
	fmt.Printf(color.YellowString("Warning: ")+f+"\n", v...)
}

func Infof(f string, v ...interface{}) {
	fmt.Printf(color.GreenString("Info: ")+f+"\n", v...)
}
