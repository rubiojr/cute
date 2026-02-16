# Layouts & Widgets

## Layouts

Layouts arrange widgets automatically. Nest them with `do...end` blocks.

### vbox — Vertical Layout

Stacks widgets top-to-bottom:

```ruby
cute.vbox do
  cute.label("First")
  cute.label("Second")
  cute.label("Third")
end
```

### hbox — Horizontal Layout

Stacks widgets left-to-right:

```ruby
cute.hbox do
  cute.button("Left")
  cute.spacer()
  cute.button("Right")
end
```

### Inline Layout Props

Pass a props hash as the first argument to set spacing, margins, and other layout properties inline:

```ruby
cute.vbox({spacing: 8, margins: [12, 8, 12, 8]}) do
  cute.label("Nicely spaced")
  cute.button("OK")
end

cute.hbox({spacing: 4, margins: [8, 6, 8, 6]}) do
  cute.label("Name:")
  cute.input("Enter name")
end
```

The props hash is optional — `cute.vbox do...end` works without it.

### Nesting

Layouts nest naturally:

```ruby
cute.vbox do
  cute.hbox do
    cute.label("Name:")
    cute.input("Enter name")
  end
  cute.hbox do
    cute.spacer()
    cute.button("OK")
    cute.button("Cancel")
  end
end
```

### scroll — Scrollable Area

Wraps content in a scrollable viewport:

```ruby
cute.scroll do
  for i in 100
    cute.label("Item #{i}")
  end
end
```

## Widgets

All widget functions return the raw Qt handle, so you can call any Qt method on them directly.

### label

Static text or reactive state-bound label:

```ruby
# Static text
lbl = cute.label("Hello, world!")
lbl.set_text("Updated!")       # change text later

# With inline props
cute.label("Title", {css: "font-weight: bold;"})
```

State-aware — auto-updates when the state changes:

```ruby
count = cute.state(0)

# Simplest form — displays the value as a string
cute.label(count)

# With a transform — maps the value to display text
cute.label(count, fn(v) "Count: #{v} times" end)

# With transform and props
cute.label(count, fn(v) "#{v}" end, {width: 200})
```

### button

The click handler can be a `do...end` block, an explicit `fn()`, or omitted entirely:

```ruby
cute.button("Click") do
  puts "clicked"
end

# With inline props
cute.button("Submit", {css: "color: green;"}) do
  save_data()
end

# Explicit callback
cute.button("Click", fn() puts "clicked" end)

# No handler
cute.button("Disabled")
```

### input

Text input field with optional placeholder and inline props:

```ruby
field = cute.input("Type here...")
# field.text() returns the current text

# With props
cute.input("Search...", {width: 200})
```

Two-way state binding — the input and state stay in sync automatically:

```ruby
name = cute.state("")
cute.input("Enter name", {state: name})

# name.get() always reflects the current input text
# name.set("new value") updates the input text
```

### checkbox

```ruby
cute.checkbox("Enable notifications", fn(state)
  if state == 2
    puts "checked"
  else
    puts "unchecked"
  end
end)
```

### combo

Dropdown with a list of options:

```ruby
cute.combo(["Red", "Green", "Blue"], fn(text)
  puts "Selected: #{text}"
end)
```

### list_widget

Scrollable list for displaying items manually:

```ruby
list = cute.list_widget()
list.add_item("First item")
list.add_item("Second item")
list.clear()                    # remove all items
```

### list — Reactive List

A state-bound list that re-renders automatically when the data changes:

```ruby
items = cute.state(["Apple", "Banana", "Cherry"])

cute.list(items, fn(item, i) "#{i + 1}. #{item}" end, fn(row)
  puts "Selected row #{row}"
end)

# Later — list re-renders automatically:
items.set(["X", "Y", "Z"])
```

The second argument is a render function `fn(item, index)` that returns display text. The optional third argument is a selection callback `fn(row)`.

### spacer

Pushes subsequent widgets to the end of the layout:

```ruby
cute.hbox do
  cute.label("Left")
  cute.spacer()
  cute.label("Right")          # pushed to the right edge
end
```

### separator

Horizontal line:

```ruby
cute.vbox do
  cute.label("Above")
  cute.separator()
  cute.label("Below")
end
```

### container

A styled wrapper widget with a vertical layout inside. Children created in the `do...end` block are stacked top-to-bottom. Layout props (`spacing`, `margins`) go to the inner layout; all other props (`css`, `width`, etc.) go to the widget itself.

```ruby
cute.container({css: "background: #333; border-radius: 8px;", spacing: 6, margins: [10, 10, 10, 8]}) do
  cute.label("Card Title", {css: "font-weight: bold;"})
  cute.label("Some description text")
  cute.button("Action") do
    puts "clicked"
  end
end
```

### image / load_pixmap

`image` creates a QLabel sized for displaying images. Use the `placeholder` prop to show text while loading. `load_pixmap` creates a QPixmap from raw bytes (PNG, JPEG, etc.).

```ruby
img = cute.image({width: 200, height: 150, placeholder: "Loading..."})

# Later, e.g. inside a fetch callback:
cute.fetch(fn() http.get(url) end, fn(resp)
  pm = cute.load_pixmap(resp.body)
  img.set_pixmap(pm)
end)
```

### text_area

Multi-line text editor:

```ruby
editor = cute.text_area("Initial text", {min_height: 200})

# Read the content later:
content = editor.to_plain_text()
```

### progress

Progress bar. Set value with `.set_value(n)`, range with `.set_range(min, max)`:

```ruby
bar = cute.progress({min_width: 200})
bar.set_range(0, 100)
bar.set_value(42)
```

### slider

Horizontal slider with min/max bounds and an optional change callback:

```ruby
cute.slider(0, 100, fn(val)
  puts "Value: #{val}"
end)

# With inline props:
cute.slider(0, 100, {width: 200}, fn(val)
  volume.set(val)
end)
```

### group

Titled group box. Children in the `do...end` block are laid out vertically inside the border:

```ruby
cute.group("Settings") do
  cute.checkbox("Dark mode") do |s|
    toggle_theme(s)
  end
  cute.slider(8, 24, fn(v)
    set_font_size(v)
  end)
end
```

### tabs / tab

Tabbed interface. Use `cute.tab` inside a `cute.tabs` block to create pages:

```ruby
cute.tabs do
  cute.tab("General") do
    cute.label("General settings here")
    cute.checkbox("Enable feature") do |s| ... end
  end
  cute.tab("Advanced") do
    cute.label("Advanced settings here")
    cute.slider(0, 100, fn(v) ... end)
  end
end
```

`cute.tab` is only valid inside a `cute.tabs` block.

## Components

Extract reusable UI subtrees as regular Rugo functions:

```ruby
def toolbar(on_save, on_quit)
  cute.hbox do
    cute.button("Save") do
      on_save()
    end
    cute.spacer()
    cute.button("Quit") do
      on_quit()
    end
  end
end

cute.app("My App", 600, 400) do
  cute.vbox do
    toolbar(
      fn() puts "saving..." end,
      fn() cute.quit() end
    )
    cute.label("Content goes here")
  end
end
```

Components are just `def` functions that call `cute.*` widget builders. They work because cute uses a module-level layout stack — widgets are always added to whatever layout is currently active.

---
Next: [Reactive State](03-reactive-state.md)
