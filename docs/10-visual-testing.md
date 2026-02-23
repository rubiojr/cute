# Visual Testing

Cute includes a visual testing framework that lets you programmatically
interact with widgets and verify UI behavior through screenshots.

## Quick Start

Create a test script that builds the UI, interacts with it, and takes
screenshots:

```ruby
require "cute"
require "cute/lib/test" as "t"

use "os"

t.init({
  screenshot: fn(path) cute.test_screenshot(path) end,
  click: fn(id) cute.test_click(id) end,
  set_text: fn(id, text) cute.test_set_text(id, text) end,
  find: fn(id) cute.test_find(id) end,
  quit: fn() cute.quit() end
})

cute.app("My Test", 400, 300) do
  cute.button("Click Me", {id: "btn"}) do
    # ...
  end

  cute.after(500, fn()
    t.screenshot("initial")
    t.click("btn")
    t.screenshot("after_click")
    t.done()
  end)
end
```

## Running Tests

```bash
# Generate reference images (first time)
rugo run tests/visual/run.rugo update

# Verify against references (all backends)
rugo run tests/visual/run.rugo

# Run only Qt tests
rugo run tests/visual/run.rugo qt

# Run only GTK tests
rugo run tests/visual/run.rugo gtk

# Run a single test (both backends)
rugo run tests/visual/run.rugo counter

# Run a single backend-specific test
rugo run tests/visual/run.rugo qt counter
```

Qt tests run headlessly using `QT_QPA_PLATFORM=offscreen` -- no display
server needed. GTK tests use `xvfb-run` for headless execution.

## Test API

### `t.init(fns)`

Initialize the test module with cute function references. Must be called
before any other test function.

### `t.screenshot(name)`

Capture the main window and save as `actual/<name>.png`. In normal mode,
compares against `expected/<name>.png` using ImageMagick with a 100-pixel
tolerance. In update mode (`CUTE_TEST_UPDATE=1`), saves as the new
reference.

### `t.click(id)`

Simulate a click on the widget with the given `id:` prop. The widget must
have been created with an `id:` in its props hash.

### `t.set_text(id, text)`

Set the text content of an input widget found by `id:`.

### `t.find(id)`

Return the raw widget handle for the given `id:`.

### `t.wait(ms)`

Pause for the given number of milliseconds.

### `t.done()`

Print pass/fail summary and exit. Call at the end of every test.

## Writing Tests

### Widget IDs

Any widget can be targeted by the test module by passing an `id:` prop:

```ruby
cute.button("Submit", {id: "submit-btn"}) do ... end
cute.input("Name", {id: "name-input"})
cute.label(state, fn(v) v end, {id: "status-label"})
```

### Timing

Use `cute.after(ms, fn() ... end)` to delay test execution until the UI
has rendered. The test steps run inside this callback:

```ruby
cute.after(500, fn()
  t.screenshot("initial")
  t.click("my-button")
  t.screenshot("after_click")
  t.done()
end)
```

## Backend-specific Tests

Test scripts are identical except for the `require` line at the top:

```ruby
# Qt test
require "./../../../../" as "cute"

# GTK test
require "./../../../../backend/gtk" as "cute"
```

GTK tests use a slightly longer `cute.after()` delay (500ms vs 0ms) to
allow the virtual display to fully render the window before screenshots.

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CUTE_TEST_UPDATE=1` | Save screenshots as new references |
| `CUTE_TEST_DIR=path` | Override the test directory |

## Directory Structure

```
tests/visual/
  run.rugo                  # Test runner (handles both backends)
  qt/
    counter/
      main.rugo             # Qt test script
      expected/             # Qt reference screenshots (committed)
      actual/               # Generated during test runs (gitignored)
  gtk/
    counter/
      main.rugo             # GTK test script
      expected/             # GTK reference screenshots (committed)
      actual/               # Generated during test runs (gitignored)
```

## How It Works

1. The **widget registry** maps `id:` props to widget references. Both Qt
   and GTK backends populate this registry automatically in `props.apply()`.

2. **`test_screenshot()`** uses `QWidget::grab()` (Qt) or the
   `GtkWidgetPaintable → GskRenderer → GdkTexture` pipeline (GTK) to
   capture the window as a PNG.

3. **`test_click()`** calls `QAbstractButton::click()` (Qt) or
   `gtk_widget_activate()` (GTK) to programmatically trigger the button's
   click handler.

4. **Image comparison** uses ImageMagick's `magick compare -metric SSIM`
   to measure structural similarity. A threshold accommodates minor
   rendering differences across environments.

### Headless Backends

- **Qt**: `QT_QPA_PLATFORM=offscreen` -- built-in offscreen rendering,
  no display server needed.
- **GTK**: `xvfb-run -a` -- virtual X11 framebuffer. GTK4 has no built-in
  offscreen mode, so a virtual display is required. The `screenshotops` Go
  module uses the native renderer from the display surface to capture
  widgets to PNG.

## GTK Screenshot Pipeline

The GTK backend uses a Go bridge module (`screenshotops/`) that chains
GTK4/GDK/GSK C APIs via purego:

1. `gtk_widget_paintable_new(window)` -- creates a paintable observer
2. `gtk_snapshot_new()` -- creates a snapshot context
3. `gdk_paintable_snapshot()` -- renders widget tree into snapshot
4. `gtk_snapshot_free_to_node()` -- converts to a render node
5. `gsk_renderer_render_texture()` -- rasterizes to a GDK texture
6. `gdk_texture_save_to_png()` -- saves the texture to disk

## Adding a New Backend

To add visual testing for a new backend:

1. Create `backend/<name>/lib/testing.rugo` with `screenshot()`, `click()`,
   `set_text()`, and `find()` functions.
2. Wire `test_*` functions in the backend's `main.rugo`.
3. Add the backend directory under `tests/visual/<name>/`.
4. Update `run.rugo` to handle the new backend's headless execution.
