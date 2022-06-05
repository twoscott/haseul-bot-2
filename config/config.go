package config

import (
	"os"
	"strings"
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"gopkg.in/yaml.v3"
)

type BotConfig struct {
	Discord struct {
		Token        string            `yaml:"token"`
		LogChannelID discord.ChannelID `yaml:"logChannelID"`
		RootGuildID  discord.GuildID   `yaml:"rootGuildID"`
	} `yaml:"discord"`
	PostgreSQL struct {
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"postgresql"`
	LastFm struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
	} `yaml:"lastFm"`
	Patreon struct {
		CampaignID   string `yaml:"campaignID"`
		AccessToken  string `yaml:"accessToken"`
		RefreshToken string `yaml:"refreshToken"`
		Secret       string `yaml:"secret"`
	} `yaml:"patreon"`
	Twitter struct {
		ConsumerKey    string `yaml:"consumerKey"`
		ConsumerSecret string `yaml:"consumerSecret"`
	} `yaml:"twitter"`
}

const fileName = "config.yml"

var (
	config *BotConfig
	once   sync.Once
)

// GetInstance returns the instance of Config.
func GetInstance() *BotConfig {
	once.Do(func() {
		cwd, _ := os.Getwd()
		sep := os.PathSeparator

		var path strings.Builder
		path.WriteString(cwd)
		path.WriteRune(sep)
		path.WriteString(fileName)

		file, err := os.Open(path.String())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)

		config = new(BotConfig)
		err = decoder.Decode(config)
		if err != nil {
			panic(err)
		}
	})

	return config
}
