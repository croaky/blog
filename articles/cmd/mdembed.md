# cmd / mdembed

`mdembed` is a command line tool to embed file contents in Markdown.
It is [open source](https://github.com/croaky/mdembed).

Install:

```bash
go install github.com/croaky/mdembed
```

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
block wrapped by the magic comments. It is aware of code comment styles for
Bash, CSS, Go, Haml, HTML, JavaScript, Lua, Ruby, Sass, SQL, and TypeScript.

If you reference another Markdown file, `mdembed` will embed its contents
directly, recursively embedding its code blocks.

## So what?

I wanted the following workflow in Vim:

1. Open `tmp.md` in my project.
2. Write a prompt for an LLM (Large Language Model).
3. Reference other files, or subsets of files, in my project.
4. Hit a key combo (e.g. `Space+r`) to send the contents to an LLM.
5. Open a visual split to render the LLM's response.

`mdembed` handles step 3.
I use [mods](https://github.com/charmbracelet/mods) for the LLM steps:

```bash
go install github.com/charmbracelet/mods@latest
```

So, my Vim config runs the following Unix pipeline in a visual split:

```bash
cat example.md | mdembed | mods
```
