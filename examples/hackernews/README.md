# Hacker News Reader

A native desktop Hacker News client built with Cute, the
Shoes-inspired UI toolkit for Rugo.

## Features

- **Five feeds** — browse Top, New, Best, Ask HN, and Show HN stories
- **Concurrent fetching** — 30 stories loaded in parallel via `spawn`
- **Open in browser** — double-click a story or press Enter to open the URL
- **HN comments** — Ctrl+Enter opens the HN comments page
- **Keyboard shortcuts** — Ctrl+R to refresh, Ctrl+Q to quit
- **HN-orange styling** — themed with Qt stylesheets

## Prerequisites

- Qt6 development libraries (`apt install qt6-base-dev` on Debian/Ubuntu)
- Go 1.22+
- Internet connection (fetches from the HN Firebase API)

## Build and run

```bash
rugo build main.rugo -o hackernews && ./hackernews
```

## Keyboard shortcuts

| Key | Action |
|-----|--------|
| `Ctrl+Q` | Quit |
| `Ctrl+R` | Refresh current feed |
| `Enter` | Open selected story in browser |
| `Ctrl+Enter` | Open HN comments page |

## How it works

The app uses Rugo's `spawn` for concurrent HTTP fetching — all 30 stories
are fetched in parallel, then displayed in a QListWidget. The `w.after()`
timer defers the initial load until after the window appears, so the UI
renders immediately.

```ruby
cute.app("Hacker News", 900, 640, fn(w)
  # ... build UI ...

  w.after(50, fn()
    load("top")
  end)
end)
```
