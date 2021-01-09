# Arke

*Arke is a no frills, discord bot with some basic functionality.*

## Commands

A discriminator is used to identify each command. The exclamation mark is serving this very purpose.

### `!ping || !pong`

Returns the string "Pong!" and "Ping!" respectively.

### `!serverinfo`

Returns some basic info of the joined guild (server).

### `!botinfo`

Returns some basic info of the bot.

### `!delete i`

Deletes the previous i messages sent in the current channel, where i is an integer. The upper limit is set to 100
messages by Discords API. If used without a number, it deletes the last message.

### `!whois s`

Returns some basic info about user s, where s is a string.

### `!mute s || !unmute s || !mute s i`

Mutes, unmutes or temporarily mutes user s, where s is a string and i is an integer.

### `!dog || !cat`

Returns a lovely dog or cat image from the web.
