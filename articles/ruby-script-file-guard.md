# Ruby script file guard

One of the most common snippets of Ruby code I write is:

```ruby
if $0 == __FILE__
  # pp something
end
```

I use this snippet to ensure that certain code is only executed when the file is
run directly, and not when it is required or loaded as a module by another
script.

`$0` contains the file name of the Ruby script currently being run. If the
script is being executed directly, `$0` will hold the name of the script file.
If the script is being required or loaded as a module, `$0` will hold the name
of the main script that initiated the loading.

`__FILE__` is a built-in constant in Ruby that holds the current file's path. It
represents the path to the file that contains the current line of code.

Here's a longer example:

```embed
code/ruby/lib/github/client.rb
```

I might put this file in `lib/github/client.rb` and
reference it from other files in `lib/**/*.rb`,
which will not run the code in the `if $0 == __FILE__` guard.

But, I can also run the file directly to print the return value
of the `#get` method, which offers a tight feedback loop for testing.
