# Introduction to Sapphire
Sapphire is a powerful framework to build discord bots in Golang with a lot of tools you will need.

This guide will teach you the fundamentals of using sapphire, check out the [documentation](https://godoc.org/github.com/sapphire-cord/sapphire) afterwards to learn more.

## Basic Bot
```go
package main

import (
  "github.com/sapphire-cord/sapphire"
  "github.com/bwmarrin/discordgo"
)

func main() {
  dg, _ := discordgo.New("token")
  bot := sapphire.New(dg)
  bot.SetPrefix("!")
  bot.LoadBuiltins()
  bot.Connect()
  bot.Wait()
}
```
That's as basic as it can get.

- First we create a `*sapphire.Bot` with `sapphire.New` passing in our `*discordgo.Session`
- Next we set a prefix with `bot.SetPrefix` (the default is `!`)
- Next we load the builtin commands including `help` with `bot.LoadBuiltins()`
- Finally we connect the bot to discord using `bot.Connect()`,

> **Note:** As said in discordgo's documentation, you must prefix the token with `Bot` for bot accounts.

If you need dynamic prefixes you can also supply a function that is called everytime sapphire needs the prefix
```go
bot.SetPrefixHandler(func(bot *sapphire.Bot, msg *discordgo.Message, dm bool) string {
  // Call database here, etc and return the prefix string
  // dm is true if this command is invoked inside a DM
  return "!"
})
```

Sapphire's APIs is also chainable so you can do it in a fancy way
```go
sapphire.New(dg).SetPrefix("!").LoadBuiltins().Connect().Wait()
```

Next [let's write some commands](Commands.md)
