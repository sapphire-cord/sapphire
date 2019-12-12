# Command Flags
Sapphire allows optional command flags in every command invokation.

It's very simple when running a command sapphire parses the flags in the form `--flag=value` value is optional so `--flag` also works, these type of arguments don't appear as arguments they are stripped of the message before the command handler processes the message.

These flags can be accessed via `ctx.Flag`/`ctx.HasFlag`

- `ctx.HasFlag(name)` checks whether the flag `name` is passed to this message, useful for boolean flags.
- `ctx.Flag(name)` returns the flag `name`'s value or an empty string "" if the flag didn't have a value or not specified.

You don't need to pass the `--` to these functions.
