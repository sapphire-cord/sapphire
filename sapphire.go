// Sapphire is a complete discord bot framework for Golang with discordgo
package sapphire

import (
  "github.com/bwmarrin/discordgo"
  "github.com/dustin/go-humanize"
  "runtime"
  "strings"
  "fmt"
  "time"
  "os"
  "syscall"
  "os/signal"
)

// VERSION is a constant representing the current version of the framework.
const VERSION = "1.0.0"
// COLOR is the color for sapphire's embed colors.
const COLOR = 0x7F139E

type PrefixHandler func(b *Bot, m *discordgo.Message, dm bool) string
type LocaleHandler func(b *Bot, m *discordgo.Message, dm bool) string
type ErrorHandler func(b *Bot, err interface{})

// Bot represents a bot with sapphire framework features.
type Bot struct {
  Session *discordgo.Session // The discordgo session.
  Prefix PrefixHandler // The handler called to get the prefix. (default: !)
  Language LocaleHandler // The handler called to get the language (default: en-US)
  Commands map[string]*Command // Map of commands.
  CommandsRan int // Commands ran.
  Monitors map[string]*Monitor // Map of monitors.
  aliases map[string]string
  CommandCooldowns map[string]map[string]time.Time
  CommandEdits map[string]string
  OwnerID string // Bot owner's ID (default: fetched from application info)
  InvitePerms int // Permissions bits to use for the invite link. (default: 3072)
  Languages map[string]*Language // Map of languages.
  DefaultLocale *Language // Default locale to fallback. (default: en-US)
  CommandTyping bool // Wether to start typing when a command is being ran. (default: true)
  ErrorHandler ErrorHandler // The handler to catch panics in monitors (which includes commands).
  MentionPrefix bool // Wether to allow @mention of the bot to be used as a prefix too. (default: true)
  sweepTicker *time.Ticker
  Application *discordgo.Application // The bot's application.
  Uptime time.Time // The time the bot hit ready event.
}

// New creates a new sapphire bot, pass in a discordgo instance configured with your token.
func New(s *discordgo.Session) *Bot {
  bot := &Bot{
    Session: s,
    Prefix: func(_ *Bot, _ *discordgo.Message, _ bool) string {
      return "!" // A very common prefix, sigh, so we will make it the default.
    },
    Language: func(_ *Bot, _ *discordgo.Message, _ bool) string {
      return "en-US"
    },
    ErrorHandler: func(_ *Bot, err interface{}) {
      fmt.Printf("Panic recovered: %v\n", err)
    },
    Commands: make(map[string]*Command),
    aliases: make(map[string]string),
    Languages: make(map[string]*Language),
    CommandsRan: 0,
    InvitePerms: 3072,
    CommandCooldowns: make(map[string]map[string]time.Time),
    CommandEdits: make(map[string]string),
    Monitors: make(map[string]*Monitor),
    CommandTyping: true,
    sweepTicker: time.NewTicker(1 * time.Hour),
    Application: nil,
    MentionPrefix: true,
  }
  bot.AddLanguage(English)
  bot.SetDefaultLocale("en-US")
  bot.AddMonitor(NewMonitor("commandHandler", CommandHandlerMonitor).AllowEdits())
  s.AddHandler(monitorListener(bot))
  s.AddHandler(monitorEditListener(bot))
  s.AddHandlerOnce(func(s *discordgo.Session, ready *discordgo.Ready) {
    bot.Uptime = time.Now()

    // Sweeps all cooldowns/edits every hour to prevent infinite memory usage
    // While even active cooldowns gets reset it is fine though, as its only hourly
    // and is not too common for users to even notice it, same for edits.
    go func() {
      <-bot.sweepTicker.C
      bot.CommandCooldowns = make(map[string]map[string]time.Time)
      bot.CommandEdits = make(map[string]string)
    }()

    // TODO: for some reason it says bots cannot use this endpoint, i've seen a similar usecase before
    // try to figure out a way.
    /*app, err := s.Application(ready.User.ID)
    if err != nil {
      bot.ErrorHandler(bot, err)
      return
    p}
    bot.Application = app
    if bot.OwnerID == "" { bot.OwnerID = app.Owner.ID }*/
  })
  return bot
}

// SetMentionPrefix toggles the usage of the bot's @mention as a prefix.
func (bot *Bot) SetMentionPrefix(toggle bool) *Bot {
  bot.MentionPrefix = toggle
  return bot
}

// SetInvitePerms sets the permissions to request for in the bot invite link.
// The default is 3072 which is [VIEW_CHANNEL, SEND_MESSAGES]
func (bot *Bot) SetInvitePerms(bits int) *Bot {
  bot.InvitePerms = bits
  return bot
}

// SetErrorHandler sets the function to handle panics that happens in monitors (which includes commands)
func (bot *Bot) SetErrorHandler(fn ErrorHandler) *Bot {
  bot.ErrorHandler = fn
  return bot
}

// Sets the default locale to fallback when the bot can't find a key in the current locale.
// Panics if locale isn't registered.
func (bot *Bot) SetDefaultLocale(locale string) *Bot {
  if lang, ok := bot.Languages[locale]; !ok {
    panic(fmt.Sprintf("The language '%s' cannot be found.", locale))
  } else {
    bot.DefaultLocale = lang
  }
  return bot
}

func (bot *Bot) SetLocaleHandler(handler LocaleHandler) *Bot {
  bot.Language = handler
  return bot
}

// SetPrefixHandler sets the prefix handler, the function is responsible to return the right prefix for the command call.
// Use this for dynamic prefixes, e.g fetch prefix from database.
func (bot *Bot) SetPrefixHandler(prefix PrefixHandler) *Bot {
  bot.Prefix = prefix
  return bot
}

// SetPrefix sets a constant string as the prefix, use SetPrefixHandler if you need dynamic per-guild prefixes.
func (bot *Bot) SetPrefix(prefix string) *Bot {
  bot.Prefix = func(_ *Bot, _ *discordgo.Message, _ bool) string {
    return prefix
  }
  return bot
}

// Wait makes the bot wait until CTRL + C is pressed, this is used to keep the process alive.
// It closes the session when CTRL + C is pressed and you are free to do any extra cleanup after the call returns.
func (bot *Bot) Wait() {
  // Wait for an interrupt signal, e.g CTRL + C
  sc := make(chan os.Signal, 1)
  signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
  <-sc
  // Cleanly close down the Discord session.
  bot.Session.Close()
  bot.sweepTicker.Stop()
}

func (bot *Bot) AddCommand(cmd *Command) *Bot {
  c, ok := bot.Commands[cmd.Name]
  // If we are overriding an existing command ensure we unload any state it loaded in the bot, mainly the aliases.
  if ok {
    for _, a := range c.Aliases {
      delete(bot.aliases, a)
    }
  }
  bot.Commands[cmd.Name] = cmd
  for _, alias := range cmd.Aliases {
    bot.aliases[alias] = cmd.Name
  }
  return bot
}

// GetCommand returns a command by name, it also searches by aliases, returns nil if not found.
func (bot *Bot) GetCommand(name string) *Command {
  cmd, ok := bot.Commands[name]
  if ok { return cmd }
  alias, ok := bot.aliases[name]
  if ok { return bot.Commands[alias] }
  return nil
}

// Connect is an alias to discordgo's Session.Open
func (bot *Bot) Connect() error {
  return bot.Session.Open()
}

// MustConnect is like Connect but panics if there is an error.
func (bot *Bot) MustConnect() {
  if err := bot.Connect(); err != nil {
    panic(err)
  }
}

// AddLanguage adds the specified language.
func (bot *Bot) AddLanguage(lang *Language) *Bot {
  bot.Languages[lang.Name] = lang
  return bot
}

func (bot *Bot) AddMonitor(m *Monitor) *Bot {
  bot.Monitors[m.Name] = m
  return bot
}

// CheckCooldown checks the cooldown for userID for a command
// the first return is a bool indicating if the user can run the command.
// The second value is if user can't run then it will be the amount of seconds
// to wait before being able to.
// Note this function assumes the user will run the command and will place the user on cooldown if it isn't already.
func (bot *Bot) CheckCooldown(userID, command string, cooldownSec int) (bool, int) {
  if cooldownSec == 0 {
    return true, 0
  }

  cooldown := time.Duration(cooldownSec) * time.Second
  user, ok := bot.CommandCooldowns[userID]

  if !ok {
    bot.CommandCooldowns[userID] = make(map[string]time.Time)
    user = bot.CommandCooldowns[userID]
  }

  last, ok := user[command]

  if !ok {
    user[command] = time.Now()
    return true, 0
  }

  if !time.Now().After(last.Add(cooldown)) {
    return false, int(time.Until(last.Add(cooldown)).Seconds())
  }

  user[command] = time.Now()
  return true, 0
}

// LoadBuiltins loads the default set of builtin command, they are:
// ping, help, stats, invite, enable, disable, gc
// Some of the must have commands. (or rather commands that i feel good to have.)
func (bot *Bot) LoadBuiltins() *Bot {
  // To keep things simple all commands are declared here, we shouldn't need that much of builtins anyway.
  // And to keep the code easier to jump around this function is always the last.
  bot.AddCommand(NewCommand("ping", "General", func(ctx *CommandContext) {
    msg, err := ctx.ReplyLocale("COMMAND_PING")
    // Should never happen but if it did, avoid panics.
    if err != nil { return }
    usertime, err := ctx.Message.Timestamp.Parse()
    if err != nil { return }
    bottime, err := msg.Timestamp.Parse()
    if err != nil { return }
    ctx.EditLocale(msg, "COMMAND_PING_PONG", bottime.Sub(usertime).Milliseconds(), ctx.Session.HeartbeatLatency().Milliseconds())
  }).SetDescription("Pong! Responds with Bot latency."))

  bot.AddCommand(NewCommand("help", "General", func(ctx *CommandContext) {
    if ctx.HasArgs() { // User passed an argument, give help information on that command only.
      cmd := bot.GetCommand(ctx.Args[0].AsString())
      if cmd == nil {
        ctx.Reply("Unknown Command.")
        return
      }
      var aliases string = "None"

      if len(cmd.Aliases) > 0 {
        aliases = strings.Join(cmd.Aliases, ", ")
      }

      ctx.BuildEmbed(NewEmbed().
        SetDescription(fmt.Sprintf("**Name:** %s\n**Description:** %s\n**Category:** %s\n**Aliases:** %s\n**Usage:** %s",
          cmd.Name,
          cmd.Description,
          cmd.Category,
          aliases,
          fmt.Sprintf("%s%s %s", ctx.Prefix, cmd.Name, HumanizeUsage(cmd.UsageString)),
        )).SetColor(COLOR).SetTitle("Command Help"))
        return
      }
      // Send all commands.

      categories := make(map[string][]string)
      for _, v := range bot.Commands {
        _, ok := categories[v.Category]
        if !ok {
          categories[v.Category] = []string{}
        }
        if !v.OwnerOnly || ctx.Author.ID == ctx.Bot.OwnerID {
          categories[v.Category] = append(categories[v.Category], v.Name)
        }
      }

      // Filter out empty categories, e.g a whole category is full of owner commands and the commands are filtered
      // because the user isn't the owner, so avoid even mentioning that empty category.
      for k, v := range categories {
        if len(v) == 0 { delete(categories, k) }
      }

      var embed = &discordgo.MessageEmbed{
        Title: "Commands",
        Color: COLOR,
        Footer: &discordgo.MessageEmbedFooter{Text: "For more info on a command use: " + ctx.Prefix + "help <command>"},
        Author: &discordgo.MessageEmbedAuthor{IconURL: ctx.Author.AvatarURL("256"), Name: ctx.Author.Username},
      }

      for cat, cmds := range categories {
        var field = &discordgo.MessageEmbedField{Name: cat, Value: ""}
        field.Value = strings.Join(cmds, ", ")
        field.Inline = true
        embed.Fields = append(embed.Fields, field)
      }
      ctx.ReplyEmbed(embed)
    }).SetDescription("Shows a list of all commands.").SetUsage("[command:string]").AddAliases("h", "cmds", "commands"))

  bot.AddCommand(NewCommand("stats", "General", func(ctx *CommandContext) {
    // Runtime stats, these stats for some reason makes me feel really good.
    stats := &runtime.MemStats{}
    runtime.ReadMemStats(stats)
    // Counters
    var guilds, users, channels int
    guilds = len(ctx.Session.State.Guilds)
    for _, guild := range ctx.Session.State.Guilds {
      users += guild.MemberCount
      channels += len(guild.Channels)
    }

    ctx.BuildEmbed(NewEmbed().
      SetTitle("Stats").
      SetAuthor(ctx.Session.State.User.Username, ctx.Session.State.User.AvatarURL("256")).
      SetColor(COLOR).
      AddField("Go Version", strings.TrimPrefix(runtime.Version(), "go")).
      AddField("DiscordGo Version", discordgo.VERSION).
      AddField("Sapphire Version", VERSION).
      AddField("Bot Stats", fmt.Sprintf("**Guilds:** %d\n**Users:** %d\n**Channels:** %d\n**Uptime:** %s", guilds, users, channels, humanize.RelTime(bot.Uptime, time.Now(), "", ""))).
      AddField("Command Stats", fmt.Sprintf("**Total Commands:** %d\n**Commands Ran:** %d", len(bot.Commands), bot.CommandsRan)).
      AddField("Memory Stats", fmt.Sprintf("**Used:** %s / %s\n**Garbage Collected:** %s\n**GC Cycles:** %d\n**Forced GC Cycles:** %d\n**Last GC:** %s\n**Next GC Target:** %s\n**Goroutines:** %d",
        humanize.Bytes(stats.Alloc),
        humanize.Bytes(stats.Sys),
        humanize.Bytes(stats.TotalAlloc - stats.Alloc),
        stats.NumGC,
        stats.NumForcedGC,
        humanize.Time(time.Unix(0, int64(stats.LastGC))),
        humanize.Bytes(stats.NextGC),
        runtime.NumGoroutine(),
      )).
      AddField("Technical Info", fmt.Sprintf("**CPU Cores:** %d\n**OS/Arch:** %s/%s",
        runtime.NumCPU(),
        runtime.GOOS,
        runtime.GOARCH,
      )))
  }).SetDescription("Stats for nerds.").AddAliases("botstats", "info"))

  bot.AddCommand(NewCommand("invite", "General", func(ctx *CommandContext) {
    ctx.ReplyLocale("COMMAND_INVITE", fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&permissions=%d&scope=bot",
      ctx.Session.State.User.ID, bot.InvitePerms))
  }).SetDescription("Invite me to your server!").AddAliases("inv"))

  bot.AddCommand(NewCommand("enable", "Owner", func(ctx *CommandContext) {
    command := ctx.Bot.GetCommand(ctx.Arg(0).AsString())
    if command == nil {
      ctx.ReplyLocale("COMMAND_NOT_FOUND", ctx.Arg(0))
      return
    }
    if command.Enabled {
      ctx.ReplyLocale("COMMAND_ENABLE_ALREADY")
      return
    }
    command.Enable()
    ctx.ReplyLocale("COMMAND_ENABLE_SUCCESS", ctx.Arg(0))
  }).SetDescription("Enables a disabled command.").SetOwnerOnly(true).SetUsage("<command:string>"))

  bot.AddCommand(NewCommand("disable", "Owner", func(ctx *CommandContext) {
    command := ctx.Bot.GetCommand(ctx.Arg(0).AsString())
    if command == nil {
      ctx.ReplyLocale("COMMAND_NOT_FOUND", ctx.Arg(0).AsString())
      return
    }

    if !command.Enabled {
      ctx.ReplyLocale("COMMAND_DISABLE_ALREADY")
      return
    }
    command.Disable()
    ctx.ReplyLocale("COMMAND_DISABLE_SUCCESS", ctx.Arg(0).AsString())
  }).SetDescription("Disables an enabled command.").SetOwnerOnly(true).SetUsage("<command:string>"))

  bot.AddCommand(NewCommand("gc", "Owner", func(ctx *CommandContext) {
    before := &runtime.MemStats{}
    runtime.ReadMemStats(before)
    // Additionally we will collect extra garbage by freeing these stuff aswell, since this command is meant to be ran
    // in memory critical situations losing them doesn't hurt at all.
    bot.CommandCooldowns = make(map[string]map[string]time.Time)
    bot.CommandEdits = make(map[string]string)
    runtime.GC()
    after := &runtime.MemStats{}
    runtime.ReadMemStats(after)
    ctx.Reply("Forced Garbage Collection.\n  - Freed **%s**\n  - %d Objects Collected.\n  - Took **%d**μs",
      humanize.Bytes(before.Alloc - after.Alloc), after.Frees - before.Frees, after.PauseTotalNs - before.PauseTotalNs)
  }).SetDescription("Forces a garbage collection cycle.").AddAliases("garbagecollect", "forcegc", "runtime.GC()").SetOwnerOnly(true))
  return bot
}
