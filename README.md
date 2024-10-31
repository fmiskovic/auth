![go workflow](https://github.com/softika/auth/actions/workflows/test.yml/badge.svg)
![lint workflow](https://github.com/softika/auth/actions/workflows/lint.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/softika/auth)](https://goreportcard.com/report/github.com/softika/auth)


# JWT Auth Middleware

This Go library provides authentication middleware and utilities for handling JWT tokens. 
It includes functionality for validating tokens, extracting claims, and managing user context.
If the given token is valid, it extracts the claims and adds them to the request context.

## Installation

To install the library, use `go get`:

```sh
go get github.com/softika/auth
```

## Examples

### Standard Library
```go
package main

import (
    "log"
    "net/http"

    "github.com/softika/auth"
)

func main() {
    mux := http.NewServeMux()

    // init auth middleware
    a := auth.New(auth.Config{
        Secret: "your-secret",
    })

    // wrap the handler with the auth middleware
    mux.Handle("/protected", a.Handler(http.HandlerFunc(HandleProtectedInfo)))

    log.Fatal(http.ListenAndServe(":3000", mux))
}

func HandleProtectedInfo(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Protected Info!"))
}
```

### Go Chi

```go
package main

import (
    "log"
    "net/http"

    "github.com/go-chi/chi"
    "github.com/softika/auth"
)

func main() {
    r := chi.NewRouter()

    // init auth configuration
    cfg := auth.Config{
        Secret: "your-secret",
    }

    // wrap the handler with the auth middleware
    r.Use(auth.Handle(cfg))
    r.Get("/protected", HandleProtectedInfo)

    log.Fatal(http.ListenAndServe(":3000", r))
}

func HandleProtectedInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Protected Info!"))
}
```




## Contributing

Contributions are welcome! Please open an issue or submit a pull request.