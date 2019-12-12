# Sending Embeds in Sapphire
Sending embeds is done via the `ctx.ReplyEmbed` that takes a regular discordgo embed, nothing too special.

But sapphire offers an easier way to construct embeds with a [discord.js](https://discord.js.org) inspired embed builder.

```go
func SendEmbed(ctx *sapphire.CommandContext) {
  ctx.ReplyEmbed(sapphire.NewEmbed().SetTitle("This is an embed!").SetColor(0xFF0000).SetDescription("Hello, World!").Build())
}
```
Easy! To avoid needing `.Build()` we offer an alternative method `BuildEmbed` that calls it for you.
```go
func SendEmbed(ctx *sapphire.CommandContext) {
  ctx.BuildEmbed(sapphire.NewEmbed().SetTitle("This is an embed!").SetColor(0xFF0000).SetDescription("Hello, World!"))
}
```
The embed builder also takes in account embed limits, so if you ever accidentally go over the limit the builder will truncate them for you!
