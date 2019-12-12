# Sapphire Generate (spgen)
`spgen` is a command line tool to generate the code that registers the commands, it's job is to find and identify all the commands and generate a source file that does the AddCommand calls, and the best part of it is that you can still express some command options at the source level and the generator will pass them to the command settings!

`spgen` is very useful for people coming from a background such as JavaScript or any other languages you could dynamically load commands, in Go that's not possible but `spgen` generates a source file that does that.

Install it via
```sh
$ go get github.com/sapphire-cord/spgen
```
It will be installed at `$GOPATH/bin` either add that to your PATH or invoke it via the path or however you want to do it.

Here's how it works
```go
package general

import (
  "github.com/sapphire-cord/sapphire"
)

// Responds with Pong!
// aliases: pong
func Ping(ctx *sapphire.CommandContext) {
  ctx.Reply("Pong!")
}
```
Assuming that file is in `commands/general/ping.go` and you are currently outside the commands folder run
```sh
$ spgen -import yourimport
```
And then it generates `commands/init.go` with the following
```go
// Package commands is the main entry point where all commands are registered.
// Auto-Generated with spgen DO NOT EDIT.
// To use this file import the package in your entry file and initialize it with commands.Init(bot)
package commands

import (
  "github.com/sapphire-cord/sapphire"
  "yourimport/commands/general"
)

func Init(bot *sapphire.Bot) {
  bot.AddCommand(sapphire.NewCommand("ping", "Responds with Pong!", general.Ping).AddAliases("pong"))
}
```
Cool huh? you just document your code like you do with godoc and it generates the initialization code.

It works by parsing the AST (Abstract Syntax Tree) and trying to find a go declaration that looks like `func CommandName(ctx *sapphire.Context)` `ctx` can have any name, and it also supports aliasing `sapphire` to something else!

The command's name is taken from the function's name lowercased.

Then you load that file in your entry file just like you did before manually, now all you need to do is run `spgen -import ...` everytime you add a new command, it takes care that you do not forget to register a command.

Here is all the supported options available via comments
- Literal `disabled` in a line: Disables the command, uses command.Disable() so you can still enable it at runtime.
- `Usage: usage string` a string to be shown in help.
- `Aliases: alias,*` comma separated aliases.

Additionally whitespaces doesn't matter and the keys are case insensitive.

You can also make the command usable by the owner only by prefixing the name with `Owner`
```go
func OwnerEval(ctx *sapphire.CommandContext) {}
```
The `Owner` part gets processed to turn the command owner only but is stripped out of the name so the result is an owner-only command with the name `eval`

## JSON
`spgen` is also powerful because not only it can generate a source file but it allows for various analysis on all of the commands found one of them is to generate JSON output, this could be used for API consumption etc.
```sh
$ spgen -json
```
It writes the json to the standard output (your console, terminal) but you can also save it to a file.
```sh
$ spgen -json -o commands.json
```
In the future we will expose the ability for users to write an analyser themselves to do various other things but for now only json is available.

Please keep in mind that `spgen` is new and still in development, it mostly works for usual cases but for some cases it could break, we covered as much as cases as possible but no code is perfect and it may not be suitable for all use cases, we are open to bug reports and feedback on how to improve the system.

Next [let's take a look at the builtins](Builtins.md)
