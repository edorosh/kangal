package main

import (
	"log"

	"github.com/hellofresh/kangal/cmd"
)

var version = "0.0.0-dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)

	if err := rootCmd.Execute(); err != nil {
		// There is no way to output logs in JSON format for parsing by logstash
		log.Fatal(err)
	}
}
