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

### cute.stylesheet(css)

Applies a Qt stylesheet (CSS string) to the entire window.

```ruby
cute.stylesheet("QLabel { color: red; }")
```

---

## Layouts

### cute.vbox(block)

Vertical box layout. Widgets stack top-to-bottom. Returns the `QVBoxLayout` handle.

```ruby
cute.vbox do
  cute.label("Top")
  cute.label("Bottom")
end
```

### cute.hbox(block)

Horizontal box layout. Widgets stack left-to-right. Returns the `QHBoxLayout` handle.

### cute.scroll(block)

Scrollable area. Content inside the block is placed in a scrollable viewport. Returns `QScrollArea` handle.

### cute.spacer()

Inserts a stretch spacer that pushes subsequent widgets to the end of the layout.

### cute.separator()

Inserts a horizontal separator line. Returns `QFrame` handle.

---

## Widgets

All widget functions automatically add to the current layout and return the raw Qt handle.

### cute.label(text)

Text label. Returns `QLabel`.

### cute.button(text, on_click)

Push button. `on_click` is a `fn()` callback or `nil`. Use `do...end` to pass the callback as a trailing block. Returns `QPushButton`.

```ruby
cute.button("Click") do
  puts "clicked"
end
```

### cute.input(placeholder)

Text input field. `placeholder` is a hint string or `nil`. Returns `QLineEdit`.

### cute.checkbox(text, on_change)

Checkbox. `on_change` receives the check state (`0` = off, `2` = on). Returns `QCheckBox`.

### cute.combo(items, on_change)

Dropdown. `items` is a string array. `on_change` receives the selected text. Returns `QComboBox`.

### cute.list_widget()

Scrollable list. Use `.add_item(text)` to add items and `.clear()` to remove all. Returns `QListWidget`.

---

## Properties

### cute.props(widget, hash)

Applies a property hash to any widget or layout. Returns the widget.

| Key | Type | Effect |
|-----|------|--------|
| `id` | string | `set_object_name` (for CSS `#id` targeting) |
| `width` | int | `set_fixed_width` |
| `height` | int | `set_fixed_height` |
| `visible` | bool | `set_visible` |
| `enabled` | bool | `set_enabled` |
| `tooltip` | string | `set_tool_tip` |
| `css` | string | `set_style_sheet` (per-widget CSS) |
| `wrap` | bool | `set_word_wrap` |
| `align` | string | `set_alignment` — `"left"`, `"center"`, `"right"` |
| `margins` | array | `set_contents_margins` — `[l, t, r, b]` or `[h, v]` |

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
- `.on(fn(value))` — registers a change callback

```ruby
count = cute.state(0)
count.on(fn(v) lbl.set_text("#{v}") end)
count.set(count.get() + 1)
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
