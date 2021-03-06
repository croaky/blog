# Config Block with Ruby

A common interface for configuring a third-party library in Ruby:

```embed
code/config-block.rb block
```

It can be implemented like this:

```embed
code/config-block.rb module
```

The `config` class method
stores a global `Config` object
in the `Service` module.

Each config setting can be accessed like this:

```embed
code/config-block.rb access
```
