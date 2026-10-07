package envutil

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type line struct {
	key   string
	raw   string
	isVar bool
}

type EnvFile struct {
	path  string
	lines []line
}

func LoadEnv(pathToEnv string) (*EnvFile, error) {
	if err := validateEnv(pathToEnv); err != nil {
		return nil, err
	}

	env := &EnvFile{path: pathToEnv}

	file, err := os.Open(pathToEnv)
	if err != nil {
		if os.IsNotExist(err) {
			return env, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		trimmed := strings.TrimSpace(text)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			env.lines = append(env.lines, line{raw: text, isVar: false})
			continue
		}

		parts := strings.SplitN(text, "=", 2)
		if len(parts) == 2 {
			env.lines = append(env.lines, line{
				key:   strings.TrimSpace(parts[0]),
				raw:   text,
				isVar: true,
			})
		}
	}

	return env, scanner.Err()
}

func (e *EnvFile) Set(key, value string) {
	key = strings.TrimSpace(key)
	newValue := fmt.Sprintf("%s=%s", key, strings.TrimSpace(value))

	for i, l := range e.lines {
		if l.isVar && l.key == key {
			e.lines[i].raw = newValue
			return
		}
	}

	e.lines = append(e.lines, line{key: key, raw: newValue, isVar: true})
}

func (e *EnvFile) Save() error {
	var output strings.Builder
	for _, l := range e.lines {
		output.WriteString(l.raw + "\n")
	}
	return os.WriteFile(e.path, []byte(output.String()), 0644)
}

func validateEnv(env string) error {
	ext := filepath.Ext(env)
	if ext != ".env" {
		return fmt.Errorf("invalid extension: %s", env)
	}

	if _, err := os.Stat(env); os.IsNotExist(err) {
		return fmt.Errorf("config file is not exist: %s", env)
	}
	return nil
}

func (ef *EnvFile) Path() string {
	return ef.path
}
