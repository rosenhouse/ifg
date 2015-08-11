package main

import (
	"fmt"
	"os"

	"github.com/rosenhouse/ifg/application"
)

func main() {
	config := application.Config{
		RootPath: LoadOrFail("ROOT_PATH"),
		Port:     LoadOrFail("PORT"),
	}

	app, err := application.NewApplication(config)
	if err != nil {
		Fail(err.Error())
	}

	err = app.Boot()
	if err != nil {
		Fail(err.Error())
	}
}

func LoadOrFail(variable string) string {
	value := os.Getenv(variable)
	if value == "" {
		Fail("'%s' is a required environment variable", variable)
	}

	return value
}

func Fail(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
