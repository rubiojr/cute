# AGENTS.md -- Coding Agent Guidelines for Cute

Cute is a declarative UI toolkit for **Rugo** (a Ruby-inspired language that
compiles to native binaries via Go). It wraps Qt6 widgets through the
[miqt](https://github.com/mappu/miqt) Go bindings into a clean `do...end`
block API.

## Project Structure

```
main.rugo              # The entire library (single file, ~470 lines)
lib/                   # Reserved for future library code
examples/
  counter/main.rugo    # Minimal counter demo
  hackernews/          # Full HN reader (main.rugo + theme.rugo)
docs/                  # Numbered tutorial series (01- through 07-)
```

## Build / Run / Test

### Prerequisites

- Qt6 development libraries (`qt6-base-dev`)
- Go 1.22+
- The `rugo` compiler

### Build Commands

```bash
# Build the counter example
rugo build main.rugo -o counter
```

```bash
# Build the hackernews example (from its directory)
rugo build main.rugo -o hackernews
```

There is no Makefile, no CI pipeline, and no linter configuration.

### Testing

There are **no automated tests** in this project. No `_test.rugo` files exist.
The examples under `examples/` serve as manual integration tests.

If tests are added in the future, Rugo uses the RATS (Rugo Automated Testing
System) framework with `_test.rugo` files.

## Code Style Guidelines

### File Layout

1. File-level header comment block describing purpose, usage, and prerequisites.
2. `require` statements (external/versioned modules first, then local modules).
3. `use` statements (standard library modules).
4. Module-level constants and state.
5. Private helpers (prefixed with `_`).
6. Public API functions.

Separate sections with horizontal-rule comment dividers:

```ruby
# -- Private helpers ----------------------------------------------------------
# -- Public API ---------------------------------------------------------------
```

### Indentation and Formatting

- **2-space indentation** everywhere.
- Lines generally under 80 characters.
- No trailing commas in hash literals or argument lists.
- One blank line between function definitions.
- No blank lines inside short function bodies; use blank lines in longer
  functions to separate logical sections.

### Imports

```ruby
# External/versioned requires first
require "github.com/mappu/miqt/qt6@v0.13.0"
require "github.com/mappu/miqt/qt6/mainthread@v0.13.0" as "_mt"

# Local module requires
require "theme"

# Standard library uses last
use "http"
use "json"
use "str"
```

- `require` for external and local modules comes before `use`.
- One import per line; no blank lines between import groups.
- Alias with `as` when the default name would conflict or is internal.

### Naming Conventions

| Element              | Convention            | Examples                              |
|----------------------|-----------------------|---------------------------------------|
| Public functions     | `snake_case`          | `app`, `vbox`, `label`, `state`       |
| Private functions    | `_snake_case`         | `_current`, `_add`, `_push`, `_pop`   |
| Local variables      | `snake_case`          | `story_list`, `detail_title`          |
| Short-scope locals   | 1-3 letter abbrev     | `s`, `v`, `cb`, `lbl`, `btn`, `inp`   |
| Constants            | `UPPER_SNAKE_CASE`    | `PAGE_SIZE`, `FEED_URLS`, `THEME`     |
| Private constants    | `_UPPER_SNAKE_CASE`   | `_CTX`                                |
| Files                | `snake_case.rugo`     | `main.rugo`, `theme.rugo`             |
| Lambda parameters    | short descriptive     | `fn(v)`, `fn(text)`, `fn(row)`        |

Entry point files are always named `main.rugo`.

### Types and Data Structures

Rugo is dynamically typed. Use hashes as ad-hoc objects:

```ruby
s = {__val__: initial, __observers__: []}
s["get"] = fn() s["__val__"] end
s["set"] = fn(v) ... end
```

- Hash keys use Ruby-style symbols: `{stack: [], win: nil}`.
- Internal/private hash fields use double-underscore wrapping: `__val__`,
  `__observers__`.
- "Methods" are closures stored as hash values.

### Error Handling

This codebase uses **nil-guard** patterns instead of explicit error propagation:

```ruby
if props == nil
  return nil
end

if on_click != nil
  btn.on_clicked(on_click)
end
```

- Check for `nil` before using a value.
- Return `nil` early when input is absent.
- No custom error types or error wrapping.

### Comments and Documentation

- All comments use `#` (no block comment syntax).
- File headers: multi-line `#` block with purpose, usage example, prerequisites.
- Function docs: placed immediately above the function. Include parameter/key
  descriptions with aligned columns for complex functions:

```ruby
# Applies a props hash to a Qt widget.
#
# Supported keys:
#   id:       string  -- set_object_name
#   width:    int     -- set_fixed_width
#   height:   int     -- set_fixed_height
def _apply_props(widget, props)
```

- Simple functions get a single-line comment:

```ruby
# Quit the application.
def quit()
```

- Use inline comments sparingly, only for non-obvious logic.
- Em dash (`--`) in prose, not `---` or `---`.

### String Handling

- Use string interpolation: `"Count: #{v} times"`.
- Use heredoc syntax for multi-line strings (e.g., CSS):

```ruby
css = <<~'CSS'
  QWidget { background: #1a1a2e; }
CSS
```

### Threading

- Long-running operations (network, I/O) go in `spawn` blocks.
- UI updates from background threads must use `cute.ui()`:

```ruby
spawn
  data = http.get(url)
  cute.ui(fn()
    label.set_text(data)
  end)
end
```

### Module Design

- The library is a single-file module (`main.rugo`).
- Multi-file apps split concerns via `require "module_name"` for local modules.
- Modules export functions, not raw constants -- provide accessor functions
  for encapsulation (e.g., `theme.css()` rather than exposing `THEME` directly).
- The implicit layout stack pattern (`_CTX`) means child widgets created inside
  `do...end` blocks are auto-parented; avoid passing explicit parent arguments.
