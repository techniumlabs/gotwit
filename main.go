package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/techniumlabs/gotwit/cmd"
)

var (
	versionString = "undefined"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	cmd.Execute()
}
