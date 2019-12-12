package sapphire

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "strconv"
  "regexp"
)

// ----- Argument casting -----

// Argument represents an argument, it has methods to grab the right type.
type Argument struct {
  value interface{}
  provided bool
}

// The methods do not check for errors and casts rightaway, because such validations are done at argument parsing time
// And the errors are reported and the command execution is aborted, so the type will always be the type the user asked for.
// But it can panic if the usage is specifying a different type than what the user used, in that case it's their fault.

// We need to cover as much as types as possible because it isn't easy for the user to extend these.

// Returns the argument as a string.
func (arg *Argument) AsString() string {
  return arg.value.(string)
}

func (arg *Argument) AsInt() int {
  return arg.value.(int)
}

func (arg *Argument) AsFloat() float64 {
  return arg.value.(float64)
}

// IsProvided checks if this argument is provided, for optional arguments you must use this before casting.
func (arg *Argument) IsProvided() bool {
  return arg.provided
}

func (arg *Argument) AsUser() *discordgo.User {
  return arg.value.(*discordgo.User)
}

func (arg *Argument) AsMember() *discordgo.Member {
  return arg.value.(*discordgo.Member)
}

func (arg *Argument) AsGuild() *discordgo.Guild {
  return arg.value.(*discordgo.Guild)
}

func (arg *Argument) AsRole() *discordgo.Role {
  return arg.value.(*discordgo.Role)
}

func (arg *Argument) AsBool() bool {
  return arg.value.(bool)
}

func (arg *Argument) AsMessage() *discordgo.Message {
  return arg.value.(*discordgo.Message)
}

// ----- Argument parsing -----

// quick helper so i don't repeat provided:true
func arg(val interface{}) *Argument {
  return &Argument{provided:true, value:val}
}

// The Regexp used for matching user mentions.
var MentionRegex = regexp.MustCompile("^(?:<@!?)?(\\d{17,19})>?$")

// Parses the raw argument as specified in tag in context of ctx
func ParseArgument(ctx *CommandContext, tag *UsageTag, raw string) (*Argument, error) {
  if raw == "" {
    return &Argument{provided:false}, nil
  }
  switch tag.Type {
  case "str":
    fallthrough
  case "string":
    return arg(raw), nil
  case "num":
    fallthrough
  case "number":
    fallthrough
  case "int":
    val, err := strconv.Atoi(raw)
    return arg(val), err
  case "member":
    match := MentionRegex.FindStringSubmatch(raw)
    if len(match) < 2 {
      return nil, fmt.Errorf("**%s** must be a valid member mention or ID.", tag.Name)
    }
    member := ctx.Member(match[1])
    if member == nil {
      return nil, fmt.Errorf("That member cannot be found in this server.")
    }
    return arg(member), nil
  case "user":
    match := MentionRegex.FindStringSubmatch(raw)
    if len(match) < 2 {
      return nil, fmt.Errorf("**%s** must be a valid user mention or ID.", tag.Name)
    }
    user := ctx.User(match[1])
    if user == nil {
      return nil, fmt.Errorf("That user cannot be found.")
    }
    return arg(user), nil
  case "literal":
    if raw != tag.Name {
      return nil, fmt.Errorf("Literal argument must be **%s**", tag.Name)
    }
    return arg(raw), nil
  default:
    return nil, fmt.Errorf("The argument type '%s' is invalid.", tag.Type)
  }
}
