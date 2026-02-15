# Reactive State

Cute provides `cute.state()` for managing values that change over time and automatically notifying the UI.

## Creating State

```ruby
count = cute.state(0)           # initial value: 0
name = cute.state("world")      # initial value: "world"
items = cute.state([])           # initial value: empty array
```

`state()` returns a hash with three methods: `.get()`, `.set()`, and `.on()`.

## Reading and Writing

```ruby
count = cute.state(0)

puts count.get()                 # 0
count.set(42)
puts count.get()                 # 42
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

The typical pattern: create a widget, then bind state to it with `.on()`:

```ruby
cute.app("Counter", 400, 300) do
  count = cute.state(0)

  cute.vbox do
    lbl = cute.label("Count: 0")
    count.on(fn(v) lbl.set_text("Count: #{v}") end)

    cute.button("Increment") do
      count.set(count.get() + 1)
    end
  end
end
```

When the button is clicked, `count.set()` updates the value and the `.on()` callback updates the label text automatically.

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

  status_lbl = cute.label(status.get())
  status.on(fn(text) status_lbl.set_text(text) end)
end

# Later, from any callback:
status.set("Loading...")
status.set("Done — 30 items loaded")
```

The status label updates automatically each time `status.set()` is called.

---
Next: [Styling](04-styling.md)
