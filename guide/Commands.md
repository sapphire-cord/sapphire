# Sapphire Commands
Let's write some commands.

We will start with a basic ping pong example

> **Note:** There is already a ping command if you use LoadBuiltins but for the sake of the tutorial we will override it to see how commands are made.

Commands are created with `sapphire.NewCommand` and is setup with the command APIs and finally registered via `bot.AddCommand()`

Commands aren't restricted to a specific directory structure or package but it is neater to follow the structure we did here.

Create a `commands/` folder and create a subdirectory for each category you need, e.g (`commands/general`, `commands/moderation`)

We will put ping command in `commands/general/ping.go`

Ready to start coding? let's go.
```go
package general // our category package

import (
  "github.com/sapphire-cord/sapphire"
)

func Ping(ctx *sapphire.CommandContext) {
  ctx.Reply("Pong!")
}
```
As easy as that!

But wait, we have to first register them to be able to use it!

To do this we create a `commands/init.go` in our project to initialize the commands, the code looks like
```go
package commands

import (
  "github.com/sapphire-cord/sapphire"
  "...your_project.../commands/general"
  // Repeat importing any other categories you have.
)

func Init(bot *sapphire.Bot) {
  bot.AddCommand(sapphire.NewCommand("ping", "General", general.Ping).SetDescription("Pong!"))
  // Repeat for any additional commands you want to add.
}
```
`NewCommand` takes `name`/`category`/`CommandHandler` and returns the `*Command` it has chainable APIs to set properties easily.

Finally we call the init function we created in our main entry file where we initialized our sapphire.Bot, make sure you import the commands package and do
```go
commands.Init(bot)
```
And that's it we are ready to run our `!ping` in chat.

**But ugh i don't want to register every possible commands there, can't i get autoloading or something?** That is how Go works, it compiles to a single binary and loses the ability to understand Go source so we can't dynamically load commands at runtime, however we can dynamically generate the registration code before runtime and we made a tool for it! Meet [spgen](SPGen.md)

Next [let's see how to use arguments](Arguments.md)
