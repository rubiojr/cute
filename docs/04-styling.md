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
import "os" as go_os

cute.app("My App", 600, 400) do
  cute.stylesheet(go_os.read_file("style.css"))
  # ...
end
```

### Targeting Widgets by ID

Use `cute.props()` to set an object name, then target it in CSS with `#id`:

```ruby
cute.stylesheet("#header { background: #ff6600; }")

cute.app("App", 600, 400) do
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
| `visible` | `set_visible` |
| `enabled` | `set_enabled` |
| `tooltip` | `set_tool_tip` |
| `css` | `set_style_sheet` — per-widget CSS |
| `wrap` | `set_word_wrap` (QLabel) |
| `align` | `set_alignment` — `"left"`, `"center"`, or `"right"` |
| `margins` | `set_contents_margins` — `[l, t, r, b]` or `[h, v]` |

`props()` returns the widget, so you can chain or use inline:

```ruby
cute.props(cute.label("Title"), {id: "title", css: cute.style({bold: true})})
```

---
Next: [Events & Shortcuts](05-events-and-shortcuts.md)
