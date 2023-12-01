# Pixo Platform CLI

## Prerequisites
- [Pixo Account](https://apex.pixovr.com)

API Key or Username/Password are needed to authenticate with the Pixo Platform APIs.



## Installation
### MacOS - HomeBrew
```bash
brew tap PixoVR/pixo-golang-clients
brew install pixo-cli
pixo help
```

### Windows - Snap
Unfortunately the Pixo CLI is not yet available on Windows via a package manager.
The CLI can be installed by downloading the latest release from the [releases page](https://github.com/PixoVR/pixo-golang-clients/releases)
or building from source.
```
git clone git@github.com:PixoVR/pixo-golang-clients.git
cd pixo-golang-clients/pixo-platform/cmd/platform-cli
make build
./bin/pixo help
```

## Configuration
Credentials can be configured using flags, environment variables or a configuration file.  

Configurations are prioritized in the following order:
1. Flags
2. Environment Variables
3. Configuration File

### Sample Environment Variables:
```bash
# API Key
export SECRET_KEY=secretkey

# Username/Password
export PIXO_USERNAME=username
export PIXO_PASSWORD=password

# API URLs
export PIXO_PLATFORM_API_URL=https://primary.apex.dev.pixovr.com
export PIXO_LEGACY_API_URL=https://api.apex.dev.pixovr.com
export PIXO_MATCHMAKING_API_URL=wss://match.apex.dev.pixovr.com
```

### Sample Configuration File:
```yaml
# ~/.pixo/config.yaml

# API Key
secret-key: secretkey

# Username/Password
username: username
password: password

# API URLs
platform-api-url: https://primary.apex.dev.pixovr.com
legacy-api-url: https://api.apex.dev.pixovr.com
matchmaking-api-url: wss://match.apex.dev.pixovr.com

# Default Module ID
module-id: 271
```

## Initialization
```bash
pixo init # Initialize the configuration file
pixo auth login # Authenticate with the Pixo Platform API
```

## Deploying a Module Gameserver Version
```bash
# Check if server version with semantic version already exists
pixo mp serverVersions deploy \
    --pre-check \
    --module-id 1 \
    --server-version 1.00.00

# Deploy a new version
pixo mp serverVersions deploy \
    --module-id 1 \
    --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest
```

## Testing Multiplayer Matchmaking
```bash
# Request a match
pixo mp matchmake \
    --module-id 1 \
    --server-version 1.00.00

# Request a match and connect to the game server
pixo mp matchmake \
    --module-id 1 \
    --server-version 1.00.00 \
    --connect

# If a match was previously found, the gameserver connection will be saved and can be used to reconnect
pixo mp --connect

# Enter message to send to gameserver: hello
# ACK: hello
```

## Mock Matchmaking Server

Run a mock matchmaking server to test matchmaking functionality locally.
It has a single websocket endpoint, `/matchmaking/matchmake`, that accepts a single message (which it ignores)
and sends a single message back containing a reference to the game server to connect to.

```bash
# Run the server - Ctrl-C to stop
pixo mp matchmake mockserver

# Expected output:
Mon, 02 Jan 2006 15:04:05 MST INF Starting mock server on port 8080
```

To customize the response, use the command line flags or create a yaml file at `./.pixo/server.yaml`.

Defaults to the equivalent of the following config file:
```yaml
server-port: 8080
gameserver-ip: 127.0.0.1
gameserver-port: 7777
map-name: Default
session-name: Test
session-id: FB0HIFBMY8NAME99IS7C3WALKERB4D76
owning-user-name: PixoServer
server-version: 1.00.00
module-id: 1
org-id: 1
```

