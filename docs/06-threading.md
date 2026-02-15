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

The typical pattern for loading data without freezing the UI:

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

The HN example fetches 30 stories concurrently, then updates the list on the main thread:

```ruby
load = fn(feed)
  status.set("Loading...")
  story_list.clear()

  spawn
    result = fetch_stories(feed)
    cute.ui do
      stories = result
      i = 0
      for story in stories
        story_list.add_item(story_line(story, i))
        i += 1
      end
      status.set("#{len(stories)} stories")
    end
  end
end
```

The UI shows "Loading..." immediately, stays responsive while stories are fetched in parallel via `spawn`, and updates the list when results arrive.

---
Next: [API Reference](07-api-reference.md)
