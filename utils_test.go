package sapphire

import (
  "testing"
)

func TestEscape(t *testing.T) {
  everyone := Escape("hello @everyone")
  here := Escape("come @here")
  if everyone != "hello @\u200beveryone" {
    t.Error("Escape din't return the expected output for @everyone")
  }

  if here != "come @\u200bhere" {
    t.Error("Escape didn't return the expectd output for @here")
  }
}
