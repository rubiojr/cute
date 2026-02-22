# Custom Widgets

In the previous tutorials we used built-in widgets like labels, buttons, and
sliders. But the real power of a UI toolkit is the ability to draw your own
widgets from scratch.

In this tutorial we'll build a **PowerBar** meter — a vertical stack of colored
segments driven by a slider. Along the way you'll learn how `cute.canvas()`
uses Qt's `paintEvent` under the hood, how to combine custom drawing with
reactive state, and how to compose everything into a reusable component.

## What We're Building

A vertical bar meter with configurable colored segments. A slider at the
bottom controls how many segments light up. The bar redraws automatically
on value changes and window resizes — no manual refresh needed.

## Getting Started

Create `examples/powerbar/main.rugo` with the basic app shell:

```ruby
require "github.com/rubiojr/cute@latest" as "cute"

cute.app("PowerBar", 160, 420) do
  level = cute.state(0)

  cute.vbox({spacing: 10, margins: [10, 10, 10, 10]}) do
    cute.label("PowerBar", {css: "font-size: 14px; font-weight: bold;"})

    cute.canvas({min_height: 280, state: level}, fn(ctx, size, v)
      ctx.fill("#1e1e2e")
    end)

    cute.slider(0, 100, fn(val)
      level.set(val)
    end)
  end
end
```

Build and run:

```bash
rugo build main.rugo -o powerbar && ./powerbar
```

You'll see a dark rectangle above a slider. Drag the slider and… nothing
visible happens yet. But behind the scenes, `level` is updating and the
canvas is repainting — it just fills with the same background color every
time.

## How Canvas Drawing Works

`cute.canvas(props, draw_fn)` creates a native drawing surface. On the Qt
backend it overrides `paintEvent` on a QWidget; on GTK it uses a
`GtkDrawingArea` with a draw callback. Either way, the widget repaints
automatically on resize and state changes. Cute invokes your `draw_fn`
with three arguments:

| Argument | Description                                |
|----------|--------------------------------------------|
| `ctx`    | Drawing context with `fill`, `rect`, `line`, `text` |
| `size`   | Hash with `width` and `height` keys        |
| `value`  | Current value of the bound `state`, or nil  |

The `state` prop in the canvas hash connects a `cute.state()` container.
When the state changes, Cute schedules a repaint through the toolkit's
event loop — no timer polling, no manual invalidation.

### Drawing primitives

| Method                          | Description                         |
|---------------------------------|-------------------------------------|
| `ctx.fill(color)`               | Fill entire canvas background       |
| `ctx.rect(x, y, w, h, color)`  | Draw a filled rectangle             |
| `ctx.line(x1, y1, x2, y2, color)` | Draw a line                      |
| `ctx.text(x, y, value, color)` | Draw text at position               |

Colors are hex strings (`"#ff0000"`) or named colors (`"red"`).

## Showing the Value

Let's prove the canvas redraws by displaying the current slider value as
text:

```ruby
cute.canvas({min_height: 280, state: level}, fn(ctx, size, v)
  ctx.fill("#1e1e2e")
  ctx.text(12, 24, "Level: #{v}", "#cdd6f4")
end)
```

Drag the slider and you'll see the text update live. The reactive loop is:

```
slider → level.set(val) → state notifies observers → canvas.update() → paintEvent → draw_fn
```

## Drawing the Bars

Now for the real thing. We need to divide the canvas height into equal
segments and fill them from bottom to top based on the current value.

The math is straightforward:

```
d_height = canvas_height - 2 × padding
step_size = d_height / n_steps
bar_height = step_size × solid_percent
spacer = step_size × (1 - solid_percent) / 2
```

Each segment `i` (counting from 0 at the bottom) is drawn at:

```
y = padding + d_height - (i + 1) × step_size + spacer
```

We subtract from `d_height` because canvas y-coordinates start at the top.

Extract the drawing into a function:

```ruby
COLORS = [
  "#49006a", "#7a0177", "#ae017e", "#dd3497",
  "#f768a1", "#fa9fb5", "#fcc5c0", "#fde0dd", "#fff7f3"
]
PADDING = 4
BAR_SOLID = 0.8
BG_COLOR = "#1e1e2e"

def draw_bar(ctx, size, value, steps, solid, pad)
  w = size["width"]
  h = size["height"]
  ctx.fill(BG_COLOR)

  n_steps = len(steps)
  d_height = h - (pad * 2)
  d_width = w - (pad * 2)
  step_size = d_height / n_steps
  bar_h = step_size * solid
  spacer = step_size * (1 - solid) / 2

  pc = value / 100.0
  n_draw = pc * n_steps

  i = 0
  while i < n_draw
    y = pad + d_height - ((i + 1) * step_size) + spacer
    ctx.rect(pad, y, d_width, bar_h, steps[i])
    i += 1
  end
end
```

Then use it in the canvas:

```ruby
cute.canvas({min_height: 280, state: level}, fn(ctx, size, v)
  draw_bar(ctx, size, v, COLORS, BAR_SOLID, PADDING)
end)
```

Each segment gets its own color from the `COLORS` array. The number of
lit segments is proportional to the slider value (0–100).

## Customisation

Since the drawing function takes all its configuration as arguments, you
can easily change the look by adjusting the constants:

```ruby
# Fewer bars, single color
COLORS = ["#f38ba8"] * 5
PADDING = 2
BAR_SOLID = 0.9

# Gradient from blue to yellow
COLORS = ["#1e66f5", "#209fb5", "#40a02b", "#df8e1d", "#fe640b"]

# Many thin bars
COLORS = ["#cba6f7"] * 20
BAR_SOLID = 0.6
PADDING = 6
```

For runtime customisation, store the configuration in state:

```ruby
bar_cfg = cute.state({
  colors: COLORS,
  solid: 0.8,
  padding: 4
})

cute.canvas({min_height: 280, state: level}, fn(ctx, size, v)
  cfg = bar_cfg.get()
  draw_bar(ctx, size, v, cfg["colors"], cfg["solid"], cfg["padding"])
end)
```

> **Note:** The canvas redraws when its bound `state` changes. If you want
> the canvas to also redraw when `bar_cfg` changes, bind both states with
> an observer:
> ```ruby
> bar_cfg.on(fn(v) level.set(level.get()) end)
> ```

## Making It a Component

Extract the whole power bar into a reusable function:

```ruby
def powerbar(min_val, max_val, colors, on_change = nil)
  n = len(colors)
  level = cute.state(0)

  cute.canvas({min_height: 200, state: level}, fn(ctx, size, v)
    draw_bar(ctx, size, v, colors, 0.8, 4)
  end)

  cute.slider(min_val, max_val, fn(val)
    level.set(val)
    if on_change != nil
      on_change(val)
    end
  end)
end
```

Now you can drop it into any layout:

```ruby
cute.app("Meters", 500, 400) do
  cute.hbox({spacing: 20, margins: [20, 20, 20, 20]}) do
    cute.vbox do
      cute.label("Volume")
      powerbar(0, 100, ["#89b4fa"] * 8, fn(v)
        puts "Volume: #{v}"
      end)
    end

    cute.vbox do
      cute.label("Temperature")
      powerbar(0, 100, [
        "#89b4fa", "#94e2d5", "#a6e3a1",
        "#f9e2af", "#fab387", "#f38ba8"
      ])
    end
  end
end
```

Components in Cute are just functions — no classes, no inheritance. They
work because Cute's layout stack automatically parents widgets to whatever
layout is currently active.

## Reactive State vs Signals

The PyQt6 version of this tutorial uses signals and slots to wire the dial
to the bar. In Cute, `cute.state()` replaces that pattern:

| PyQt6                                  | Cute                              |
|----------------------------------------|-----------------------------------|
| `dial.valueChanged.connect(bar.update)` | `level.on(fn(v) ... end)`        |
| `bar.paintEvent(self, e)`              | `canvas({state: level}, draw_fn)` |
| `self.update()`                        | Automatic on state change         |
| `class _Bar(QWidget)`                  | `def draw_bar(ctx, size, ...)`    |

State containers decouple the data flow from the widget tree. Any number
of widgets can observe the same state, and updates propagate automatically.

## Cross-Backend Styling

Qt and GTK use different CSS dialects — `QSlider::groove:horizontal` vs
`scale trough`. Use `cute.backend()` to apply the right stylesheet:

```ruby
css = qt_css
if cute.backend() == "gtk"
  css = gtk_css
end
cute.stylesheet(css)
```

`cute.backend()` returns `"qt"` or `"gtk"` so you can branch on it anywhere.

## The Final Code

The complete `examples/powerbar/main.rugo`:

```ruby
require "github.com/rubiojr/cute@latest" as "cute"

COLORS = [
  "#49006a", "#7a0177", "#ae017e", "#dd3497",
  "#f768a1", "#fa9fb5", "#fcc5c0", "#fde0dd", "#fff7f3"
]
PADDING = 4
BAR_SOLID = 0.8
BG_COLOR = "#1e1e2e"

def draw_bar(ctx, size, value, steps, solid, pad)
  w = size["width"]
  h = size["height"]
  ctx.fill(BG_COLOR)

  n_steps = len(steps)
  d_height = h - (pad * 2)
  d_width = w - (pad * 2)
  step_size = d_height / n_steps
  bar_h = step_size * solid
  spacer = step_size * (1 - solid) / 2

  pc = value / 100.0
  n_draw = pc * n_steps

  i = 0
  while i < n_draw
    y = pad + d_height - ((i + 1) * step_size) + spacer
    ctx.rect(pad, y, d_width, bar_h, steps[i])
    i += 1
  end
end

cute.app("PowerBar", 160, 420) do
  qt_css = <<~'CSS'
    QWidget {
      background: #1e1e2e;
      color: #cdd6f4;
      font-family: sans-serif;
    }
    QSlider::groove:horizontal {
      background: #313244;
      height: 6px;
      border-radius: 3px;
    }
    QSlider::handle:horizontal {
      background: #cba6f7;
      width: 16px;
      margin: -5px 0;
      border-radius: 8px;
    }
  CSS

  gtk_css = <<~'CSS'
    window {
      background: #1e1e2e;
      color: #cdd6f4;
      font-family: sans-serif;
    }
    scale trough {
      background: #313244;
      min-height: 6px;
      border-radius: 3px;
    }
    scale slider {
      background: #cba6f7;
      min-width: 16px;
      min-height: 16px;
      border-radius: 8px;
    }
  CSS

  css = qt_css
  if cute.backend() == "gtk"
    css = gtk_css
  end
  cute.stylesheet(css)

  level = cute.state(0)

  cute.vbox({spacing: 10, margins: [10, 10, 10, 10]}) do
    cute.label("PowerBar", {css: "font-size: 14px; font-weight: bold;"})

    cute.canvas({min_height: 280, state: level}, fn(ctx, size, v)
      draw_bar(ctx, size, v, COLORS, BAR_SOLID, PADDING)
    end)

    cute.slider(0, 100, fn(val)
      level.set(val)
    end)

    cute.label(level, fn(v)
      pct = "%"
      "Level: #{v}#{pct}"
    end, {css: "color: #a6adc8;"})
  end
end
```

Build and run:

```bash
cd examples/powerbar
rugo build main.rugo -o powerbar && ./powerbar
```

---
Previous: [Layout Manipulation](08-layout-manipulation.md)
