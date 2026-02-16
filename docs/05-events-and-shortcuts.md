# Events & Shortcuts

## Button Callbacks

Buttons take a click handler as their second argument. Use `nil` for no handler, or `do...end` for a trailing block:

```ruby
cute.button("Save", fn() save_data() end)

cute.button("Save") do
  save_data()
end
```

## Widget Signals

Widget functions return raw Qt handles. You can connect to any Qt signal directly:

```ruby
list = cute.list_widget()
list.on_current_row_changed(fn(row)
  puts "Selected row: #{row}"
end)

field = cute.input("Search...")
field.on_text_changed(fn(text)
  puts "Typed: #{text}"
end)
```

Signals use Rugo's `snake_case` naming. Common ones:

| Widget | Signal | Callback |
|--------|--------|----------|
| QPushButton | `on_clicked` | `fn()` |
| QLineEdit | `on_text_changed` | `fn(text)` |
| QComboBox | `on_current_text_changed` | `fn(text)` |
| QListWidget | `on_current_row_changed` | `fn(row)` |
| QListWidget | `on_item_double_clicked` | `fn(item)` |
| QCheckBox | `on_state_changed` | `fn(state)` |

Use `rugo doc github.com/mappu/miqt/qt6` to discover all available signals for any widget type.

## Keyboard Shortcuts

Bind keyboard shortcuts with `cute.shortcut()`:

```ruby
cute.shortcut("Ctrl+Q", fn() cute.quit() end)
cute.shortcut("Ctrl+S", fn() save_data() end)
cute.shortcut("F5", fn() refresh() end)
```

Key strings use Qt's `QKeySequence` format: `Ctrl+Key`, `Shift+Key`, `Alt+Key`, `Return`, `Escape`, `F1`–`F12`, etc.

## Timers

### One-shot Timer

Run a callback once after a delay:

```ruby
cute.after(1000, fn()
  puts "one second later"
end)
```

Commonly used to defer work until after the window appears:

```ruby
cute.app("App", 600, 400) do
  # ... build UI ...

  cute.after(50, fn()
    load_initial_data()
  end)
end
```

### Repeating Timer

Run a callback at regular intervals:

```ruby
t = cute.timer(1000, fn()
  puts "tick"
end)

# Stop it later:
t.stop()
```

## Double-Click Detection

Qt doesn't expose mouse event overrides (like `mouseDoubleClickEvent`) through the Rugo bindings. Cute provides a helper that creates a transparent overlay button and uses click timing to detect double-clicks on any widget:

```ruby
cute.on_double_click(card, fn()
  show_detail(card)
end)
```

The overlay shows a pointing-hand cursor on hover and a subtle highlight. Two clicks within 400ms trigger the callback. It returns the overlay button handle.

The overlay auto-resizes with the widget, so it works correctly with responsive layouts.

## Resize Detection

Qt's virtual method overrides (like `OnResizeEvent`) are not available in the Rugo bindings -- only signal connections are exposed. Cute provides a polling-based helper that fires a callback when a widget's size changes:

```ruby
cute.on_resize(cute.window(), fn(w, h)
  puts "Window is now #{w}x#{h}"
end)
```

The default poll interval is 200 ms. You can customize it:

```ruby
# Check every 500 ms instead
cute.on_resize(cute.window(), fn(w, h)
  reflow(w)
end, 500)
```

`on_resize` returns the underlying `QTimer` handle, so you can stop polling when no longer needed:

```ruby
t = cute.on_resize(widget, fn(w, h) relayout(w, h) end)
# later:
t.stop()
```

This is useful for responsive layouts that need to reflow content based on available width -- see [Layout Manipulation](08-layout-manipulation.md) for the `flow` helper that builds on this.

## Dialogs

### Alert

```ruby
cute.alert("Info", "Operation completed successfully.")
```

### Confirm

Returns `true` if the user clicks Yes:

```ruby
if cute.confirm("Delete", "Are you sure?")
  delete_item()
end
```

---
Next: [Threading](06-threading.md)
