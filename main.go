//go:generate goversioninfo -icon=./resource/icon.ico ./resource/versioninfo.json

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/spf13/viper"
)

func debugOutput(debug bool, a ...interface{}) (n int, err error) {
	if !debug {
		return
	}
	return fmt.Fprintln(os.Stdout, a...)
	//return fmt.Fprintln(os.Stdout, append([]interface{}{"[Debug]"}, a...))
}

func getBinaryDirectory(args []string) (dir string, err error) {
	if len(args) < 1 {
		err = errors.New("too few arguments")
		return
	}

	dir, err = filepath.Abs(filepath.Dir(args[0]))
	return dir, err
}

func sortConfigBrowserPriority(input []domain) (output []domain, err error) {
	sort.Slice(input[:], func(i, j int) bool {
		return input[i].Priority < input[j].Priority
	})
	output = input
	return
}

func getUrl(args []string, config configuration) (url string, err error) {
	debugOutput(config.Debug, "Arguments:", args)

	if len(args) < 1 {
		err = errors.New("missing parameters")
		return
	}

	for index, element := range args {
		debugOutput(config.Debug, "Element:", index, element)
		start, err := regexp.Compile("(http|https|ftp|ftps|ftpes|file).*")
		if err != nil {
			fmt.Println(err)
			fmt.Scanln()
			os.Exit(1)
		}
		if start.MatchString(element) {
			url = args[index]
			break
		}
	}

	if url == "" {
		err = errors.New("no url found")
		return
	}

	debugOutput(config.Debug, "URL: ", url)
	return
}

func getFqdnFromUrl(url string, config configuration) (protocol string, fqdn string, err error) {
	// Regex match FQDN
	// https?:\/\/([^\/]*)\/?.*
	r, err := regexp.Compile("(http|https|ftp|ftps|ftpes)://([^/]*)/?.*")
	if err != nil {
		return
	}
	matches := r.FindStringSubmatch(url)

	debugOutput(config.Debug, "Matches: ", matches)

	if len(matches) < 3 {
		err = errors.New("invalid url: " + url)
		return
	}

	protocol = matches[1]
	fqdn = matches[2]

	debugOutput(config.Debug, "Protocol: ", protocol, " | FQDN: ", fqdn)
	return
}

func main() {
	dir, err := getBinaryDirectory(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
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
		fmt.Scanln()
		return
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}

	config.Domain, err = sortConfigBrowserPriority(config.Domain)
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}

	// Get url from arguments
	url, err := getUrl(os.Args, config)
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}

	_, fqdn, err := getFqdnFromUrl(url, config)
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		return
	}

	// Check rules to select browser
	selector := len(config.Domain) - 1
	for index, element := range config.Domain {
		match, _ := regexp.MatchString(element.Regex, fqdn)
		if match {
			selector = index
			debugOutput(config.Debug, "Match found", element.Browser, "Priority:", element.Priority, "Regex:", element.Regex)
			break
		}
	}

	// Check if browser exists
	if _, ok := config.Browser[config.Domain[selector].Browser]; !ok {
		if config.Debug {
			fmt.Println("Browser not found in configuration:", config.Domain[selector].Browser)
			fmt.Scanln()
		}
		os.Exit(1)
	}

	// Start browser
	var command = config.Browser[config.Domain[selector].Browser].Exec
	var cmdArgs []string
	if config.Browser[config.Domain[selector].Browser].Script == "" {
		// Exe + "FQDN"
		//cmdArgs = append(cmdArgs, "\""+url+"\"")
		cmdArgs = append(cmdArgs, url)
	} else {
		// Exe + Script + "FQDN"
		cmdArgs = append(cmdArgs, config.Browser[config.Domain[selector].Browser].Script, "\""+url+"\"")
	}

	debugOutput(config.Debug, command, cmdArgs)
	cmd := exec.Command(command, cmdArgs...)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		fmt.Scanln()
		os.Exit(1)
	}

	// Stop execution to show debug output
	if config.Debug {
		fmt.Println()
		fmt.Scanln()
	}
}
