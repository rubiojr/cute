# API Reference

Complete reference for all public `cute.*` functions.

## Application

### cute.app(title, width, height, block)

Creates a window and runs the Qt event loop. The `block` builds the UI tree.

```ruby
cute.app("My App", 800, 600) do
  # build UI here
end
```

### cute.quit()

Exits the application.

### cute.window()

Returns the raw `QWidget` handle for the main window.

### cute.stylesheet(css_or_state)

Applies a Qt stylesheet (CSS string) to the entire window. Also accepts a state or computed state — the stylesheet re-applies automatically when it changes.

```ruby
cute.stylesheet("QLabel { color: red; }")

# Reactive — auto-updates when font_size changes
css = cute.computed(font_size, fn(size)
  result = "QWidget { font-size: #{size}px; }"
  result
end)
cute.stylesheet(css)
```

---

## Layouts

### cute.vbox(props = nil, block)

Vertical box layout. Widgets stack top-to-bottom. Returns the `QVBoxLayout` handle.

```ruby
cute.vbox do
  cute.label("Top")
  cute.label("Bottom")
end

cute.vbox({spacing: 8, margins: [12, 8, 12, 8]}) do
  cute.label("Spaced")
end
```

### cute.hbox(props = nil, block)

Horizontal box layout. Widgets stack left-to-right. Returns the `QHBoxLayout` handle.

```ruby
cute.hbox({spacing: 4}) do
  cute.button("Left")
  cute.spacer()
  cute.button("Right")
end
```

### cute.scroll(props = nil, block)

Scrollable area. Content inside the block is placed in a scrollable viewport. Returns `QScrollArea` handle.

```ruby
cute.scroll({css: "background: #fff;"}) do
  # scrollable content
end
```

### cute.spacer()

Inserts a stretch spacer that pushes subsequent widgets to the end of the layout.

### cute.separator()

Inserts a horizontal separator line. Returns `QFrame` handle.

---

## Widgets

All widget functions automatically add to the current layout and return the raw Qt handle.

### cute.label(text_or_state, arg2 = nil, arg3 = nil)

Text label. Returns `QLabel`.

Static text:

```ruby
cute.label("Hello")
cute.label("Title", {css: "font-weight: bold;"})
```

State-aware — auto-updates when the state changes:

```ruby
cute.label(count)                                   # displays value as string
cute.label(count, fn(v) "Count: #{v}" end)          # with transform
cute.label(count, fn(v) "#{v}" end, {width: 200})   # with transform and props
```

### cute.button(text, arg1 = nil, arg2 = nil)

Push button. Use `do...end` for the click handler. Optional props hash for inline styling. Returns `QPushButton`.

```ruby
cute.button("Click") do
  puts "clicked"
end

cute.button("Submit", {css: "color: green;"}) do
  save_data()
end
```

### cute.input(placeholder = "", props = nil)

Text input field. `placeholder` is a hint string. Returns `QLineEdit`.

Two-way state binding via the `state:` prop key:

```ruby
name = cute.state("")
cute.input("Enter name", {state: name})
# name.get() reflects input text; name.set() updates it
```

### cute.checkbox(text, arg1 = nil, arg2 = nil)

Checkbox. `on_change` receives the check state (`0` = off, `2` = on). Returns `QCheckBox`.

```ruby
cute.checkbox("Dark mode") do |state|
  toggle_theme(state)
end

cute.checkbox("Dark mode", {css: "color: #ccc;"}) do |state|
  toggle_theme(state)
end
```

### cute.combo(items, arg1 = nil, arg2 = nil)

Dropdown. `items` is a string array. `on_change` receives the selected text. Returns `QComboBox`.

```ruby
cute.combo(["Red", "Green", "Blue"], fn(text)
  puts "Selected: #{text}"
end)

cute.combo(["A", "B", "C"], {width: 120}, fn(text)
  puts text
end)
```

### cute.list_widget(props = nil)

Scrollable list for manual item management. Use `.add_item(text)` to add items and `.clear()` to remove all. Returns `QListWidget`.

### cute.list(items_state, render_fn, on_select = nil)

Reactive list bound to a state. Re-renders automatically when the state changes. Returns `QListWidget`.

`render_fn` receives `(item, index)` and returns display text. `on_select` receives the selected row index.

```ruby
items = cute.state(["Apple", "Banana", "Cherry"])
cute.list(items, fn(item, i) "#{i + 1}. #{item}" end, fn(row)
  puts "Selected row #{row}"
end)

# Re-renders automatically:
items.set(["X", "Y", "Z"])
```

---

## Properties

### cute.props(widget, hash)

Applies a property hash to any widget or layout. Returns the widget.

| Key | Type | Effect |
|-----|------|--------|
| `id` | string | `set_object_name` (for CSS `#id` targeting) |
| `width` | int | `set_fixed_width` |
| `height` | int | `set_fixed_height` |
| `min_width` | int | `set_minimum_width` |
| `min_height` | int | `set_minimum_height` |
| `max_width` | int | `set_maximum_width` |
| `max_height` | int | `set_maximum_height` |
| `visible` | bool | `set_visible` |
| `enabled` | bool | `set_enabled` |
| `tooltip` | string | `set_tool_tip` |
| `css` | string | `set_style_sheet` (per-widget CSS) |
| `wrap` | bool | `set_word_wrap` |
| `align` | string | `set_alignment` — `"left"`, `"center"`, `"right"` |
| `margins` | array | `set_contents_margins` — `[l, t, r, b]` or `[h, v]` |
| `spacing` | int | `set_spacing` (layouts only) |

### cute.style(hash)

Builds a CSS property string from a hash. Returns a string.

| Key | Type | CSS |
|-----|------|-----|
| `bg` | string | `background: <value>` |
| `color` | string | `color: <value>` |
| `size` | int | `font-size: <value>px` |
| `bold` | bool | `font-weight: bold` |
| `italic` | bool | `font-style: italic` |
| `family` | string | `font-family: <value>` |
| `padding` | int/array | `padding: <value>` |
| `margin` | int/array | `margin: <value>` |
| `border` | string | `border: <value>` |
| `radius` | int | `border-radius: <value>px` |

---

## State

### cute.state(initial)

Creates a reactive state container. Returns a hash with:

- `.get()` — returns the current value
- `.set(value)` — updates and notifies observers
- `.update(fn)` — transform in place: `count.update(fn(v) v + 1 end)`
- `.on(fn(value))` — registers a change callback

```ruby
count = cute.state(0)
count.update(fn(v) v + 1 end)   # 1
count.on(fn(v) puts "now #{v}" end)
count.set(10)                    # prints: now 10
```

### cute.computed(source, transform)

Creates a derived state that auto-updates when the source changes. Returns a read-only state (has `.get()` and `.on()`).

```ruby
count = cute.state(0)
display = cute.computed(count, fn(v) "Count: #{v}" end)
cute.label(display)   # auto-updates when count changes
```

Works anywhere state is accepted: labels, stylesheets, bind, etc.

### cute.bind(state, widget, prop, transform = nil)

Binds a state to a widget property. Sets the initial value and auto-updates on changes.

Supported property names: `"text"`, `"visible"`, `"enabled"`, `"css"`, `"tooltip"`.

```ruby
status = cute.state("Ready")
lbl = cute.label("Ready")
cute.bind(status, lbl, "text")

# With transform:
cute.bind(count, lbl, "text", fn(v) "Count: #{v}" end)
```

---

## Events

### cute.shortcut(key, callback)

Binds a keyboard shortcut. `key` uses Qt's key sequence format (e.g., `"Ctrl+Q"`, `"F5"`). Returns `QShortcut` handle.

### cute.after(ms, callback)

Single-shot timer. Fires `callback` once after `ms` milliseconds. Returns `QTimer` handle.

### cute.timer(ms, callback)

Repeating timer. Fires `callback` every `ms` milliseconds. Call `.stop()` on the returned handle to cancel.

---

## Dialogs

### cute.alert(title, message)

Shows an information dialog.

### cute.confirm(title, message)

Shows a Yes/No dialog. Returns `true` if Yes is clicked.

---

## Threading

### cute.ui(callback)

Runs `callback` on the Qt main thread and blocks until it completes. Use this from inside `spawn` blocks to safely update widgets.

```ruby
spawn
  data = fetch_data()
  cute.ui do
    lbl.set_text(data)
  end
end
```

### cute.fetch(work, on_done)

Runs `work` in a background thread and dispatches the result to the UI thread. Combines `spawn` + `cute.ui()` into one call.

```ruby
cute.fetch(fn() http.get(url) end, fn(resp)
  lbl.set_text(resp.body)
end)
```
