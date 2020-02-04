//go:generate goversioninfo -icon=./resource/icon.ico ./resource/versioninfo.json

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/spf13/viper"
)

type configuration struct {
	Domain  []domain
	Browser map[string]*browser
	Debug   bool
}

type domain struct {
	Browser  string
	Regex    string
	Priority uint
}

type browser struct {
	Exec   string
	Script string
}

func debug(debug bool, a ...interface{}) (n int, err error) {
	if !debug {
		return
	}
	return fmt.Fprintln(os.Stdout, a...)
	//return fmt.Fprintln(os.Stdout, append([]interface{}{"[Debug]"}, a...))
}

func main() {
	// Get location of running binary
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Load config
	var config configuration
	viper.AddConfigPath(".")
	viper.AddConfigPath(dir)
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
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
	args := os.Args
	debug(config.Debug, "Arguments:", args)

	var uri string
	if len(args) > 1 {
		for index, element := range args {
			debug(config.Debug, "Element:", index, element)
			start, err := regexp.Compile("(http|https|ftp|ftps|ftpes).*")
			if err != nil {
				fmt.Println(err)
				fmt.Scanln()
				os.Exit(1)
			}
			if start.MatchString(element) {
				uri = args[index]
				break
			}
		}
	}

	if uri == "" {
		fmt.Println("Missing URL")
		os.Exit(1)
	}

	debug(config.Debug, "URI: ", uri)

	// Regex match FQDN
	// https?:\/\/([^\/]*)\/?.*
	r, err := regexp.Compile("(http|https|ftp|ftps|ftpes)://([^/]*)/?.*")
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}
	fqdn := r.FindStringSubmatch(uri)[2]
	debug(config.Debug, "FQDN: ", fqdn)

	// iterate over config
	selector := len(config.Domain) - 1
	for index, element := range config.Domain {
		match, _ := regexp.MatchString(element.Regex, fqdn)
		if match {
			selector = index
			debug(config.Debug, "Match found", element.Browser, "Priority:", element.Priority, "Regex:", element.Regex)
			break
		}
	}

	// Check if browser exists
	if _, ok := config.Browser[config.Domain[selector].Browser]; !ok {
		fmt.Println("Browser not found in configuration")
		fmt.Scanln()
		os.Exit(1)
	}

	// Start browser
	var command = config.Browser[config.Domain[selector].Browser].Exec
	var cmdArgs []string
	if config.Browser[config.Domain[selector].Browser].Script == "" {
		// Exe + "FQDN"
		//cmdArgs = append(cmdArgs, "\""+uri+"\"")
		cmdArgs = append(cmdArgs, uri)
	} else {
		// Exe + Script + "FQDN"
		cmdArgs = append(cmdArgs, config.Browser[config.Domain[selector].Browser].Script, "\""+uri+"\"")
	}

	fmt.Println(command, cmdArgs)
	cmd := exec.Command(command, cmdArgs...)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}

	if config.Debug {
		fmt.Println()
		fmt.Scanln()
	}
}
