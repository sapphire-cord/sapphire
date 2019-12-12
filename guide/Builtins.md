# Sapphire Builtins
Sapphire comes with a little basic commands builtin if you choose to load them. (`bot.LoadBuiltins()`)

This document explains all the builtins and what they do.

### Ping
When learning a new programming language the first thing you do is make a Hello world program, ping command is sort of that for when making discord bots, however sometimes it comes in handy if it can show the latency and the one in sapphire does.

### Help
One of the most must-have commands in Discord bots is a help command, it documents all available commands, sapphire's builtin help does just that in a clean style.

### Invite
If your bot is public then the invite command is one of the must have ones to allow people to invite it in their guilds. If your bot is not public then sapphire makes the invite command owner only.

### Enable/Disable
A command broke? A critical vulneribility found and you can't fix it right now? Fear not the disable builtin allows you to temporarily disable a command and likewise enable does the opposite and enables a disabled command.

### GC
GC triggers a cycle of garbage collection, this is useful for when your critically low on memory as it cleans some garbage to buy you some time.

## Overriding a builtin
Sometimes you may want to edit a command's behaviour, nothing suits everyone, so we tried to make that easy on you.

- If you just don't want that command to exist at all use `bot.RemoveCommand("name")`
- If you need to change responses you may want to edit the locale keys (see [Localization](Localization.md))
- If you want to modify the behaviour just a *little bit* copy the code from sapphire.go and place it up in a new command.
- If you want to do something else entirely with that name just go ahead sapphire overwrites existing commands.
