# Hash Embed

`hashembed` is a thin wrapper around [embed.FS](https://pkg.go.dev/embed) to allow accessing files with a content hash.

`hashembed` is useful if you are embedding static assets directly into your application and want to
facilitate serving these files with very long duration client-side caching.

***Note**: You should use middleware instead of `hashembed`. While `hashembed` is functional it was built as an exercise in learning Go.*

# Usage

```go
package main

import (
  "embed"
  "github.com/42z-io/hashembed"
)

//go:embed testdata/*
var embedded embed.FS

func main() {
  embedded, _ := Generate(embedded)
  path := embedded.Reverse("testdata/test.css")
  fmt.Printf(path)
  // Output: testdata/test.8d77f04c3be2abcd554f262130ba6c30f277318e66588b6a0d95f476c4ae7c48.css
  data, _ := embedded.ReadFile(path)
  fmt.Println(string(data[:])
  // Output: body { width: 100%; }
}
```

# Use Case

Here is a psuedo example using [Fiber](https://gofiber.io/), and [Templ](https://templ.guide/)

```go
//go:embed dist/*
var data embed.FS
var hashedData := hashembed.Generate(data)

// Template
templ {
    <script src={ "/static/" + hashedData.Reverse("dist/my_file.css") }>
}
// <script src="/static/dist/my_file.HASH_CODE_HERE.css">

// Filesystem middleware
app.Use("/static", filesystem.New(
  filesystem.Config{
      Root:       http.FS(hashedData),
      Browse:     false,
      MaxAge:     600000,
  }
))
```
