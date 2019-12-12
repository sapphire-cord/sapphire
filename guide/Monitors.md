# Sapphire Monitors
Sapphire has monitors that is called everytime a message is sent, it allows multiple monitors to process the same message.

The command handling is also implemented as a monitor and is one of the monitors ran.

Each monitor is started in a seperate goroutine.

Monitors can be created via `sapphire.NewMonitor` and added via `bot.AddMonitor`

Example usecases for monitors would be a point system (track every messages and reward the user accordingly) and a word filter (check incoming messages if it contains a filtered/blacklisted word and act accordingly)

Monitors are components just like commands and the same rules apply for registering them and such.

Let's write a message logger as an example to try monitors. Create a `monitors/log.go` and put the following code
```go
package monitors

import (
  "github.com/sapphire-cord/sapphire"
  "github.com/bwmarrin/discordgo"
  "fmt" // or any logging library you use.
)

func Log(bot *sapphire.Bot, msg *discordgo.Message) {
  fmt.Println(msg.Content) // print content.
}
```
Create `monitors/init.go` as our entry for registering monitors.
```go
package monitors

import (
  "github.com/sapphire-cord/sapphire"
)

func Init(bot *sapphire.Bot) {
  bot.AddMonitor(sapphire.NewMonitor("logger", Log))
  // Repeat for any additional monitors you want to register.
}
```
By default monitors ignore bots and webhooks which is what you want in most cases but if you ever need to listen to those messages as well call `AllowBots`/`AllowWebhooks` in a monitor to configure them, e.g to allow our logger monitor to log webhooks and bots we would do```go
bot.AddMonitor(sapphire.NewMonitor("logger", Log).AllowBots().AllowWebhooks())
```

Finally in our main entry file where we connect our bot we make sure we load our monitors
```go
monitors.Init(bot)
```
Be sure to import the monitors package.

Next [let's try localizing our bot](Localization.md)

