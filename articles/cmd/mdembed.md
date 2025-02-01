# cmd / mdembed

`mdembed` is a command line tool to embed programming file contents in Markdown.

It is [open source](https://github.com/croaky/mdembed).

## Install

```sh
go install github.com/croaky/mdembed
```

For that to work, Go must be installed and
<a href="https://go.dev/wiki/SettingGOPATH" target="_blank">`$(go env GOPATH)/bin`</a>
must be on your `$PATH`:

```sh
export PATH="$HOME/go/bin:$PATH"
```

## Use

In `input1.md`:

````md
# Examples

Embed a whole file:

```embed
f1.rb
```

Embed multiple whole files:

```embed
f1.rb
f2.go
```

Embed specific lines in a file:

```embed
f3.js log
```

Embed multiple blocks within the same file:

```embed
f4.html h1
f4.html ul
```

Embed Markdown files and their embeds recursively:

```embed
input2.md
```
````

In `f1.rb`:

```rb
puts "hi"
```

In `f2.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("hi")
}
```

In `f3.js`:

```js
console.log("Not embedded");

// emdo log
console.log("hi");
// emdone log

console.log("Not embedded");
```

In `f4.html`:

```html
<!doctype html>
<html>
  <head>
    <title>Not embedded</title>
  </head>
  <body>
    <!-- emdo h1 -->
    <h1>h1</h1>
    <!-- emdone h1 -->
    <p>not embedded</p>
    <!-- emdo ul -->
    <ul>
      <li>1</li>
      <li>2</li>
    </ul>
    <!-- emdone ul -->
  </body>
</html>
```

In `input2.md`:

````md
## Input2

Embed from within an embedded Markdown file:

```embed
f1.rb
```
````

Run:

```bash
cat input.md | mdembed
```

The output will be:

````md
# Examples

Embed a whole file:

```rb
# f1.rb
puts "hi"
```

Embed multiple whole files:

```rb
# f1.rb
puts "hi"
```

```go
// f2.go
package main

import "fmt"

func main() {
    fmt.Println("hi")
}
```

Embed specific lines in a file:

```js
// f3.js
console.log("hi");
```

Embed multiple blocks within the same file:

```html
<!-- f4.html -->
<h1>h1</h1>
```

```html
<!-- f4.html -->
<ul>
  <li>1</li>
  <li>2</li>
</ul>
```

Embed Markdown files and their embeds recursively:

## Input2

Embed from within an embedded Markdown file:

```rb
# f1.rb
puts "hi"
```
````

`mdembed` embeds code blocks in the output Markdown,
removing surrounding whitespace.
The file extension is used as the code fence attribute.

If `emdo` and `emdone` magic comments were used, it will only embed the code
block wrapped by the magic comments.

It is aware of code comment styles for Ada, Assembly, Awk, Bash, C, Clojure,
COBOL, C++, C#, CSS, D, Dart, Elm, Erlang, Elixir, Fortran, F#, Gleam, Go, Haml,
Haskell, HTML, Java, Julia, JavaScript, Kotlin, Lisp, Logo, Lua, MATLAB, OCaml,
Objective-C, Mojo, Nim, Pascal, PHP, Perl, Prolog, Python, R, Ruby, Rust, Scala,
Schema, Sass, Shell, Solidity, SQL, Swift, Tcl, TypeScript, VBScript, Visual
Basic, Wolfram, and Zig.

If you reference another Markdown file, `mdembed` will embed its contents
directly, recursively embedding its code blocks.

## So what?

I wanted the following workflow in Vim:

1. Open `tmp.md` in my project.
2. Write a prompt for an LLM (Large Language Model),
   referencing other files, or subsets of files, in my project.
3. Hit a key combo (`<Leader>+r`) to send the contents to an LLM
   and render the LLM's response in a vertical split.

`mdembed` handles the Markdown parsing steps.
[mods](https://github.com/charmbracelet/mods) handles the LLM steps:

```bash
go install github.com/charmbracelet/mods@latest
```

So, my Vim config runs the following Unix pipeline in a vertical split:

```bash
cat example.md | mdembed | mods
```
