package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	configFile = os.ExpandEnv("$HOME/.ssh/accountswitcherconfig.json")
	tempConfig = os.ExpandEnv("$HOME/.ssh/accountswitcherconfig.txt")
	conf       map[string]AccountConfig
)

// AccountConfig holds the configuration for each GitHub account.
type AccountConfig struct {
	Prefix string `json:"prefix"`
	Email  string `json:"email"`
}

// readConfigFile reads the JSON configuration file into a map.
func readConfigFile() (map[string]AccountConfig, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config map[string]AccountConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// write writes the current account to the temporary config file.
func write(account string) error {
	return ioutil.WriteFile(tempConfig, []byte(strings.TrimSpace(account)), 0644)
}

// fixURL constructs a GitHub URL based on the account prefix and repository name.
func fixURL(account, repo string) (string, error) {
	data, exists := conf[account]
	if !exists {
		return "", fmt.Errorf("account %s not found", account)
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid repo format: %s", repo)
	}

	url := fmt.Sprintf("%s:%s/%s.git", data.Prefix, parts[0], parts[1])
	fmt.Println(url)

	return url, nil
}

// changeAccount changes the Git configuration to match the email associated with the specified account.
func changeAccount(account string) error {
	if err := write(account); err != nil {
		return err
	}

	data, exists := conf[account]
	if !exists {
		return fmt.Errorf("account %s not found", account)
	}

	emailCommand := fmt.Sprintf("git config --global user.email %s", data.Email)
	fmt.Println(emailCommand)
	return runCommand(emailCommand)
}

// cloneRepo clones the repository using the constructed URL.
func cloneRepo(url string) error {
	command := fmt.Sprintf("git clone %s", url)
	return runCommand(command)
}

// runCommand executes the specified shell command.
func runCommand(command string) error {
	fmt.Println("running command =>", command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// info prints the content of the temporary config file.
func info() error {
	data, err := ioutil.ReadFile(tempConfig)
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: accountswitcher <account> [repo]")
		os.Exit(1)
	}

	var err error
	conf, err = readConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		os.Exit(1)
	}

	account := os.Args[1]
	if len(os.Args) > 2 {
		repo := os.Args[2]
		if err := changeAccount(account); err != nil {
			fmt.Fprintf(os.Stderr, "Error changing account: %v\n", err)
			os.Exit(1)
		}

		url, err := fixURL(account, repo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fixing URL: %v\n", err)
			os.Exit(1)
		}

		if err := cloneRepo(url); err != nil {
			fmt.Fprintf(os.Stderr, "Error cloning repo: %v\n", err)
			os.Exit(1)
		}
	} else {
		if account == "personal" || account == "work" {
			if err := changeAccount(account); err != nil {
				fmt.Fprintf(os.Stderr, "Error changing account: %v\n", err)
				os.Exit(1)
			}
		} else if account == "info" {
			if err := info(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading info: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Unknown command")
			os.Exit(1)
		}
	}
}
