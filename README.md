# vaskafka

### How to use
```go
package main

import "github.com/software-engr-full-stack/vaskafka"

func main() {
    email := "abc@example.com"
    err := vaskafka.Produce("user-created", email)
    if err != nil {
        panic(err)
    }
}
```
