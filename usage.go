package sapphire

import (
  "errors"
  "strings"
  "regexp"
)

type UsageTag struct {
  Name string // Name of the tag, e.g for <reason:string> the name is reason.
  Type string // Type of the tag, e.g for <reason:string> the type is string.
  Rest bool // If this is rest of the arguments, e.g for <reason:string...> it is true.
  Required bool // If this argument is required, e.g <name> is required but [name] is not.
}

// Parse a usage string into tags.
func ParseUsage(usage string) ([]*UsageTag, error) {
  // TODO: We'll need to handle more cases to improve error handling.
  tags := make([]*UsageTag, 0)
  // The current tag we are parsing and building.
  current := &UsageTag{Required:false, Rest:false, Type:"", Name:""}
  // true if we are currently parsing the type, otherwise the name.
  typeMode := false
  for _, c := range usage {
    if c == '<' {
      current.Required = true
    } else if c == '>' {
      // Beginning the tag name with @@ is a syntactic sugar for member and beginning with @ is for user.
      // And if the type is missing we assume literal.
      if strings.HasPrefix(current.Name, "@@") && current.Type == "" {
        current.Name = strings.TrimPrefix(current.Name, "@@")
        current.Type = "member"
      } else if strings.HasPrefix(current.Name, "@") && current.Type == "" {
        current.Name = strings.TrimPrefix(current.Name, "@")
        current.Type = "user"
      } else if current.Type == "" { current.Type = "literal" }
      tags = append(tags, current)
      current = &UsageTag{Required:false, Rest:false, Type:"", Name:""}
      typeMode = false
    } else if c == '[' {
      if current.Required {
        return tags, errors.New("Cannot open an optional tag after opening a required one.")
      }
    } else if c == ']' {
      if strings.HasPrefix(current.Name, "@@") && current.Type == "" {
        current.Name = strings.TrimPrefix(current.Name, "@@")
        current.Type = "member"
      } else if strings.HasPrefix(current.Name, "@") && current.Type == "" {
        current.Name = strings.TrimPrefix(current.Name, "@")
        current.Type = "user"
      } else if current.Type == "" { current.Type = "literal" }
      tags = append(tags, current)
      current = &UsageTag{Required:false, Rest:false, Type:"", Name:""}
      typeMode = false
    } else if c == ' ' {
      continue
    } else if c == ':' {
      typeMode = true
    } else {
      if typeMode { current.Type += string(c) } else { current.Name += string(c) }
    }
  }
  // Now that we know enough about the tags and how many are there we can validate rest args.
  for i, tag := range tags {
    if strings.HasSuffix(tag.Type, "...") {
      if i != len(tags) - 1 {
        return tags, errors.New("Rest parameters can only appear last.")
      }
      tag.Type = strings.TrimSuffix(tag.Type, "...")
      tag.Rest = true
    }
  }
  return tags, nil
}

// HumanizeUsageRegex is the regexp used for HuamnizeUsage
var HumanizeUsageRegex = regexp.MustCompile("(<|\\[)(\\w+):[^.]+?(\\.\\.\\.)?(>|\\])")

// HumanizeUsage removes the unneccessary types and shows only the names.
// e.g <hello:string> <user:user> [rest:int...] => <hello> <user> [rest...]
func HumanizeUsage(usage string) string {
  return HumanizeUsageRegex.ReplaceAllString(usage, "$1$2$3$4")
}
