# sports / bsfeeds

[BS Feeds](https://bsfeeds.com) is a Bluesky feed generator
providing curated sports feeds.

Available feeds:

- [Golf](https://bsky.app/profile/bsfeeds.com/feed/golf)
- [MLB](https://bsky.app/profile/bsfeeds.com/feed/mlb)
- [NBA](https://bsky.app/profile/bsfeeds.com/feed/nba)
- [NCAAB](https://bsky.app/profile/bsfeeds.com/feed/ncaab)
- [NCAAF](https://bsky.app/profile/bsfeeds.com/feed/ncaaf)
- [NFL](https://bsky.app/profile/bsfeeds.com/feed/nfl)
- [NHL](https://bsky.app/profile/bsfeeds.com/feed/nhl)
- [PWHL](https://bsky.app/profile/bsfeeds.com/feed/pwhl)

## Architecture

I wrote the server as a single Go binary that connects to
the Bluesky firehose and a Postgres database.

It contains a few services:

- **Web**: HTTP handlers for
  [Bluesky's feed generator protocol](https://docs.bsky.app/docs/starter-templates/custom-feeds)
- **Stream**: subscribes to the Bluesky firehose via WebSocket,
  receiving every public post in real-time
- **Builder**: matches incoming posts against feed definitions,
  writes matches to Postgres
- **Maintenance**: prunes feed items older than one week
  to keep the database small

It has no message queues or worker processes.
Goroutines handle concurrency within the binary.

## Feed definitions

Each feed is defined in Go code with:

- **Text terms**: team names, player names, playoff terms
- **Account DIDs**: official team accounts, beat reporters
- **Game patterns**: regex for matchups like "Bills vs Cowboys"

For example, the NFL feed includes all 32 team names,
hundreds of player names, and DIDs for accounts like
`@detroitlions.bsky.social`.

Terms are chosen carefully to avoid false positives.
Ambiguous names are excluded. Substring conflicts are tested
(e.g., "Josh Allen" shouldn't match posts about a different Josh Allen).

I use [Warp](/ai/warp) to help me update content periodically as rosters change.

## Matching

When a post arrives from the firehose, the matcher checks:

1. Does the post contain any text term (case-insensitive)?
2. Does it match a game pattern like "Cowboys-Bills"?
3. Is the author a listed account DID?

Matches are written to Postgres with the feed ID and post URI.
When Bluesky clients request a feed, the web handler
queries recent matches and returns them in
[AT Protocol](https://atproto.com/) format.
