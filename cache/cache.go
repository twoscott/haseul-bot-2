package cache

import (
	"log"
	"sync"

	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/twoscott/haseul-bot-2/database"
	"github.com/twoscott/haseul-bot-2/database/guilddb"
	"github.com/twoscott/haseul-bot-2/router"
)

type (
	configMap map[discord.GuildID]*guilddb.GuildConfig

	// Cache wraps the database and stores data in memory for quicker access.
	Cache struct {
		*database.DB
		configs configMap
	}
)

var (
	cache *Cache
	once  sync.Once
)

// GetInstance returns the instance of Cache.
func GetInstance() *Cache {
	once.Do(func() {
		cache = &Cache{
			DB:      database.GetInstance(),
			configs: make(configMap),
		}
	})

	return cache
}

func (c *Cache) onStartup(rt *router.Router, ev *gateway.ReadyEvent) {
	// TEMPORARY
	guilds, _ := rt.State.Guilds()
	for _, guild := range guilds {
		c.Guilds.Add(guild.ID)
	}

	configs, err := c.Guilds.GetConfigs()
	if err != nil {
		panic(err)
	}

	for _, config := range configs {
		c.configs[config.GuildID] = &config
	}
}

func (c *Cache) onGuildJoin(rt *router.Router, guild *state.GuildJoinEvent) {
	log.Println(guild.Name)
	c.AddGuild(guild.ID)
}

// AddGuild adds a guild config entry for the given guild ID to the cache.
func (c *Cache) AddGuild(guildID discord.GuildID) error {
	added, err := c.Guilds.Add(guildID)
	if err != nil {
		return err
	}
	if !added {
		return nil
	}

	config, err := c.Guilds.GetConfig(guildID)
	if err != nil {
		return err
	}

	c.configs[guildID] = config
	return nil
}

// GetGuildConfig returns a guild config for a given guild ID.
func (c Cache) GetGuildConfig(guildID discord.GuildID) *guilddb.GuildConfig {
	return c.configs[guildID]
}

// GetPrefix returns the bot's prefix for a guild ID.
func (c Cache) GetPrefix(guildID discord.GuildID) string {
	config := c.configs[guildID]
	if config == nil {
		return ""
	}

	return config.Prefix
}
