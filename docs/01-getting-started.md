# Getting Started

## Prerequisites

- **Qt6 development libraries** — `apt install qt6-base-dev` on Debian/Ubuntu
- **Go 1.22+**
- **Rugo** — `go install github.com/rubiojr/rugo@latest`

The miqt Qt6 bindings are fetched automatically on first build.

## Hello World

Create `main.rugo`:

```ruby
require "github.com/rubiojr/cute@v0.3.0"

cute.app("Hello", 400, 200) do
  cute.vbox do
    cute.label("Hello, world!")
    cute.button("Quit") do
      cute.quit()
    end
  end
end
```

Build and run:

```bash
rugo build main.rugo -o hello && ./hello
```

## How It Works

Cute is a Rugo module — you load it with `require "github.com/rubiojr/cute@v0.3.0"` and call its functions through the `cute` namespace.

The `cute.app()` function creates a window and runs the Qt event loop. Everything inside the `do...end` block builds the UI tree. Widgets are automatically added to the current layout — no manual parenting needed.

`do...end` is Rugo's trailing block syntax. It's sugar for passing a `fn()` as the last argument:

```ruby
# These are equivalent:
cute.vbox(fn()
  cute.label("Hello")
end)

cute.vbox do
  cute.label("Hello")
end
```

## Project Structure

A typical Cute app looks like:

```
my-app/
  main.rugo      # App entry point
  style.css      # Qt stylesheet (optional)
```

For larger apps, extract components into separate `.rugo` files and load them with `require`.

---
Next: [Layouts & Widgets](02-layouts-and-widgets.md)
