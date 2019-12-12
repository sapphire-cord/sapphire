# Command Arguments in Sapphire
Arguments in sapphire commands is easy, it is done via usage strings which tells the required and optional arguments and their types.

For example this can be the usage string for a kick command.
```
<member:member> [reason:string...]
```
it takes a required member, and an optional reason, the `...` shows that it is a rest argument, as in it parses whatever else after it using the `reason` part of usage. (`...` can only appear on the last tag.)

The name after a colon `:` is the type, here we want a member from the server it's ran on so we can kick them.

Here is how we can access the arguments for that usage string.

```go
func Kick(ctx *sapphire.CommandContext) {
  user := ctx.Arg(0).AsUser()
  reason := ctx.JoinedArgs(1)
  // user is the first argument, it will always exist.
  // reason is optional and is an empty string "" if not provided
}
```
Additionally usage strings gives a human readable clue to the user on how to use the command, it gets documented in help.

You can access the raw arguments via the `ctx.RawArgs` slice that doesn't follow usage strings, and you can join all the raw arguments with a space via `ctx.JoinedArgs`, see the documentation for more details.

Usage string functionality is still in it's early stages, it works but is less powerful, in the future we have plans to add multiple args `<add|remove>` etc.

Currently the following types are supported, more will be added and suggestions are welcome:
- `int`/`num`/`number` - A number like `5`
- `string`/`str` - A string or text input.
- `user` - A user on discord, searches globally from all guilds.
- `member` A member from the current guild the command is ran on.

**TODO** These are types are planned to be added, check this before suggesting, contributions are welcome.
- `server`/`guild` - A Discord server
- `codeblock`/`code` parses a codeblock's contents.

When an argument is required sapphire will take care that it is provided so you can just assume it always exists.

For optionals if it isn't provided it's ignored and returns an empty arg which you can call `IsProvided` on it to check for existence, as soon as an optional is provided it gets treated like a required one and won't pass until it is parsed successfully without errors.

An example to check argument existence:
```go
arg := ctx.Arg(0)
var user *discordgo.User
if arg.IsProvided() {
  user = arg.AsUser()
} else {
  // handle non-existence, e.g a default
  user = ctx.Author
}
```

Additionally for the user and member types there is an alias to make it easier, `@user` is same as `user:user` and `@@member` is the same as `member:member`

Also you must be very aware what `As*` cast functions you are calling, it must be what you defined in the usage string because it casts blindly and assumes the argument is present as said in usage string, failing to do so can lead to panics.

Next [let's try monitors](Monitors.md)
