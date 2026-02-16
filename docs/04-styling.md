# Styling

Cute supports Qt stylesheets (CSS subset) for theming, plus helper functions for building styles programmatically.

## Global Stylesheets

Apply a stylesheet to the entire window with `cute.stylesheet()`:

```ruby
cute.app("Styled", 400, 300) do
  cute.stylesheet("QLabel { color: #333; font-size: 14px; }")

  cute.vbox do
    cute.label("I'm styled!")
  end
end
```

### Reactive Stylesheets

`cute.stylesheet()` also accepts a state or computed state. The stylesheet re-applies automatically whenever the state changes:

```ruby
font_size = cute.state(14)

css = cute.computed(font_size, fn(size)
  result = "QWidget { font-size: #{size}px; }"
  result
end)
cute.stylesheet(css)

# When font_size changes, the stylesheet updates automatically
cute.button("A+") do
  font_size.update(fn(v) v + 1 end)
end
```

This is how the Hacker News example implements live font-size adjustment.

### Loading from a File

Keep styles in a separate `.css` file for readability:

```css
/* style.css */
QWidget {
  font-family: sans-serif;
  font-size: 13px;
  background: #f6f6ef;
}

QPushButton {
  padding: 6px 12px;
  border-radius: 3px;
}
```

```ruby
use "os"

cute.app("My App", 600, 400) do
  cute.stylesheet(os.read_file("style.css"))
  # ...
end
```

### Targeting Widgets by ID

Use `cute.props()` to set an object name, then target it in CSS with `#id`:

```ruby
cute.app("App", 600, 400) do
  cute.stylesheet("#header { background: #ff6600; }")

  header = cute.hbox do
    cute.label("Title")
  end
  cute.props(header, {id: "header"})
end
```

## cute.style() — CSS Builder

Build CSS property strings from a hash instead of writing raw CSS:

```ruby
css = cute.style({bg: "#ff6600", color: "white", bold: true, size: 15})
# → "background: #ff6600; color: white; font-weight: bold; font-size: 15px; "
```

### Supported Keys

| Key | CSS Property | Value |
|-----|-------------|-------|
| `bg` | `background` | color string |
| `color` | `color` | color string |
| `size` | `font-size` | integer (px) |
| `bold` | `font-weight: bold` | `true` to enable |
| `italic` | `font-style: italic` | `true` to enable |
| `family` | `font-family` | string |
| `padding` | `padding` | int or array `[h, v]` or `[t, r, b, l]` |
| `margin` | `margin` | int or array (same as padding) |
| `border` | `border` | shorthand string (e.g., `"1px solid red"`) |
| `radius` | `border-radius` | integer (px) |

### Using style() with props()

Apply inline styles to individual widgets:

```ruby
lbl = cute.label("Important!")
cute.props(lbl, {css: cute.style({bold: true, color: "red", size: 18})})
```

## cute.props() — Widget Properties

Apply a hash of properties to any widget or layout:

```ruby
lbl = cute.label("Hello")
cute.props(lbl, {id: "greeting", width: 200, tooltip: "A greeting label"})

header = cute.hbox do
  cute.label("Title")
end
cute.props(header, {id: "header", margins: [8, 6, 8, 6]})
```

### Supported Keys

| Key | Effect |
|-----|--------|
| `id` | `set_object_name` — for CSS targeting with `#id` |
| `width` | `set_fixed_width` |
| `height` | `set_fixed_height` |
| `min_width` | `set_minimum_width` |
| `min_height` | `set_minimum_height` |
| `max_width` | `set_maximum_width` |
| `max_height` | `set_maximum_height` |
| `visible` | `set_visible` |
| `enabled` | `set_enabled` |
| `tooltip` | `set_tool_tip` |
| `css` | `set_style_sheet` — per-widget CSS |
| `wrap` | `set_word_wrap` (QLabel) |
| `align` | `set_alignment` — `"left"`, `"center"`, or `"right"` |
| `margins` | `set_contents_margins` — `[l, t, r, b]` or `[h, v]` |
| `spacing` | `set_spacing` (layouts only) |

### Inline Props on Widgets

Most widget constructors accept a props hash directly, so you don't need a separate `cute.props()` call:

```ruby
# Props on a label
cute.label("Title", {css: cute.style({bold: true, size: 18})})

# Props on a button
cute.button("Submit", {css: "color: green;"}) do
  save_data()
end

# Props on layouts
cute.hbox({spacing: 4, margins: [8, 6, 8, 6]}) do
  cute.label("Name:")
  cute.input("Enter name")
end
```

---
Next: [Events & Shortcuts](05-events-and-shortcuts.md)
