package version

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	log "geeny/log"
)

const Application = "Geeny CLI"

var Version = buildVersion()
var UpdateChecked bool
var UserAgent = "Geeny-CLI/" + version

var version string
var major string
var minor string
var patch string
var timestamp string
var updateCheckResult bool
var updateCheckMessage string

const versionRegex = "^v?([0-9]+)\\.([0-9]+)\\.?([0-9a-zA-Z_-]*)"
const versionUrl = "https://developers.geeny.io/downloads/cli/version.txt"
const messageTemplate = "New %s version available: %s\n%s\n%s\n"
const helpMsg = "More details and options available at https://developers.geeny.io/documentation/cli/"
const downloadMsg = "curl https://developers.geeny.io/downloads/cli/osx/amd64/geeny -o geeny && chmod +x geeny"
const winDownloadMsg = "Download from https://developers.geeny.io/downloads/cli/windows/amd64/geeny.exe"

func CheckUpdate(timeout bool) (bool, string) {
	if !UpdateChecked {
		log.Tracef("Checking for new CLI version")
		req, err := http.NewRequest("GET", versionUrl, nil)
		if err != nil {
			log.Warn("Could not create version check request", err)
			return false, ""
		}
		req.Header.Set("User-Agent", UserAgent)
		timeoutDuration := 1000 * time.Millisecond
		if !timeout {
			timeoutDuration = 0
		}
		client := &http.Client{Timeout: timeoutDuration}
		resp, err := client.Do(req)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Warn("Could not perform version check", err)
			return false, ""
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn("Could not read version check response", err)
			return false, ""
		}
		updateCheckResult, updateCheckMessage = checkVersion(string(body))
		UpdateChecked = true
	}
	return updateCheckResult, updateCheckMessage
}

func buildVersion() string {
	if major != "" && minor != "" {
		return versionDescriptor()
	} else if version != "" {
		regex := regexp.MustCompile(versionRegex)
		match := regex.FindStringSubmatch(version)
		if match != nil && len(match) >= 3 {
			major = string(match[1])
			minor = string(match[2])
			if len(match) >= 4 {
				patch = string(match[3])
			}
			return versionDescriptor()
		}
	}
	return fmt.Sprintf("%s %s/%s", Application, runtime.GOOS, runtime.GOARCH)
}

func versionDescriptor() string {
	if timestamp != "" {
		return fmt.Sprintf("%s v%s (%s) %s/%s", Application, versionSemantic(), timestamp, runtime.GOOS, runtime.GOARCH)
	}
	return fmt.Sprintf("%s v%s %s/%s", Application, versionSemantic(), runtime.GOOS, runtime.GOARCH)
}

func versionSemantic() string {
	if len(patch) == 0 {
		return fmt.Sprintf("%s.%s", major, minor)
	}
	return fmt.Sprintf("%s.%s.%s", major, minor, patch)
}

func checkVersion(newVersionRaw string) (bool, string) {
	newVersion := strings.TrimSpace(newVersionRaw)
	if len(newVersion) < 1 {
		return false, ""
	}
	regex := regexp.MustCompile(versionRegex)
	match := regex.FindStringSubmatch(newVersion)
	if match != nil && len(match) > 2 {
		var dMsg string
		if runtime.GOOS == "windows" {
			dMsg = winDownloadMsg
		} else {
			dMsg = downloadMsg
		}

		log.Tracef("Comparing current version %v to my version %s", match, version)
		if string(match[1]) > major {
			log.Tracef("New major version available: %s", newVersion)
			return true, fmt.Sprintf(messageTemplate, "major", newVersion, dMsg, helpMsg)
		} else if string(match[1]) == major && string(match[2]) > minor {
			log.Tracef("New minor version available: %s", newVersion)
			return true, fmt.Sprintf(messageTemplate, "minor", newVersion, dMsg, helpMsg)
		} else if string(match[2]) == minor && len(match) > 3 && string(match[3]) > patch {
			log.Tracef("New patch version available: %s", newVersion)
			return true, fmt.Sprintf(messageTemplate, "patch", newVersion, dMsg, helpMsg)
		}
	}
	log.Tracef("Already at newest version: %s", version)
	return false, ""
}
