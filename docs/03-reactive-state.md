# Reactive State

Cute provides `cute.state()` for managing values that change over time and automatically notifying the UI.

## Creating State

```ruby
count = cute.state(0)           # initial value: 0
name = cute.state("world")      # initial value: "world"
items = cute.state([])           # initial value: empty array
```

`state()` returns a hash with four methods: `.get()`, `.set()`, `.update()`, and `.on()`.

## Reading and Writing

```ruby
count = cute.state(0)

puts count.get()                 # 0
count.set(42)
puts count.get()                 # 42
```

## Updating with a Transform

Use `.update(fn)` to transform the current value in place — no need to read and write separately:

```ruby
count = cute.state(0)

count.update(fn(v) v + 1 end)   # 1
count.update(fn(v) v + 1 end)   # 2
count.update(fn(v) v * 10 end)  # 20
```

This is equivalent to `count.set(fn(count.get()))` but cleaner, especially when used in callbacks:

```ruby
cute.button("+1") do
  count.update(fn(v) v + 1 end)
end
```

## Observing Changes

Register callbacks that fire whenever the value changes:

```ruby
count = cute.state(0)
count.on(fn(v) puts "count is now #{v}" end)

count.set(1)   # prints: count is now 1
count.set(2)   # prints: count is now 2
```

Multiple observers are supported — they fire in registration order.

## Binding State to Widgets

### State-aware label (recommended)

The simplest way — pass a state directly to `cute.label()`:

```ruby
cute.app("Counter", 400, 300) do
  count = cute.state(0)

  cute.vbox do
    cute.label(count, fn(v) "Count: #{v}" end)

    cute.button("Increment") do
      count.update(fn(v) v + 1 end)
    end
  end
end
```

The label auto-subscribes and updates whenever `count` changes.

### Manual binding with .on()

For full control, use `.on()` directly:

```ruby
lbl = cute.label("Count: 0")
count.on(fn(v) lbl.set_text("Count: #{v}") end)
```

### cute.bind() — bind state to any widget property

Bind a state to any widget property (`"text"`, `"visible"`, `"enabled"`, `"css"`, `"tooltip"`):

```ruby
status = cute.state("Ready")
lbl = cute.label("Ready")
cute.bind(status, lbl, "text")

# With a transform:
count = cute.state(0)
lbl = cute.label("0")
cute.bind(count, lbl, "text", fn(v) "Count: #{v}" end)
```

## Computed State

`cute.computed()` creates a derived state that auto-updates when its source changes:

```ruby
count = cute.state(0)
display = cute.computed(count, fn(v) "Clicked #{v} times" end)

cute.label(display)   # auto-updates when count changes
```

Computed states are read-only — they update automatically when the source state changes. They have `.get()` and `.on()` like regular state, so they work anywhere state is accepted (labels, stylesheets, bind, etc.).

```ruby
# Reactive theme that re-applies when font_size changes
font_size = cute.state(14)
css = cute.computed(font_size, fn(size)
  result = "QWidget { font-size: #{size}px; }"
  result
end)
cute.stylesheet(css)
```

**Note:** Rugo's `if/else/end` does not work as a return expression in lambdas. Use an explicit variable instead:

```ruby
# ✗ WRONG — if/else/end doesn't return a value
cute.computed(mode, fn(v)
  if v == "dark"
    dark_css()
  else
    light_css()
  end
end)

# ✓ CORRECT — use an explicit variable
cute.computed(mode, fn(v)
  result = light_css()
  if v == "dark"
    result = dark_css()
  end
  result
end)
```

## State vs Plain Variables

Use `cute.state()` when:

- A value drives UI updates (labels, status bars, etc.)
- You want to decouple data changes from UI updates
- Multiple widgets need to react to the same value

Use plain variables when:

- The value is only read once (e.g., building a list)
- The value is internal to a callback and doesn't affect the UI

## Example: Status Bar

```ruby
status = cute.state("Ready")

cute.vbox do
  # ... app content ...

  cute.label(status)             # auto-updates on changes
end

# Later, from any callback:
status.set("Loading...")
status.set("Done — 30 items loaded")
```

The status label updates automatically each time `status.set()` is called.

---
Next: [Styling](04-styling.md)
