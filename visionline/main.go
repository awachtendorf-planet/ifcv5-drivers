package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weareplanet/ifcv5-main/app"
	"github.com/weareplanet/ifcv5-main/config"

	driver "github.com/weareplanet/ifcv5-drivers/visionline"
)

var (
	configFile  = ""
	serviceName = "visionline"
)

func main() {

	flag.StringVar(&configFile, "c", "", fmt.Sprintf("use config file (default \"%s.toml\")", serviceName))
	flag.StringVar(&serviceName, "n", serviceName, "use service name")

	flag.Parse()

	if app.HasVersionCommand() {
		fmt.Println(config.VersionString())
		os.Exit(0)
	}

	ifc, err := driver.New(configFile)

	if !app.HasServiceCommand() {
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}

	app.Service(ifc, serviceName, configFile)

}
