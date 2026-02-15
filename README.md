<p align="center">
  <img src="images/logo.svg" alt="Cute — Light up your desktop" width="600"/>
</p>

# Cute Documentation

A declarative UI toolkit for [Rugo](https://github.com/rubiojr/rugo).

Build native desktop apps with `do...end` blocks, reactive state, and zero ceremony.

```ruby
require "github.com/rubiojr/cute@v0.1.1"

cute.app("Counter", 400, 300) do
  count = cute.state(0)

  cute.vbox do
    lbl = cute.label("Clicked: 0 times")
    count.on(fn(v) lbl.set_text("Clicked: #{v} times") end)

    cute.button("Click Me") do
      count.set(count.get() + 1)
    end

    cute.hbox do
      cute.button("Reset") do
        count.set(0)
      end
      cute.button("Quit") do
        cute.quit()
      end
    end
  end

  cute.shortcut("Ctrl+Q", fn() cute.quit() end)
end
```

## Tutorial

1. [Getting Started](docs/01-getting-started.md) — Prerequisites, hello world, build & run
2. [Layouts & Widgets](docs/02-layouts-and-widgets.md) — vbox, hbox, label, button, input, and more
3. [Reactive State](docs/03-reactive-state.md) — `state()`, `.get()`, `.set()`, `.on()`
4. [Styling](docs/04-styling.md) — Stylesheets, `style()` helper, `props()`
5. [Events & Shortcuts](docs/05-events-and-shortcuts.md) — Callbacks, keyboard shortcuts, timers
6. [Threading](docs/06-threading.md) — Background work with `spawn` and `cute.ui`
7. [API Reference](docs/07-api-reference.md) — Complete function reference

## Examples

- [Counter](examples/cute_counter/) — Minimal app with reactive state
- [Hacker News Reader](examples/cute_hackernews/) — Full app with networking, lists, and threading

![](/images/shot.png)
