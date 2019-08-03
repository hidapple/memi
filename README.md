# memi
`memi` is the very cute slack bot.

## Usage
Set following environment variables.
```text
- SLACK_TOKEN      # Slack API token for memi bot user
- KIBELA_TOKEN     # Kibela API token
- KIBELA_TEAM      # Your Kibela team name
- KIBELA_LINK_NOTE # Kibela note ID which memi update by link command
```

Run `memi`
```shell
./memi
```

## Commands

### `ping`
Test reachability of memi bot.
```
@memi ping
```

### `link`
Append markdown link to the Kibela note of $KIBELA_LINK_NOTE.
```
@memi link $URL
```
Title is optional.
```
@memi link $URL This is title
```
