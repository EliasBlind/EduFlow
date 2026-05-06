package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Translator struct {
	translations map[string]map[string]string
	defaultLang  string
}

func NewTranslator(fds fs.FS, dirPath string, defaultLang string) (*Translator, error) {
	t := &Translator{
		translations: make(map[string]map[string]string),
		defaultLang:  defaultLang,
	}

	files, err := fs.ReadDir(fds, dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locales dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		lang := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		fullPath := filepath.Join(dirPath, file.Name())
		data, err := fs.ReadFile(fds, fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", fullPath, err)
		}

		var m map[string]string
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("failed to unmarshal %s: %w", fullPath, err)
		}

		t.translations[lang] = m
	}

	if _, ok := t.translations[defaultLang]; !ok {
		return nil, fmt.Errorf("default language '%s' not found in loaded locales", defaultLang)
	}

	return t, nil
}

func (t *Translator) Translate(key, lang string) string {
	messages, ok := t.translations[lang]
	if !ok {
		messages = t.translations[t.defaultLang]
	}

	if val, ok := messages[key]; ok {
		return val
	}

	if val, ok := t.translations[t.defaultLang][key]; ok {
		return val
	}

	return key
}
