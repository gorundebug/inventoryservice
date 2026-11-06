package main

import (
	"context"
	"flag"
	"log"
	"os"

	inventoryservice "github.com/gorundebug/inventoryservice/internal/app"
)

func main() {
	valuesPathArg := flag.String("values", "./config/overrides.yaml", "service config values path")
	configPathArg := flag.String("config", "./config/config.yaml", "service config path")
	flag.Parse()

	if err := inventoryservice.Start(
		context.Background(),
		nil,
		nil,
		configPathArg,
		valuesPathArg,
		os.Getenv("APP_VERSION"),
		os.Getenv("APP_COMMIT"),
	); err != nil {
		log.Printf("start service error: %s", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
