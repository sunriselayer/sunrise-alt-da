package config

import (
	toml "github.com/pelletier/go-toml"
)

type Config struct {
	Chain struct {
		AddrPrefix       string `toml:"addr_prefix"`
		PublisherAccount string `toml:"publisher_account"`
		HomePath         string `toml:"home_path"`
		KeyringBackend   string `toml:"keyring_backend"`
		Fees             string `toml:"fees"`
		CometbftRPC      string `toml:"cometbft_rpc"`
	}
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	configTree, err := toml.LoadFile("config.toml")
	if err != nil {
		return nil, err
	}
	err = configTree.Unmarshal(config)
	return config, err
}
