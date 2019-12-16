package sapphire

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "io"
  "strings"
)

type CommandHandler func(ctx *CommandContext)

// Command represents a command in the sapphire framework.
type Command struct {
  Name string // The command's name. (default: required)
  Aliases []string // Aliases that point to this command. (default: [])
  Run CommandHandler // The handler that actually runs the command. (default: required)
  Enabled bool // Wether this command is enabled. (default: true)
  Description string // The command's brief description. (default: "No Description Provided.")
  Category string // The category this command belongs to. (default: required)
  OwnerOnly bool // Wether this command can only be used by the owner. (default: false)
  GuildOnly bool // Wether this command can only be ran on a guild. (default: false)
  UsageString string // Usage string for this command. (default: "")
  Usage []*UsageTag // Parsed usage tags for this command.
  Cooldown int // Command cooldown in seconds. (default: 0)
  Editable bool // Wether this command's response will be editable. (default: true)
  RequiredPermissions int // Permissions the user needs to run this command. (default: 0)
  BotPermissions int // Permissions the bot needs to perform this command. (default: 0)
}

func NewCommand(name string, category string, run CommandHandler) *Command {
  return &Command{
    Name: name,
    Category: category,
    Run: run,
    Aliases: []string{},
    Enabled: true,
    Description: "No Description Provided.",
    OwnerOnly: false,
    GuildOnly: false,
    UsageString: "",
    Editable: true,
    Cooldown: 0,
    RequiredPermissions: 0,
    BotPermissions: 0,
    Usage: make([]*UsageTag, 0),
  }
}

// AddAliases adds aliases to this command.
func (c *Command) AddAliases(aliases ...string) *Command {
  c.Aliases = append(c.Aliases, aliases...)
  return c
}

// SetDescription sets the command's description
func (c *Command) SetDescription(description string) *Command {
  c.Description = description
  return c
}

// SetUsage sets the usage string for this command.
// Panics if there is a parse error in the usage string.
func (c *Command) SetUsage(usage string) *Command {
  c.UsageString = usage
  usg, err := ParseUsage(usage)
  if err != nil {
    panic(err)
  }
  c.Usage = usg
  return c
}

// Disable disables the command.
func (c *Command) Disable() *Command {
  c.Enabled = false
  return c
}

// Enable enables the command.
func (c *Command) Enable() *Command {
  c.Enabled = true
  return c
}

// SetOwnerOnly toggles wether the command can only be used by the bot owner.
func (c *Command) SetOwnerOnly(toggle bool) *Command {
  c.OwnerOnly = toggle
  return c
}

// SetGuildOnly toggles if this command can only be used on a guild.
func (c *Command) SetGuildOnly(toggle bool) *Command {
  c.GuildOnly = toggle
  return c
}

// SetEditable toggles wether this command will be respondable to edits.
func (c *Command) SetEditable(toggle bool) *Command {
  c.Editable = toggle
  return c
}

// SetCooldown sets the command's cooldown in seconds.
func (c *Command) SetCooldown(cooldown int) *Command {
  c.Cooldown = cooldown
  return c
}

// CommandContext represents an execution context of a command.
type CommandContext struct {
  Command *Command // The currently executing command.
  Message *discordgo.Message // The message of this command.
  Session *discordgo.Session // The discordgo session.
  Bot *Bot // The sapphire Bot.
  Channel *discordgo.Channel // The channel this command was ran on.
  Author *discordgo.User // Alias of Context.Message.Author
  Args []*Argument // List of arguments.
  Prefix string // The prefix used to invoke this command.
  Guild *discordgo.Guild // The guild this command was ran on.
  Flags map[string]string // Map of flags passed to the command. e.g --flag=yo
  Locale *Language // The current language.
  RawArgs []string // The raw args that may not match the usage string.
  InvokedName string // The name this command was invoked as, this includes the used alias.
}

// Reply replies with a string.
// It will call Sprintf() on the content if atleast one vararg is passed.
func (ctx *CommandContext) Reply(content string, args ...interface{}) (*discordgo.Message, error) {
  if !ctx.Command.Editable {
    return ctx.ReplyNoEdit(content)
  }

  // This is neccessary to avoid problems with dynamic content
  // ctx.Reply(dynamicVariable)
  // If the user doesn't intend to use the formatting then don't use Sprintf
  // Because it will mess up any '%s' etc in the content even if the user did not intent to format it.
  // Another solution is a seperate Replyf function which was my original solution
  // But i think it's cleaner to stick with one function.
  if len(args) > 0 {
    content = fmt.Sprintf(content, args...)
  }

  m, ok := ctx.Bot.CommandEdits[ctx.Message.ID]
  if !ok {
    msg, err := ctx.Session.ChannelMessageSend(ctx.Channel.ID, content)
    if err != nil {
      return nil, err
    }
    ctx.Bot.CommandEdits[ctx.Message.ID] = msg.ID
    return msg, nil
  }
  return ctx.Session.ChannelMessageEditComplex(discordgo.NewMessageEdit(ctx.Channel.ID, m).
    SetContent(content))
}

// ReplyNoEdit replies with content but does not consider editable option of the command.
func (ctx *CommandContext) ReplyNoEdit(content string, args ...interface{}) (*discordgo.Message, error) {
  // See the comments in Reply
  if len(args) > 0 {
    content = fmt.Sprintf(content, args...)
  }
  return ctx.Session.ChannelMessageSend(ctx.Channel.ID, content)
}

// ReplyLocale sends a localized key for the current context's locale.
func (ctx *CommandContext) ReplyLocale(key string, args ...interface{}) (*discordgo.Message, error) {
  res := ctx.Locale.Get(key, args...)

  if res != "" {
    return ctx.Reply(res)
  }

  // Try the default locale.
  fallback := ctx.Bot.DefaultLocale.Get(key, args...)
  if fallback != "" {
    return ctx.Reply(fallback)
  }

  // All failed, the key isn't translated, report the error.
  // We have to also watch out if the error message isn't translated!
  return ctx.Reply(ctx.Locale.GetDefault("LOCALE_NO_KEY", key,
    ctx.Bot.DefaultLocale.GetDefault("LOCALE_NO_KEY", key,
      fmt.Sprintf("No localization found for the key \"%s\" Please report this to the developers.", key))))
}

// EditLocale edits msg with a localized key
func (ctx *CommandContext) EditLocale(msg *discordgo.Message, key string, args ...interface{}) (*discordgo.Message, error) {
  res := ctx.Locale.Get(key, args...)
  if res != "" {
    return ctx.Edit(msg, res)
  }
  fallback := ctx.Bot.DefaultLocale.Get(key, args...)
  if fallback != "" {
    return ctx.Edit(msg, fallback)
  }
  // All failed, the key isn't translated, report the error.                                                                              // We have to also watch out if the error message isn't translated!
  return ctx.Edit(msg, ctx.Locale.GetDefault("LOCALE_NO_KEY", key,
    ctx.Bot.DefaultLocale.GetDefault("LOCALE_NO_KEY", key,
      fmt.Sprintf("No localization found for the key \"%s\" Please report this to the developers.", key))))
}

// Edit edits msg's content
// It will call Sprintf() on the content if atleast one vararg is passed.
func (ctx *CommandContext) Edit(msg *discordgo.Message, content string, args ...interface{}) (*discordgo.Message, error) {
  // See the comments in Reply
  if len(args) > 0 {
    content = fmt.Sprintf(content, args...)
  }
  return ctx.Session.ChannelMessageEdit(msg.ChannelID, msg.ID, content)
}

// HasArgs returns true if there is atleast one argument in the raw args.
func (ctx *CommandContext) HasArgs() bool {
  return len(ctx.RawArgs) > 0
}

// ReplyEmbed replies with an embed.
func (ctx *CommandContext) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
  if !ctx.Command.Editable {
    return ctx.ReplyEmbedNoEdit(embed)
  }
  m, ok := ctx.Bot.CommandEdits[ctx.Message.ID]
  if !ok {
    msg, err := ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
    if err != nil {
      return nil, err
    }
    ctx.Bot.CommandEdits[ctx.Message.ID] = msg.ID
    return msg, nil
  }
  return ctx.Session.ChannelMessageEditComplex(discordgo.NewMessageEdit(ctx.Channel.ID, m).SetContent("").SetEmbed(embed))
}

// ReplyEmbedNoEdits replies with an embed but not considering the editable option of the command.
func (ctx *CommandContext) ReplyEmbedNoEdit(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
  return ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}

// BuildEmbed calls ReplyEmbed(embed.Build())
func (ctx *CommandContext) BuildEmbed(embed *Embed) (*discordgo.Message, error) {
  return ctx.ReplyEmbed(embed.Build())
}

// BuildEmbedNoEdit calls ReplyEmbedNoEdit(embed.Build())
func (ctx *CommandContext) BuildEmbedNoEdit(embed *Embed) (*discordgo.Message, error) {
  return ctx.ReplyEmbedNoEdit(embed.Build())
}

// SendFile sends a file with name
func (ctx *CommandContext) SendFile(name string, file io.Reader) (*discordgo.Message, error) {
  return ctx.Session.ChannelFileSend(ctx.Channel.ID, name, file)
}

// Error invokes the bot's error handler, see bot.SetErrorHandler
func (ctx *CommandContext) Error(err interface{}, args ...interface{}) {
  // We make err an interface so it can also be invoked standalone with error objects, etc.

  // See comments in Reply
  if len(args) > 0 {
    err = fmt.Sprintf(fmt.Sprint(err), args...)
  }
  ctx.Bot.ErrorHandler(ctx.Bot, err)
}

// Flag returns the value of a commmnd flag, if it is a bool-flag use HasFlag() instead.
func (ctx *CommandContext) Flag(flag string) string {
  str, ok := ctx.Flags[flag]
  if ok {
    return str
  }
  return ""
}

// HasFlag returns a bool of wether the flag exists.
func (ctx *CommandContext) HasFlag(flag string) bool {
  _, ok := ctx.Flags[flag]
  return ok
}

// Arg returns the argument at index idx, it returns an empty arg if the index doesn't exist, useful for optional arguments.
func (ctx *CommandContext) Arg(idx int) *Argument {
  if len(ctx.Args) > idx {
    return ctx.Args[idx]
  }
  return &Argument{provided:false}
}

// Get the joined arguments as a string
// If sliced is provided then arguments are sliced by that before joining
// Examples
// ["hello", "example", "test"]
// JoinedArgs() => "hello example test"
// JoinedArgs(1) => "example test"
// This uses the raw arguments so arguments of different types are also shown in their raw form.
// This also means invalid arguments are also included but strings can never be invalid so this is useful
// for getting the rest strings.
func (ctx *CommandContext) JoinedArgs(sliced ...int) string {
  var s int = 0
  if len(sliced) > 0 { s = sliced[0] }
  return strings.Join(ctx.RawArgs[s:], " ")
  return ""
}

// Parses the raw args and fills in ctx.Args and returns true on success and on failure it replies with the error and returns false
// This is called in the command handler to process the arguments, it shouldn't be used in normal code
// It is exported to allow modification of the command handler in your own bot and avoid this line from giving errors.
func (ctx *CommandContext) ParseArgs() bool {
  // Helper to get an index without panicking.
  safeGet := func(idx int) string {
    if len(ctx.RawArgs) > idx {
      return ctx.RawArgs[idx]
    }
    return ""
  }
  // If it doesn't need arguments we are done.
  if ctx.Command.UsageString == "" {
    return true
  }
  ctx.Args = make([]*Argument, len(ctx.Command.Usage))
  for i, tag := range ctx.Command.Usage {
    v := safeGet(i)
    if tag.Required && v == "" {
      ctx.Reply("The argument **%s** is required.", tag.Name)
      return false
    }
    if tag.Rest {
      cut := ctx.RawArgs[i:]
      for i, raw := range cut {
        arg, err := ParseArgument(ctx, tag, raw)
        if err != nil {
          ctx.Reply(err.Error())
          return false
        }
        if i > len(ctx.Args) - 1 {
          ctx.Args = append(ctx.Args, arg)
        } else {
          ctx.Args[i] = arg
        }
      }
    } else {
      arg, err := ParseArgument(ctx, tag, safeGet(i))
      if err != nil {
        ctx.Reply(err.Error())
        return false
      }
      ctx.Args[i] = arg
    }
  }
  return true
}

// User gets a user by id, returns nil if not found.
func (ctx *CommandContext) User(id string) *discordgo.User {
  for _, guild := range ctx.Session.State.Guilds {
    if member, err := ctx.Session.State.Member(guild.ID, id); err == nil {
      return member.User
    }
  }
  return nil
}

// Member gets a member by id from the current guild, returns nil if not found.
func (ctx *CommandContext) Member(id string) *discordgo.Member {
  if ctx.Guild == nil { return nil }
  member, err := ctx.Session.State.Member(ctx.Guild.ID, id)
  if err != nil {
    return nil
  }
  return member
}

// GetFirstMentionedUser returns the first user mentioned in the message.
func (ctx *CommandContext) GetFirstMentionedUser() *discordgo.User {
  if len(ctx.Message.Mentions) < 1 { return nil }
  return ctx.Message.Mentions[0]
}

func (ctx *CommandContext) CodeBlock(lang, content string, args ...interface{}) (*discordgo.Message, error) {
  if len(args) > 0 {
    content = fmt.Sprintf(content, args...)
  }
  return ctx.Reply("```%s\n%s```", lang, content)
}

// React adds the reaction emoji to the message that triggered the command.
func (ctx *CommandContext) React(emoji string) error {
  return ctx.Session.MessageReactionAdd(ctx.Channel.ID, ctx.Message.ID, emoji)
}
