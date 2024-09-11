![Pixo logo](platform-cli/assets/logo.png)

# Pixo Platform Golang SDK
A collection of Golang clients to interact with the Pixo Platform.

## Clients

- [Pixo CLI](platform-cli/README.md)
- [Platform Client](#platform-client)
- [Legacy Headset Client](#legacy-headset-client)
- [Matchmaking Client](#matchmaking-client)
- [Allocator Client](#allocator-client)
- [Heartbeat Client](#heartbeat-client)


### Platform Client

The Platform Client is a Golang client that interacts with the Pixo Platform API.
It provides a simple interface to interact with the Pixo Platform API.

```go
package main

import (
    "context"
    "fmt"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
    "os"
)

func main() {
    // Create a new client using an API key
    client := platform.NewClient(urlfinder.ClientConfig{
        APIKey: os.Getenv("PIXO_API_KEY"),
        Region: os.Getenv("PIXO_REGION"),
    })
    fmt.Print(client.IsAuthenticated()) // true

    // Create a new client using basic auth
    client, err := platform.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), config)
    if err != nil {
        fmt.Println(err) // invalid login
        return
    }

    fmt.Print(client.IsAuthenticated()) // true

    modules, err := client.GetModules(context.Background())
    if err != nil {
        fmt.Println(err) // error getting modules
        return
    }

    fmt.Println(modules)
}
```

### Legacy Headset Client

The Legacy Headset Client is a Golang client that interacts with the Legacy Headset API

```go
package main

import (
    "context"
    "fmt"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/headset"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
    "os"
)

func main() {
    // Create a new client using an API key
    client := headset.NewClient(urlfinder.ClientConfig{
        APIKey: os.Getenv("PIXO_API_KEY"),
        Region: os.Getenv("PIXO_REGION"),
    })
    fmt.Print(client.IsAuthenticated()) // true

    // Create a new client using basic auth
    client, err := headset.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), config)
    if err != nil {
        fmt.Println(err) // invalid login
        return
    }

    fmt.Print(client.IsAuthenticated()) // true

    // Start session
    response, err := client.StartSession(context.Background(), headset.EventRequest{ModuleID: 1})
    if err != nil {
        fmt.Println(err) // error getting headset info
    }

    fmt.Println(response) // start session event response
}
```

### Matchmaking Client

The Matchmaking Client is a Golang client that interacts with the Matchmaking Service to
find and interact with a gameserver

```go
package main

import (
    "fmt"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/matchmaker"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
    "os"
)

func main() {
    // Create a new client using an API key
    client := matchmaker.NewClient(urlfinder.ClientConfig{
        APIKey: os.Getenv("PIXO_API_KEY"),
        Region: os.Getenv("PIXO_REGION"),
    })
    fmt.Print(client.IsAuthenticated()) // true

    // Create a new client using basic auth
    client, err := matchmaker.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), config)
    if err != nil {
        fmt.Println(err) // invalid login
        return
    }

    fmt.Print(client.IsAuthenticated()) // true

    // Find a gameserver
    addr, err := client.FindMatch(matchmaker.MatchRequest{
        ModuleID:      1,
        ServerVersion: "1.00.00",
    })
    if err != nil {
        fmt.Println(err) // error finding gameserver
    }

    fmt.Println(addr) // udp address of gameserver
}
```

### Allocator Client

The Allocator Client is a Golang client that interacts with the Allocator Service to
allocate gameservers and interact with them

```go
package main

import (
    "fmt"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/allocator"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
    "os"
)

func main() {
    // Create a new client using an API key
    client := allocator.NewClient(urlfinder.ClientConfig{
        APIKey: os.Getenv("PIXO_API_KEY"),
        Region: os.Getenv("PIXO_REGION"),
    })
    fmt.Print(client.IsAuthenticated()) // true

    // Create a new client using basic auth
    client, err := allocator.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), config)
    if err != nil {
        fmt.Println(err) // invalid login
        return
    }

    fmt.Print(client.IsAuthenticated()) // true
	
    // Allocate a gameserver
    gameserver, err := client.AllocateGameserver(allocator.AllocationRequest{
        OrgID:         1,
        ModuleID:      1,
        ServerVersion: "1.00.00",	
    })
    if err != nil {
        fmt.Println(err) // error allocating gameserver
    }
	
    fmt.Println(gameserver) // gameserver info
}
```

### Heartbeat Client

The Heartbeat Client is a Golang client that simplifies sending pulses to the 
Heartbeat Service to keep the session alive

```go
package main

import (
    "fmt"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/heartbeat"
    "github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
    "os"
)

func main() {
    // Create a new client using an API key
    client := heartbeat.NewClient(urlfinder.ClientConfig{
        APIKey: os.Getenv("PIXO_API_KEY"),
        Region: os.Getenv("PIXO_REGION"),
    })
    fmt.Print(client.IsAuthenticated()) // true

    // Create a new client using basic auth
    tokenClient, err := heartbeat.NewClientWithBasicAuth(os.Getenv("PIXO_USERNAME"), os.Getenv("PIXO_PASSWORD"), config)
    if err != nil {
        fmt.Println(err) // invalid login
        return
    }

    fmt.Print(tokenClient.IsAuthenticated()) // true

    // Send a pulse
    if err = tokenClient.SendPulse(sessionID); err != nil {
        fmt.Println(err) // error sending pulse
    }
}
```

