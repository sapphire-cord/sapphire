# Localizing your sapphire bot
Localization is easy with sapphire and we will show it all here.

Inside our commands we can simply use `ctx.ReplyLocale("key")` where key is a key in our language components.

We have to tell the locale we are going to use to sapphire with `SetLocaleHandler`
```go
bot.SetLocaleHandler(func(bot *sapphire.Bot, msg *discordgo.Message, dm bool) string {
  // The return is a string showing a language component name to use.
  // The default localization for English is called "en-US" but you can use any name for your own components.
  // Here we hardcode the language to always be English which is the default behaviour
  // But you can for example fetch the language settings of the guild from a database and return that name.
  return "en-US"
})
```
If you just want to hardcode your bot to a specific locale you can also use `bot.SetLocale("en-US")` but that defeats the purpose of localization.

Full example, let's write an actual language and see it in action, we will create a `hello` command that says hello in different languages.

The code for the hello command is
```go
package general

import (
  "github.com/sapphire-cord/sapphire"
)

func Hello(ctx *sapphire.CommandContext) {
  ctx.ReplyLocale("COMMAND_HELLO")
}
```
We assume you've already read the [Commands guide](Commands.md) and you can register that command.

Now if we run it as is the bot will reply with `No localization found for the key "COMMAND_HELLO" Please report this to the developers.` We are the developers so let's fix it, it basically says that the key isn't localized (yet)

And here we go again with creating a new component, make a `languages/` directory and put `languages/english.go`

```go
package languages

import (
  "github.com/sapphire-cord/sapphire"
)

var English = sapphire.NewLanguage("en-US").
  Set("COMMAND_HELLO", "Hello").
  Merge(sapphire.English)
```
We create a new language with name `en-US` (which will result is overriding the builtin one) and set the key `COMMAND_HELLO` to say `Hello` in English, finally we also merge the builtin English language component so the builtins don't break as we are overriding the default one. (Note if you want to change the builtin responses use Merge first and use Set to overwrite a specific key.)

Ah you know the grind now let's make a `languages/init.go`
```go
package languages

import (
  "github.com/sapphire-cord/sapphire"
)

func Init(bot *sapphire.Bot) {
  bot.AddLanguage(English)
  // Repeat for any other languages you want to add.
}
```
And don't forget to import and call `languages.Init(bot)` in your entry file.

Now run `!hello` again and, it just says "Hello" not that interesting but that's about to change.

Let's add a french localization, create `languages/french.go`
```go
package languages

import (
  "github.com/sapphire-cord/sapphire"
)

var French = sapphire.NewLanguage("fr-FR").
  Set("COMMAND_HELLO", "Bonjour")
```
sapphire's builtins currently don't have localizations for other languages apart from English so no merging is needed but the builtins will default to reply in English if you don't translate them. (Look at `language.go` in the source for the keys you can translate.)

Now hardcode the bot for a moment to speak french `bot.SetLocale("fr-FR")` and run `!hello` again and it responds in French!

When the bot can't find a key it fallbacks to the default languages and if it can't find it in the default language it replies with what we have seen before adding the localized key. To set the default languages use `bot.SetDefaultLocale("fr-FR")` now the bot speaks french when it can't find a key in the set locale.

### Locale arguments
You won't always send constant strings, sometimes you need to insert some dynamic info calculated from the command, to do this we allow language keys to have format strings and ReplyLocale can take extra args to format them, just like printf.

Next [let's send embeds in a fancy way](Embeds.md)
