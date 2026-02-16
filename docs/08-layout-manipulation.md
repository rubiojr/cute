# Layout Manipulation

The core Cute layout system (`vbox`, `hbox`, `scroll`) covers static UI trees built at startup. But many apps need to add, remove, or rearrange widgets at runtime -- loading new content, reflowing a grid on resize, or building infinite-scroll lists.

This guide covers Cute's layout manipulation helpers.

## Clearing a Layout

`cute.clear_layout(layout)` removes all items from a layout without destroying the child widgets. This lets you re-add them later in a different arrangement:

```ruby
cute.clear_layout(grid)
# grid is now empty, but widgets formerly in it are still alive
```

Items are removed back-to-front. Spacer and stretch items are discarded.

## Adding Widgets to a Layout

### cute.add_widget(layout, widget)

Appends an existing widget to a layout and shows it:

```ruby
cute.add_widget(row, card)
```

Unlike `cute.add_to`, this does not use the context stack -- it works with widgets that were created elsewhere (e.g., via `detached`).

### cute.add_stretch(layout)

Appends an expanding spacer. Useful for left-aligning items in a row:

```ruby
cute.add_widget(row, card1)
cute.add_widget(row, card2)
cute.add_stretch(row)   # pushes cards to the left
```

## Creating Detached Layouts and Widgets

### cute.new_hbox(props)

Creates a free-standing horizontal box layout that is **not** added to the current context. Use it as a row inside a grid or pass it to `layout.add_layout()`:

```ruby
row = cute.new_hbox({spacing: 12})
cute.add_widget(row, card1)
cute.add_widget(row, card2)
grid.add_layout(row)
```

### cute.detached(layout, block)

Builds a widget tree using the Cute DSL but immediately removes it from the layout, returning the top-level widget in a hidden state. The widget can be placed later with `add_widget` or `flow`:

```ruby
card = cute.detached(grid) do
  cute.container({css: "background: #333;"}) do
    cute.label("Hello")
  end
end
# card is a QWidget, hidden, not in any layout yet
```

This is the primary way to create reusable widget objects that you manage manually.

## Responsive Grid: cute.flow

`cute.flow(layout, items, cols, props)` arranges an array of widgets into a grid of horizontal rows inside a vertical layout:

```ruby
cute.flow(grid, cards, 3, {spacing: 12})
```

It clears the layout first, then creates `hbox` rows with the given column count. A trailing stretch is added to each row so items stay left-aligned when the last row is incomplete.

### Responsive Reflow

Combine `flow` with `on_resize` to make the column count adapt to the window width:

```ruby
cols = 3
grid = nil

cute.scroll do
  grid = cute.vbox({spacing: 12})
end

cute.on_resize(cute.window(), fn(w, h)
  new_cols = w / 280
  if new_cols < 1
    new_cols = 1
  end
  if new_cols != cols
    cols = new_cols
    cute.flow(grid, cards, cols, {spacing: 12})
  end
end)
```

Each time the window width changes, the grid reflows to fit as many columns as possible at 280px per card.

## Worked Example: Thumbnail Grid

A simplified responsive thumbnail grid in ~25 lines:

```ruby
require "cute"
use "http"

cute.app("Thumbnails", 800, 600) do
  cards = []
  cols = 3
  grid = nil

  cute.scroll do
    grid = cute.vbox({spacing: 12})
  end

  # Create 9 image cards
  i = 0
  while i < 9
    url = "https://picsum.photos/id/#{i + 10}/200/200"
    card = cute.detached(grid) do
      cute.container({width: 200, height: 200}) do
        cute.label("Loading...")
      end
    end
    cards = cards + [card]
    i = i + 1
  end

  # Initial flow
  cute.flow(grid, cards, cols, {spacing: 12})

  # Reflow on resize
  cute.on_resize(cute.window(), fn(w, h)
    new_cols = w / 220
    if new_cols < 1
      new_cols = 1
    end
    if new_cols != cols
      cols = new_cols
      cute.flow(grid, cards, cols, {spacing: 12})
    end
  end)
end
```

For a full example with background image loading and endless scroll, see `examples/gallery/`.

---
Next: [API Reference](07-api-reference.md)
