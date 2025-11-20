package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ydcloud-dy/leaf-api/cmd"
)

var (
	configPath string
	version    = "1.0.0"
	showVer    bool
)

func init() {
	flag.StringVar(&configPath, "config", "", "config file path (default: ./config.yaml)")
	flag.StringVar(&configPath, "c", "", "config file path (shorthand)")
	flag.BoolVar(&showVer, "version", false, "show version")
	flag.BoolVar(&showVer, "v", false, "show version (shorthand)")
}

func main() {
	flag.Parse()

	if showVer {
		fmt.Printf("Blog Admin API v%s\n", version)
		os.Exit(0)
	}

	// 运行应用
	if err := cmd.Run(configPath); err != nil {
		fmt.Printf("Failed to run application: %v\n", err)
		os.Exit(1)
	}
}
