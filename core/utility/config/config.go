package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const dirName = ".inplainsight"
const fileName = "config.json"

type Config struct {
	PoolPath           string `json:"pool_path"`
	AFKTimeout         int    `json:"afk_timeout"`
	LogoutOnScreenLock bool
}

func (c *Config) Save() error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s/", strings.TrimRight(userConfigDir, "/"), strings.Trim(dirName, "/"))

	if err := os.MkdirAll(path, 0700); err == nil {
		path += fileName

		handle, err := os.Create(path)
		if err != nil {
			return err
		}
		defer handle.Close()

		encoder := json.NewEncoder(handle)
		err = encoder.Encode(c)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func Load() (c *Config, err error) {
	c = &Config{}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	path := fmt.Sprintf(
		"%s/%s/%s",
		strings.TrimRight(userConfigDir, "/"),
		strings.Trim(dirName, "/"),
		fileName,
	)

	handle, err := os.Open(path)
	if err != nil {
		return
	}
	defer handle.Close()

	decoder := json.NewDecoder(handle)
	err = decoder.Decode(&c)

	return
}
