# Layouts & Widgets

## Layouts

Layouts arrange widgets automatically. Nest them with `do...end` blocks.

### vbox — Vertical Layout

Stacks widgets top-to-bottom:

```ruby
cute.vbox do
  cute.label("First")
  cute.label("Second")
  cute.label("Third")
end
```

### hbox — Horizontal Layout

Stacks widgets left-to-right:

```ruby
cute.hbox do
  cute.button("Left", nil)
  cute.spacer()
  cute.button("Right", nil)
end
```

### Nesting

Layouts nest naturally:

```ruby
cute.vbox do
  cute.hbox do
    cute.label("Name:")
    cute.input("Enter name")
  end
  cute.hbox do
    cute.spacer()
    cute.button("OK", nil)
    cute.button("Cancel", nil)
  end
end
```

### scroll — Scrollable Area

Wraps content in a scrollable viewport:

```ruby
cute.scroll do
  for i in 100
    cute.label("Item #{i}")
  end
end
```

## Widgets

All widget functions return the raw Qt handle, so you can call any Qt method on them directly.

### label

```ruby
lbl = cute.label("Hello, world!")
lbl.set_text("Updated!")       # change text later
```

### button

The second argument is the click callback. Use `nil` for no handler, or use `do...end` to pass the handler as a trailing block:

```ruby
# With inline callback
cute.button("Click", fn() puts "clicked" end)

# With do...end block
cute.button("Click") do
  puts "clicked"
end

# No handler
cute.button("Disabled", nil)
```

### input

Text input field with optional placeholder:

```ruby
field = cute.input("Type here...")
# field.text() returns the current text
```

### checkbox

```ruby
cute.checkbox("Enable notifications", fn(state)
  if state == 2
    puts "checked"
  else
    puts "unchecked"
  end
end)
```

### combo

Dropdown with a list of options:

```ruby
cute.combo(["Red", "Green", "Blue"], fn(text)
  puts "Selected: #{text}"
end)
```

### list_widget

Scrollable list for displaying items:

```ruby
list = cute.list_widget()
list.add_item("First item")
list.add_item("Second item")
list.clear()                    # remove all items
```

### spacer

Pushes subsequent widgets to the end of the layout:

```ruby
cute.hbox do
  cute.label("Left")
  cute.spacer()
  cute.label("Right")          # pushed to the right edge
end
```

### separator

Horizontal line:

```ruby
cute.vbox do
  cute.label("Above")
  cute.separator()
  cute.label("Below")
end
```

## Components

Extract reusable UI subtrees as regular Rugo functions:

```ruby
def toolbar(on_save, on_quit)
  cute.hbox do
    cute.button("Save") do
      on_save()
    end
    cute.spacer()
    cute.button("Quit") do
      on_quit()
    end
  end
end

cute.app("My App", 600, 400) do
  cute.vbox do
    toolbar(
      fn() puts "saving..." end,
      fn() cute.quit() end
    )
    cute.label("Content goes here")
  end
end
```

Components are just `def` functions that call `cute.*` widget builders. They work because cute uses a module-level layout stack — widgets are always added to whatever layout is currently active.

---
Next: [Reactive State](03-reactive-state.md)
