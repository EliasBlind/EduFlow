package main

import (
	"flag"
	"fmt"
	"strings"

	envutil "github.com/EliasBlind/EduFlow/pkg/env"
	"github.com/EliasBlind/EduFlow/pkg/key"
)

type EnvUpdater struct {
	envutil.EnvFile
}

func main() {
	envPath := flag.String("env", ".env", "Path to .env file")
	keyName := flag.String("key", "", "Variable name in .env (e.g., SECRET_KEY)")
	seedValue := flag.String("sed", "", "Manual value (if empty, uses generator)")
	genType := flag.String("gen", "random", "Generation type: 'random' or 'db'")

	flag.Parse()

	if *keyName == "" {
		fmt.Println("Error: --key is required")
		return
	}

	eu := NewEnvUpdater(*envPath)
	defer eu.Save()

	finalValue := *seedValue

	if finalValue == "" {
		switch *genType {
		case "db":
			finalValue = key.GenerateDBPassword()
		default:
			finalValue = key.GenerateRandomKey()
		}
	}

	eu.Set(*keyName, formatValue(finalValue))
	fmt.Printf("Updated %s in %s\n", *keyName, *envPath)
}

func NewEnvUpdater(path string) *EnvUpdater {
	envFile, err := envutil.LoadEnv(path)
	if err != nil {
		panic(fmt.Sprintf("There is no such file: %s", err))
	}
	return &EnvUpdater{*envFile}
}

func formatValue(val string) string {
	val = strings.TrimSpace(val)
	if strings.ContainsAny(val, " #\n\t") && !strings.HasPrefix(val, "\"") {
		return fmt.Sprintf("\"%s\"", val)
	}
	return val
}
