package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"time"
)

const UserConfigFileName = "cli.json"
const ConfigEnv = "GEENY_CLI_CONFIG"

type ConfigExt struct {
	Log                bool          `json:"log.enable"`
	LogTrace           bool          `json:"log.trace"`
	LogInfo            bool          `json:"log.info"`
	LogWarn            bool          `json:"log.warn"`
	LogError           bool          `json:"log.error"`
	Spinner            bool          `json:"output.spinner"`
	ConfigEnv          bool          `json:"config.env"`
	AutoUpdateCheck    bool          `json:"update.autocheck"`
	SwaggerValidate    bool          `json:"swagger.validate"`
	CaPemFilename      string        `json:"net.ca"`
	Proxy              bool          `json:"net.proxy"`
	Timeout            time.Duration `json:"net.timeout"`
	DisplayJson        bool          `json:"output.displayJson"`
	RawJson            bool          `json:"output.rawJson"`
	ConnectUrl         string        `json:"connect.url"`
	ConnectSwaggerUrl  string        `json:"connectSwagger.url"`
	ApiUrl             string        `json:"api.url"`
	ApiSwaggerUrl      string        `json:"apiSwagger.url"`
	DefaultContentType string        `json:"swagger.defaultContentType"`
}

type ConfigInt struct {
	IsDebug       bool
	SpinnerOutput bool
}

var Default = ConfigExt{
	Log:                false,
	LogTrace:           true,
	LogInfo:            true,
	LogWarn:            true,
	LogError:           true,
	Spinner:            true,
	ConfigEnv:          true,
	AutoUpdateCheck:    true,
	SwaggerValidate:    false,
	CaPemFilename:      "",
	Proxy:              false,
	Timeout:            30 * time.Second,
	DisplayJson:        false,
	RawJson:            true,
	ConnectUrl:         "https://connect.geeny.io/",
	ConnectSwaggerUrl:  "",
	ApiUrl:             "https://api.geeny.io/",
	ApiSwaggerUrl:      "",
	DefaultContentType: "application/json",
}

var CurrentExt = getCurrent()

// - private

func getCurrent() *ConfigExt {
	user, err := getUserConfig(&Default)
	if err != nil {
		user = &Default
	}
	if user.ConfigEnv {
		cwd, _ := getEnvConfig(user)
		return cwd
	}
	return user
}

func getEnvConfig(parent *ConfigExt) (*ConfigExt, error) {
	envConfig := os.Getenv(ConfigEnv)
	if len(envConfig) == 0 {
		return parent, nil
	}
	return readConfig(parent, envConfig)
}

func getUserConfig(parent *ConfigExt) (*ConfigExt, error) {
	usr, err := user.Current()
	if err != nil {
		return parent, err
	}
	configPath := usr.HomeDir + string(os.PathSeparator) + ".geeny" + string(os.PathSeparator) + UserConfigFileName
	return readConfig(parent, configPath)
}

func readConfig(parent *ConfigExt, configPath string) (*ConfigExt, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return parent, err
	}
	config := parent
	err = json.Unmarshal(file, &config)
	if err != nil {
		return parent, err
	}
	return config, err
}
