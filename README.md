# Sapphire
Sapphire is a bot framework built for [discordgo](https://github.com/bwmarrin/discordgo)

Join our Discord at [discord.gg/ArwQrH4](https://discord.gg/ArwQrH4)

**Features:**
- Easy to configure, lot of configurations with sane defaults.
- Abstract, We don't force you to use a specific database instead we let you express your database of choice to us.
- Lightweight, Sapphire only depends on very minimal dependencies so you don't spend time and space pulling in dependencies.
- Full featured, Sapphire ain't a toy, it's a complete framework for your bot.
- Lot of tools! A lot of utilities to avoid reinventing the wheel such as a reaction paginator and many more.
- Components can be disabled/enabled on the go at runtime.
- Localization, Sapphire helps to translate your bot's responses easily.

## Install
```sh
$ go get github.com/sapphire-cord/sapphire
```

## Usage
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
  bot.LoadBuiltins() // Loads builtin commands.
  bot.Connect()
  bot.Wait() // Needed to keep the process running.
}
```
That's it! a basic bot will be launched with builtin commands such as help.

Read our [guide](guide/) for lot of more cool things you can do! (Don't be afraid our guides are easy to follow and we are open for questions.)

See also the [documentation](https://godoc.org/github.com/pollen5/sapphire) after the guides for even more possibilies!

## Contributing
Sapphire is still in it's early stages of development and there is a lot of things that can be done, we welcome contributions on everything, typo-fixes, grammar-fixes, detail improvement, new guides and contributions on the code are all welcome.

Here is a little personal TODO for myself but you can help me with if you wish so.
- Use mutexes where needed.
- Improve the arguments API.

It is incomplete but fairly usable.

## License
[MIT](LICENSE)
