# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`nyaa-fstop` is a site plugin for [`github.com/RangelReale/filesharetop`](https://github.com/RangelReale/filesharetop), the shared core that handles MongoDB storage, scoring and the web UI. This repo only contains the scraper for http://www.nyaa.se/?cats=... and two thin `main` binaries. Almost all behavior beyond HTML parsing lives in the core packages `fstoplib` (`lib/`), `fstopimp` (`importer/`) and `fstopsite` (`site/`). Read the core repo's CLAUDE.md for the data model and scoring.

This is legacy, pre-modules Go code. There is no `go.mod`, and imports assume GOPATH layout (`$GOPATH/src/github.com/RangelReale/nyaa-fstop`, with `filesharetop` next to it). Dependencies include `github.com/PuerkitoBio/goquery`, `gopkg.in/mgo.v2`, `github.com/BurntSushi/toml` and, through the core site package, `code.google.com/p/plotinum`.

## Layout

- Root package `nyaa`: `fetcher.go` implements `fstoplib.Fetcher` (`ID`, `SetLogger`, `Fetch`, `CategoryMap`). `nyaaparser.go` (`NYParser`) downloads pages and parses them with goquery into `map[id]*fstoplib.Item`.
- `nyaa-importer/`: run periodically (hourly). It fetches, calls `Importer.Import`, then `Consolidate("", 48)` and `Consolidate("weekly", 168)` into database `fstop_nyaa`.
- `nyaa-site/`: runs `fstopsite.RunServer` on port **13112** against `fstop_nyaa`, with `TopId = "weekly"`.
- Both binaries connect to MongoDB at `localhost` (hard-coded) and support `-version`.

## Scraping

For each raw category ID in `CategoryMap` (e.g. `1_37`), it scrapes 1 page sorted by seeders and 1 sorted by leechers. The parser reads each row's category from the page itself.

When editing the parser:
- Create items with `fstoplib.NewItem()` so that unavailable stats stay `-1`.
- Dedupe on the item ID. The same torrent can appear in both the seeders pass and the leechers pass. `SeedersPos` / `LeechersPos` record its rank in each pass.
- `AddDate` must be `YYYY-MM-DD`.
- `Link` must be an absolute URL.
- Rows that fail to parse are logged and skipped, not treated as errors.

## Configuration

None: there is no `Config` type, and `NewFetcher()` takes no arguments.

## Commands

```sh
go build ./nyaa-importer ./nyaa-site
go test -run TestFetcher .     # single test
```

`TestFetcher` hits the live site.
