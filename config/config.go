package config

import (
	"context"
	"log"
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Discord struct {
		Token string `env:"TOKEN,required"`
	} `env:",prefix=DISCORD_"`
	Bot struct {
		LogChannelID discord.ChannelID `env:"LOG_CHANNEL_ID"`
		HomeGuildID  discord.GuildID   `env:"HOME_GUILD_ID,required"`
		AdminUserID  discord.UserID    `env:"ADMIN_USER_ID,required"`
	} `env:",prefix=BOT_"`
	PostgreSQL struct {
		Host     string `env:"HOST"`
		Port     string `env:"PORT"`
		Database string `env:"DB"`
		Username string `env:"USER"`
		Password string `env:"PASSWORD,required"`
	} `env:",prefix=POSTGRES_"`
	LastFm struct {
		Key    string `env:"KEY"`
		Secret string `env:"SECRET"`
	} `env:",prefix=LASTFM_"`
	Patreon struct {
		CampaignID  string `env:"CAMPAIGN_ID"`
		AccessToken string `env:"ACCESS_TOKEN"`
	} `env:",prefix=PATREON_"`
	SushiiImageServer struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT"`
	} `env:",prefix=SUSHII_IMAGE_SERVER_"`
}

var (
	config *Config
	once   sync.Once
)

// GetInstance returns the instance of Config.
func GetInstance() *Config {
	once.Do(func() {
		// ignore any errors since the .env file won't be passed to the docker
		// context and instead the environment variables will be passed from
		// the file.
		godotenv.Load("local.env")

		config = new(Config)
		err := envconfig.Process(context.Background(), config)
		if err != nil {
			log.Fatalf("Failed to process environment variables: %s\n", err)
		}
	})

	return config
}
