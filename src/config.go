package src

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
)

type RunnerConfig struct {
	Name      string `toml:"name"`
	URL       string `toml:"url"`
	Token     string `toml:"token"`
	Limit     int    `toml:"limit"`
	Executor  string `toml:"executor"`
	BuildsDir string `toml:"builds_dir"`

	ShellScript string `toml:"shell_script"`
}

type BaseConfig struct {
	Concurrent int             `toml:"concurrent"`
	RootDir    string          `toml:"root_dir"`
	Runners    []*RunnerConfig `toml:"runners"`
}

type Config struct {
	BaseConfig
	ModTime time.Time
}

func (c *RunnerConfig) GetBuildsDir() string {
	if len(c.BuildsDir) == 0 {
		return "tmp/builds"
	} else {
		return c.BuildsDir
	}
}

func (c *RunnerConfig) ShortDescription() string {
	return c.Token[0:8]
}

func (config *Config) LoadConfig(config_file string) error {
	info, err := os.Stat(config_file)
	if err != nil {
		return err
	}

	if _, err = toml.DecodeFile(config_file, &config.BaseConfig); err != nil {
		return err
	}

	if config.Concurrent == 0 {
		config.Concurrent = 1
	}

	config.ModTime = info.ModTime()
	return nil
}

func (config *Config) SaveConfig(config_file string) error {
	var new_config bytes.Buffer
	new_buffer := bufio.NewWriter(&new_config)

	if err := toml.NewEncoder(new_buffer).Encode(&config.BaseConfig); err != nil {
		log.Fatalf("Error encoding TOML: %s", err)
		return err
	}

	if err := new_buffer.Flush(); err != nil {
		return err
	}

	if err := ioutil.WriteFile(config_file, new_config.Bytes(), 0600); err != nil {
		return err
	}

	return nil
}

func ReloadConfig(config_file string, config_time time.Time, reload_config chan Config) {
	for {
		time.Sleep(RELOAD_CONFIG_INTERVAL * time.Second)

		info, err := os.Stat(config_file)
		if err != nil {
			log.Errorln("Failed to stat config", err)
			continue
		}

		if config_time.Before(info.ModTime()) {
			config_time = info.ModTime()

			new_config := Config{}
			err := new_config.LoadConfig(config_file)
			if err != nil {
				log.Errorln("Failed to load config", err)
				continue
			}

			reload_config <- new_config
		}
	}
}

func (c *Config) SetChdir() {
	if len(c.RootDir) > 0 {
		err := os.Chdir(c.RootDir)
		if err != nil {
			panic(err)
		}
	}
}