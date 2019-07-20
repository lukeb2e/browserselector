package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"

	"github.com/spf13/viper"
)

type configuration struct {
	Domain []browser
	Debug  bool
}

type browser struct {
	Browser  string
	Regex    string
	Priority uint
}

func debug(debug bool, a ...interface{}) (n int, err error) {
	if !debug {
		return
	}
	return fmt.Fprintln(os.Stdout, a...)
	//return fmt.Fprintln(os.Stdout, append([]interface{}{"[Debug]"}, a...))
}

func main() {
	// Check if config file exists
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		fmt.Println("Configuration file does not exist.")
		os.Exit(1)
	}

	// Load config
	var config configuration
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Sort browser by priority
	sort.Slice(config.Domain[:], func(i, j int) bool {
		return config.Domain[i].Priority < config.Domain[j].Priority
	})

	// Get OS Arguments
	args := os.Args[1]
	debug(config.Debug, "Arguments:", args)

	// Regex match FQDN
	// https?:\/\/([^\/]*)\/?.*
	r, err := regexp.Compile("https?://([^/]*)/?.*")
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}
	fqdn := r.FindStringSubmatch(args)[1]
	debug(config.Debug, "FQDN: ", fqdn)

	// iterate over config
	selector := len(config.Domain) - 1
	for index, element := range config.Domain {
		match, _ := regexp.MatchString(element.Regex, fqdn)
		if match {
			selector = index
			debug(config.Debug, "Match found")
			break
		}
	}

	// Select correct browser
	// Start browser
	cmd := exec.Command(config.Domain[selector].Browser, args)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}
}
