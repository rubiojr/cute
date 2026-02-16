# Threading

Cute is built on Qt, which has a strict threading rule: **widgets must only be accessed from the main thread**. Rugo's `spawn` creates goroutines that run on arbitrary OS threads, so touching widgets from a `spawn` block is unsafe.

Cute provides `cute.ui` to safely dispatch work back to the main thread.

## The Problem

This crashes or corrupts state:

```ruby
# ✗ WRONG — widget access from a background thread
spawn
  data = http.get("https://api.example.com/data")
  lbl.set_text(data.body)     # 💥 wrong thread!
end
```

## The Solution: cute.ui

`cute.ui` runs a callback on the Qt main thread and blocks until it completes:

```ruby
# ✓ CORRECT — fetch in background, update UI on main thread
spawn
  data = http.get("https://api.example.com/data")
  cute.ui(fn()
    lbl.set_text(data.body)   # safe: runs on main thread
  end)
end
```

Or with `do...end`:

```ruby
spawn
  data = http.get("https://api.example.com/data")
  cute.ui do
    lbl.set_text(data.body)
  end
end
```

## Pattern: Non-Blocking Load

### cute.fetch() — the easy way

`cute.fetch()` combines `spawn` + `cute.ui()` into a single call. The first function runs in the background; the second receives the result on the UI thread:

```ruby
cute.fetch(fn() http.get(url) end, fn(resp)
  lbl.set_text(resp.body)
end)
```

A more complete example with status feedback:

```ruby
status = cute.state("Ready")

load = fn(url)
  status.set("Loading...")
  cute.fetch(fn() http.get(url) end, fn(resp)
    process(resp)
    status.set("Done")
  end)
end
```

### Manual spawn + cute.ui

For full control, use the manual pattern:

```ruby
status = cute.state("Ready")

load = fn(url)
  status.set("Loading...")

  spawn
    result = http.get(url)
    cute.ui do
      process(result)
      status.set("Done")
    end
  end
end
```

1. `status.set("Loading...")` — runs on the main thread (inside a button callback), safe
2. `spawn` — starts background work, UI stays responsive
3. `http.get()` — runs on background thread, no widget access
4. `cute.ui do...end` — dispatches back to main thread for UI updates

## What's Safe Without cute.ui

Code that runs inside these contexts is already on the main thread:

- The `cute.app() do...end` block
- Button, combo, checkbox callbacks
- `cute.shortcut()` callbacks
- `cute.after()` and `cute.timer()` callbacks
- `cute.state().on()` observers (when triggered from the main thread)

You only need `cute.ui` when updating widgets from inside a `spawn` block.

## Example: Hacker News Reader

The HN example uses `cute.fetch()` with reactive state — stories are fetched concurrently in the background, then a reactive `cute.list()` re-renders automatically when the state updates:

```ruby
stories = cute.state([])
status = cute.state("Starting...")

load = fn(feed)
  status.set("Loading #{feed} stories...")
  stories.set([])

  cute.fetch(fn() fetch_stories(feed) end, fn(result)
    stories.set(result)
    status.set("#{len(result)} #{feed} stories")
  end)
end

# The list re-renders automatically when stories changes
cute.list(stories, fn(story, i)
  story_line(story, i)
end, fn(row) handle_selection(row) end)
```

The UI shows "Loading..." immediately, stays responsive while stories are fetched in parallel, and the list re-renders when results arrive — no manual `clear()` / `add_item()` loop needed.

---
Next: [API Reference](07-api-reference.md)
