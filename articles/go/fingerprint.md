# go / fingerprint

I use file-based asset fingerprinting in Go web apps
to enable aggressive caching with [CDNs](/web/cdn).

## The problem

When serving CSS, JavaScript, or other static assets,
browsers cache them to improve performance.
But when I update a file,
I need browsers to fetch the new version instead of using the cached copy.

Common solutions like cache headers with short TTLs or `?v=123` query strings
either sacrifice caching performance or require manual version management.

## Fingerprinting

Asset fingerprinting generates a unique URL for each version of a file
by including a hash of its contents in the filename or path.
When the file changes, the hash changes, creating a new URL.

This allows me to:

- Set long cache headers (1 year) for maximum performance
- Automatically invalidate caches when files change
- Avoid manual version tracking

## Environment configuration

Load environment variables early at startup to fail fast on misconfiguration:

```go
type Env struct {
    AppEnv string // "dev" or "prod"
}

func LoadEnv() Env {
    appEnv := os.Getenv("APP_ENV")
    if appEnv == "" {
        appEnv = "dev"
    }
    return Env{AppEnv: appEnv}
}

func (e Env) Dev() bool {
    return e.AppEnv == "dev"
}

func main() {
    env := LoadEnv()
    server := NewServer(env)
    fmt.Printf("Serving at http://localhost:4000 [%s]\n", env.AppEnv)
    log.Fatal(http.ListenAndServe(":4000", server.Handler()))
}
```

## Implementation

I compute file hashes at server startup
and serve assets at fingerprinted URLs.
In development, skip fingerprinting for live reloading.

```go
type Server struct {
    env        Env
    cssPath    string            // Fingerprinted CSS path
    imgPaths   map[string]string // Original -> fingerprinted
    fontPaths  map[string]string // Original -> fingerprinted
    cssContent []byte            // Processed CSS with rewritten URLs
}

func fileDigest(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := md5.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func NewServer(env Env) *Server {
    s := &Server{
        env:       env,
        imgPaths:  make(map[string]string),
        fontPaths: make(map[string]string),
    }

    // In dev mode, skip fingerprinting, load as-is
    if s.env.Dev() {
        s.cssPath = "/ui/app.css"
        return s
    }

    // Fingerprint images
    imgFiles, _ := filepath.Glob("ui/img/*")
    for _, file := range imgFiles {
        name := filepath.Base(file)
        ext := filepath.Ext(name)
        base := name[:len(name)-len(ext)]
        if hash, err := fileDigest(file); err == nil {
            s.imgPaths[name] = fmt.Sprintf("%s-%s%s", base, hash[:8], ext)
        }
    }

    // Fingerprint fonts
    fontFiles, _ := filepath.Glob("ui/font/*")
    for _, file := range fontFiles {
        name := filepath.Base(file)
        ext := filepath.Ext(name)
        base := name[:len(name)-len(ext)]
        if hash, err := fileDigest(file); err == nil {
            s.fontPaths[name] = fmt.Sprintf("%s-%s%s", base, hash[:8], ext)
        }
    }

    // Process CSS: rewrite asset URLs to fingerprinted versions
    cssBytes, err := os.ReadFile("ui/app.css")
    if err != nil {
        log.Fatal(err)
    }
    cssContent := string(cssBytes)
    for orig, fp := range s.imgPaths {
        cssContent = strings.ReplaceAll(cssContent, "img/"+orig, "img/"+fp)
    }
    for orig, fp := range s.fontPaths {
        cssContent = strings.ReplaceAll(cssContent, "font/"+orig, "font/"+fp)
    }
    s.cssContent = []byte(cssContent)

    // Fingerprint CSS from processed content
    h := md5.New()
    h.Write(s.cssContent)
    s.cssPath = fmt.Sprintf("/ui/app-%s.css", fmt.Sprintf("%x", h.Sum(nil))[:8])

    return s
}
```

CSS requires special handling because it references other assets
via `url()`. The CSS content is modified in memory
to rewrite font and image paths to their fingerprinted versions,
then served from memory.
Images and fonts are binary blobs that don't reference other files,
so they're served directly from disk.

## Serving assets

In development, serve assets from disk with `no-cache` headers
so edits are immediately visible.
In production, serve fingerprinted assets with 1-year cache headers.

```go
func (s *Server) Handler() http.Handler {
    mux := http.NewServeMux()

    noCache := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Cache-Control", "no-cache")
            next.ServeHTTP(w, r)
        })
    }

    if s.env.Dev() {
        // Dev: serve from disk for live reloading
        mux.HandleFunc("GET /ui/app.css", func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/css")
            w.Header().Set("Cache-Control", "no-cache")
            http.ServeFile(w, r, "./ui/app.css")
        })
        mux.Handle("GET /ui/img/", noCache(http.StripPrefix("/ui/img/",
            http.FileServer(http.Dir("./ui/img")))))
        mux.Handle("GET /ui/font/", noCache(http.StripPrefix("/ui/font/",
            http.FileServer(http.Dir("./ui/font")))))
    } else {
        // Production: serve fingerprinted assets with long cache
        mux.HandleFunc("GET "+s.cssPath, func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "text/css")
            w.Header().Set("Cache-Control", "public, max-age=31536000")
            w.Write(s.cssContent)
        })

        for orig, fp := range s.imgPaths {
            origFile := orig
            mux.HandleFunc("GET /ui/img/"+fp, func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Cache-Control", "public, max-age=31536000")
                http.ServeFile(w, r, "./ui/img/"+origFile)
            })
        }

        for orig, fp := range s.fontPaths {
            origFile := orig
            mux.HandleFunc("GET /ui/font/"+fp, func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Cache-Control", "public, max-age=31536000")
                http.ServeFile(w, r, "./ui/font/"+origFile)
            })
        }
    }

    mux.HandleFunc("GET /", s.index)
    return mux
}
```

## Templates

Pass the CSS path to templates for rendering:

```go
type PageData struct {
    Title   string
    CSSPath string
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("ui/index.html"))
    tmpl.Execute(w, PageData{
        Title:   "Home",
        CSSPath: s.cssPath,
    })
}
```

In the template:

```html
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="{{.CSSPath}}" />
  </head>
  <body>
    <h1>{{.Title}}</h1>
  </body>
</html>
```
