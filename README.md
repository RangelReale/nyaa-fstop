# nyaa-fstop

[filesharetop](https://github.com/RangelReale/filesharetop) plugin for Nyaa (nyaa.se). It scrapes the tracker's top torrent lists and serves a ranked "Nyaa Top" website.

## Components

- **`nyaa-importer`** scrapes the tracker and stores a snapshot in MongoDB (database `fstop_nyaa`). It then recomputes the 48-hour and weekly rankings. Run it periodically, for example hourly.
- **`nyaa-site`** is the web UI, served on port **13112**.

Both binaries connect to MongoDB on the local host and accept `-version`.

## Configuration

None needed.

## Build

This is pre-modules Go code. Put this repo and `filesharetop` in a GOPATH workspace, then build:

```sh
go get github.com/RangelReale/nyaa-fstop/...
# or, from a GOPATH checkout:
go build ./nyaa-importer ./nyaa-site
```

## Run

```sh
# once per hour, e.g. from cron
./nyaa-importer/nyaa-importer

# web UI at http://localhost:13112/
./nyaa-site/nyaa-site
```

## Tests

```sh
go test -run TestFetcher .
```

The test scrapes the live tracker, so it needs network access.
