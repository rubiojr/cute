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

### cute.container(props = nil, block)

A `QWidget` wrapper that holds an inner vertical layout. Props are split: `spacing` and `margins` go to the inner layout, everything else (`css`, `width`, `height`, etc.) goes to the widget. Returns the `QWidget` handle.

```ruby
cute.container({css: "background: #1e1e2e;", spacing: 8}) do
  cute.label("Inside a container")
  cute.button("Click me")
end
```

### cute.new_hbox(props = nil)

Creates a free-standing `QHBoxLayout` that is **not** added to the current layout context. Use with `layout.add_layout()` or as a row in `cute.flow`.

```ruby
row = cute.new_hbox({spacing: 12})
cute.add_widget(row, card1)
cute.add_widget(row, card2)
```

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

## More Widgets

### cute.image(props = nil)

Displays a `QLabel` configured as an image placeholder. Set the pixmap later via `label.set_pixmap(pixmap)`. Returns the `QLabel` handle.

```ruby
img = cute.image({width: 200, height: 200})
# Later, after loading:
img.set_pixmap(pixmap)
```

### cute.canvas(props = nil, draw_fn = nil)

Composition-based drawing widget backed by a `QLabel`.

`draw_fn(ctx, size, state_value)` is called on creation and resize. If `props["state"]` is provided, it is also called whenever that state changes, and `state_value` receives the current state value.

`ctx` exposes drawing methods:
- `ctx.fill(color)`
- `ctx.line(x1, y1, x2, y2, color = nil)`
- `ctx.text(x, y, text, color = nil)`

`size` is a hash with `width` and `height`.

```ruby
phase = cute.state(0)
canvas = cute.canvas({min_height: 180, state: phase}, fn(ctx, size, v)
  w = size["width"]
  ctx.fill("#11111b")
  ctx.line(10, 20, w - 10, 20, "#89b4fa")
  ctx.text(12, 40, "phase: #{v}", "#cdd6f4")
end)
```

### cute.canvas_frame(width, height, bg = nil)

Creates an internal frame object for canvas drawing.

### cute.canvas_fill(frame, color)

Fills an existing frame (advanced usage).

### cute.canvas_line(frame, x1, y1, x2, y2, color = nil)

Draws a line on an existing frame (advanced usage).

### cute.canvas_text(frame, x, y, text, color = nil)

Draws text on an existing frame (advanced usage).

### cute.canvas_commit(canvas, frame)

Assigns the frame result to the canvas widget.

### cute.load_pixmap(data)

Converts raw image bytes (e.g., from `http.get().body_bytes`) into a `QPixmap`. Returns `nil` if the data cannot be decoded.

```ruby
resp = http.get("https://picsum.photos/200/200")
pix = cute.load_pixmap(resp.body_bytes)
if pix != nil
  img.set_pixmap(pix)
end
```

### cute.text_area(text = "", props = nil)

Multi-line read-only text area. Returns `QPlainTextEdit` with editing disabled.

```ruby
cute.text_area("Long text here...", {min_height: 200})
```

### cute.progress(props = nil)

Indeterminate progress bar (range 0..0). Apply props like `min_width` or `css`. Returns `QProgressBar`.

```ruby
bar = cute.progress({min_width: 200})
# To show determinate progress, set range and value on the handle:
# bar.set_range(0, 100); bar.set_value(50)
```

### cute.slider(min, max, arg1 = nil, arg2 = nil)

Horizontal slider. `arg1` can be a props hash or a change callback; `arg2` is the callback if `arg1` is props. Returns `QSlider`.

```ruby
cute.slider(0, 100, fn(val)
  puts "Value: #{val}"
end)

cute.slider(10, 50, {width: 200}, fn(val)
  volume.set(val)
end)
```

### cute.group(title, props = nil, block)

Group box with a title and inner vertical layout. Returns the `QGroupBox` handle.

```ruby
cute.group("Settings") do
  cute.checkbox("Enable notifications")
  cute.checkbox("Dark mode")
end
```

### cute.tabs(props = nil, block)

Tabbed widget container. Use `cute.tab` inside the block to add pages. Returns `QTabWidget`.

```ruby
cute.tabs do
  cute.tab("General") do
    cute.label("General settings")
  end
  cute.tab("Advanced") do
    cute.label("Advanced settings")
  end
end
```

### cute.tab(title, props = nil, block)

A single tab page inside `cute.tabs`. Creates a `QWidget` with a vertical layout and adds it as a tab. Returns the page `QWidget`.

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

### cute.on_resize(widget, callback, interval_ms = 200)

Polls `widget` for size changes every `interval_ms` milliseconds. Calls `callback(width, height)` when the size differs from the previous check. Returns the `QTimer` handle (`.stop()` to cancel).

```ruby
cute.on_resize(cute.window(), fn(w, h)
  cols = w / 280
  reflow(cols)
end)
```

### cute.on_double_click(widget, callback)

Detects double-click on any widget. Creates a transparent overlay button that fires `callback` when two clicks occur within 400ms. Shows a pointing-hand cursor and subtle hover highlight. The overlay auto-resizes with the widget. Returns the overlay button handle.

```ruby
cute.on_double_click(card, fn()
  show_detail(card)
end)
```

---

## Dynamic Layout

### cute.add_to(layout, block)

Pushes `layout` onto the context stack and runs `block`. Any `cute.*` widget calls inside the block are auto-parented to this layout. Use this to add widgets to a layout after initial construction.

```ruby
cute.add_to(card_layout) do
  cute.label("Added later")
end
```

### cute.detached(layout, block)

Builds a widget tree via the Cute DSL, then detaches the result from the layout and returns it as a hidden `QWidget`. The widget can later be placed with `cute.add_widget` or `cute.flow`.

```ruby
card = cute.detached(grid) do
  cute.container({css: "background: #333;"}) do
    cute.label("Hello")
  end
end
```

### cute.clear_layout(layout)

Removes all items from a layout back-to-front. Child widgets are **not** destroyed -- they can be re-added. Spacer and stretch items are discarded.

```ruby
cute.clear_layout(grid)
```

### cute.add_widget(layout, widget)

Appends an existing widget to a layout and calls `.show()` on it.

```ruby
cute.add_widget(row, card)
```

### cute.add_stretch(layout)

Appends an expanding spacer to a layout.

### cute.flow(layout, items, cols, props = nil)

Arranges an array of widgets into a grid inside a vertical layout. Clears the layout first, then creates `hbox` rows of `cols` columns. Each row gets a trailing stretch for left-alignment.

| Param | Type | Description |
|-------|------|-------------|
| `layout` | QVBoxLayout | The container layout to fill |
| `items` | array | Widgets to arrange |
| `cols` | int | Number of columns per row |
| `props` | hash | Optional -- passed to each row's hbox (e.g., `{spacing: 12}`) |

```ruby
cute.flow(grid, cards, 3, {spacing: 12})
```

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
