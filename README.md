# manifest-inspector

A CLI testing tool for inspecting streaming manifests (HLS today, DASH/MPD later).

**This project is not complete.** It is an early, ongoing experiment. APIs, flags, and checks will change. Do not treat it as a finished validator.

## What it does right now

Given an HLS master playlist URL, `manifest-inspector` fetches the manifest, parses `#EXT-X-STREAM-INF` variants and `#EXT-X-MEDIA` renditions, then prints a pass/fail report.

Current checks include:

- Master playlist has at least one variant
- Variants are sorted by `BANDWIDTH` (warning if not)
- `CODECS` is present on variants (warning if missing)
- Audio and subtitle `GROUP-ID` values are referenced by variants
- `DEFAULT` is `YES` or `NO`
- `STABLE-RENDITION-ID` matches the allowed character set

If the playlist is valid, it prints variant count and bitrate range. If not, it lists issues by severity (`ERROR`, `WARNING`, `INFO`).

## What is not done

- DASH / MPD support is not implemented (the `--url` flag mentions it, but only HLS `m3u8` URLs are handled)
- Media playlists (segment lists) are not inspected
- No JSON output, config file, or test suite yet
- Many HLS tags and edge cases are still missing

## Requirements

- Go 1.27+

## Usage

```bash
go build -o manifest_inspector
./manifest_inspector --url https://example.com/master.m3u8
```

Or via Make:

```bash
make run url=https://example.com/master.m3u8
```

Short flag: `-u` is the same as `--url`.

## Status

Work in progress. Expect incomplete coverage, rough edges, and breaking changes.
