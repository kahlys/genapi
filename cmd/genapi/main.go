package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/kahlys/genapi"
	"github.com/spf13/viper"
)

// flags
var (
	version  = "undefined"
	fversion = flag.Bool("version", false, "print server version")
	fconfig  = flag.String("config", "", "configuration file path")
	fdir     = flag.String("dir", "", "output directory")
)

// TODO add database structure (optional)

var restapi genapi.RestAPI

func parseConfig() error {
	viper.SetConfigName(strings.TrimSuffix(filepath.Base(*fconfig), filepath.Ext(filepath.Base(*fconfig))))
	viper.AddConfigPath(filepath.Dir(*fconfig))
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("unable to read configuration file: %v", err)
	}
	// Configuration for the micro-service
	if err := viper.Unmarshal(&restapi); err != nil {
		return fmt.Errorf("unable to decode into config struct: %v", err)
	}
	return nil
}

func verifyConfig() error {
	return nil
}

func main() {
	flag.Parse()
	if *fversion {
		fmt.Println(version)
		return
	}
	// configuration
	if err := parseConfig(); err != nil {
		log.Fatal("unable to parse configuration file:", err)
	}
	if err := verifyConfig(); err != nil {
		log.Fatal("missing configuration parameters:", err)
	}
	// generation
	if err := restapi.Generate(*fdir); err != nil {
		panic(err)
	}
}
