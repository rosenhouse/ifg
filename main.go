package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rosenhouse/ifg/application"
)

func main() {
	config := application.Config{
		RootPath:     os.ExpandEnv(LoadOrFail("ROOT_PATH")),
		Port:         LoadOrFail("PORT"),
		VCAPServices: LoadOrFail("VCAP_SERVICES"),
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
	time.Sleep(3 * time.Second)
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
