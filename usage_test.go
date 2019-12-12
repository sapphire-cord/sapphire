package sapphire

import (
  "testing"
  "fmt"
)

func TestParseUsage(t *testing.T) {
  assert := func(tag *UsageTag, name, typ string, required, rest bool) {
    if tag.Name != name {
      t.Errorf("Expected tag name to be %s but got %s", name, tag.Name)
    }
    if tag.Type != typ {
      t.Errorf("Expected tag type to be %s but got %s", typ, tag.Type)
    }
    if tag.Required != required {
      t.Errorf("Expected tag required to be %s but got %s", fmt.Sprint(required), fmt.Sprint(tag.Required))
    }
    if tag.Rest != rest {
      t.Errorf("Expected tag rest to be %s but got %s", fmt.Sprint(rest), fmt.Sprint(tag.Rest))
    }
  }
  tags, err := ParseUsage("<name:user> <reason:string> [something] <@user> <@@member> [days:int...]")
  if err != nil {
    t.Error(err)
  } else {
    assert(tags[0], "name", "user", true, false)
    assert(tags[1], "reason", "string", true, false)
    assert(tags[2], "something", "literal", false, false)
    assert(tags[3], "user", "user", true, false)
    assert(tags[4], "member", "member", true, false)
    assert(tags[5], "days", "int", false, true)
  }

  tag := "<hello:string> <name:string> <rest:int...>"
  expect := "<hello> <name> <rest...>"
  if res := HumanizeUsage(tag); res != expect {
    t.Errorf("Expected HumanizeUsage(\"%s\") to return \"%s\" but got \"%s\"", tag, expect, res)
  }
}
