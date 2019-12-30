package sapphire

import (
	"fmt"
)

type Language struct {
	Name string
	Keys map[string]string
}

// NewLanguage creates a new language with the specified name.
func NewLanguage(name string) *Language {
	return &Language{Name: name, Keys: make(map[string]string)}
}

// Merge merges the keys from the other language
func (l *Language) Merge(other *Language) *Language {
	for k, v := range other.Keys {
		l.Keys[k] = v
	}
	return l
}

func (l *Language) Set(key string, value string) *Language {
	l.Keys[key] = value
	return l
}

func (l *Language) Get(key string, args ...interface{}) string {
	v, ok := l.Keys[key]
	if ok {
		return fmt.Sprintf(v, args...)
	}
	return ""
}

func (l *Language) GetDefault(key string, def string, args ...interface{}) string {
	v := l.Get(key, args...)
	if v == "" {
		return def
	}
	return v
}

var English = NewLanguage("en-US").
	Set("LOCALE_NO_KEY", "No localization found for the key \"%s\" Please report this to the developers.").
	Set("COMMAND_ERROR", "Something went wrong, please try again later.").
	Set("COMMAND_PING", "Pong!").
	Set("COMMAND_PING_PONG", "Pong! Latency: **%d**ms, API Latency: **%d**ms").
	Set("COMMAND_ENABLE_ALREADY", "That command is already enabled!").
	Set("COMMAND_DISABLE_ALREADY", "That command is already disabled!").
	Set("COMMAND_ENABLE_SUCCESS", "Successfully enabled the command **%s**").
	Set("COMMAND_DISABLE_SUCCESS", "Successfully disabled the command **%s**").
	Set("COMMAND_NOT_FOUND", "Command '%s' not found.").
	Set("COMMAND_INVITE", "To invite me to your server: <%s>").
	Set("COMMAND_OWNER_ONLY", "This command is for the bot owner only!").
	Set("COMMAND_GUILD_ONLY", "This command can only be used in a server!").
	Set("COMMAND_COOLDOWN", "You can use this command again in %d seconds.").
	Set("COMMAND_DISABLED", "This command has been disabled globally by the bot owner.")
