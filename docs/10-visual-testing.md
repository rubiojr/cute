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
./tests/visual/run.sh --update

# Verify against references
./tests/visual/run.sh
```

Tests run headlessly using `QT_QPA_PLATFORM=offscreen` -- no display
server needed.

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

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CUTE_TEST_UPDATE=1` | Save screenshots as new references |
| `CUTE_TEST_DIR=path` | Override the test directory |

## Directory Structure

```
tests/visual/
  run.sh                  # Test runner script
  counter/
    main.rugo             # Test script
    expected/             # Reference screenshots (committed)
      initial.png
      after_clicks.png
    actual/               # Generated during test runs (gitignored)
      initial.png
      after_clicks.png
```

## How It Works

1. The **widget registry** maps `id:` props to widget references. Both Qt
   and GTK backends populate this registry automatically in `props.apply()`.

2. **`test_screenshot()`** uses `QWidget::grab()` (Qt) to capture the
   window as a PNG -- works in offscreen mode without a display server.

3. **`test_click()`** calls `QAbstractButton::click()` (Qt) to
   programmatically trigger the button's click handler.

4. **Image comparison** uses ImageMagick's `magick compare -metric AE`
   to count differing pixels. A threshold of 100 pixels accommodates
   minor rendering differences across environments.

## Future: GTK Support

The test API is backend-agnostic. Adding GTK support requires implementing
`testing.rugo` and `registry.rugo` for the GTK backend:

- Screenshot: `gtk_widget_snapshot()` or Broadway-based capture
- Click: `gtk_widget_activate()` for buttons
- Registry: already implemented in `backend/gtk/lib/registry.rugo`
