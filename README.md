# DGCP-SDK-GO

Official Go SDK for DGCP APIs.

## Installation
```bash
go get github.com/XYZLZ/dgcp-sdk-go
```

## Quick Start
```go
package main

import (
    "context"
    "log"
    
    sdk "github.com/yourcompany/my-api-sdk-go"
    "github.com/yourcompany/my-api-sdk-go/types"
)

func main() {
    // Initialize SDK
    client := dgcp.New("your-api-key")
    
    ctx := context.Background()
    
    // Use the SDK
    files, err := client.Mahoraga.Files.List(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, user := range users.Payload.Content {
        log.Printf("User: %s (%s)", user.Name, user.Email)
    }
}
```

## Configuration
```go
client := dgcp.New(
    "your-api-key",
    dgcp.WithTimeout(10*time.Second),
    dgcp.WithDebug(true),
    dgcp.WithMaxRetries(5),
    dgcp.WithCustomHeader("X-Custom", "value"),
)
```

## API Endpoints

The SDK automatically routes requests to the correct API:

- **DGCP API**: Processes, Offers, Contracts
- **Mahoraga API**: File operations

<!-- ## Error Handling
```go
user, err := client.Users.Get(ctx, "user-123")
if err != nil {
    switch e := err.(type) {
    case *client.AuthenticationError:
        log.Println("Invalid API key")
    case *client.NotFoundError:
        log.Println("User not found")
    case *client.RateLimitError:
        log.Printf("Rate limited, retry after %d seconds", e.RetryAfter)
    default:
        log.Printf("Error: %v", err)
    }
}
``` -->

## License

MIT