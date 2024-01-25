# Pixo Platform CLI
This is a CLI that allows you to interact with the Pixo Platform in various ways. It is likely to be most helpful to  
developers building modules to be deployed on the Pixo Platform. It can be used for things like authenticating  
with the platform, deploying gameserver versions, and simplifying the testing of multiplayer components.

## Prerequisites
- [Pixo Account](https://apex.pixovr.com)

API Key or Username/Password are needed to authenticate with the Pixo Platform APIs.


## Installation
### MacOS - HomeBrew
```bash
brew tap PixoVR/pixo-golang-clients
brew install pixo-cli

pixo help

# Update
brew tap PixoVR/pixo-golang-clients
```

### Windows
Unfortunately the Pixo CLI is not yet available on Windows via a package manager.
The CLI can be installed by downloading the latest release from the [releases page](https://github.com/PixoVR/pixo-golang-clients/releases)
or building from source.
```
git clone git@github.com:PixoVR/pixo-golang-clients.git
cd pixo-golang-clients/pixo-platform/cmd/platform-cli
make build
./bin/pixo help
```


## Initialization
```bash
pixo init # Initialize the configuration file
pixo auth login # Authenticate with the Pixo Platform API
```

## Configuration
Configurations can be set using flags, environment variables or a config file.

Configurations are prioritized in the following order:
1. Flags
2. Environment Variables
3. Local Configuration File `./.pixo/config.yaml`
4. Global Configuration File `~/.pixo/config.yaml`

### Sample Environment Variables:
```bash
# API Key
export SECRET_KEY=secretkey

# Username/Password
export PIXO_USERNAME=username
export PIXO_PASSWORD=password

# Pixo Platform Used - if not set, defaults to na and prod
export PIXO_LIFECYCLE=stage
export PIXO_REGION=na
```

### Sample Configuration File:
```yaml
# ~/.pixo/config.yaml

# API Key
secret-key: secretkey

# Username/Password
username: username
password: password

# Platform Used - if not set, defaults to na and prod
lifecycle: stage
region: na

# Default Module ID
module-id: 1
```

## Set via Command Line
```bash
# Requires logging in again after switching environments
pixo config set --region saudi # Switch to saudi environment
pixo config set --lifecycle dev # Switch to dev environment
pixo config set --key module-id --val 1 # Set default module id
```

## Create a User
```bash
pixo users create \
    --username testuser \
    --password testpassword \
    --first-name Test \
    --last-name User \
    --org-id 1 \
    --role developer
    
```

## API Keys

### Create an API Key
```bash
pixo apiKeys create

# Or for a specific user
pixo apiKeys create --user-id 1
```

### List API Keys
```bash
pixo apiKeys list

# Or for a specific user
pixo apiKeys list --user-id 1
```

### Delete an API Key
```bash
pixo apiKeys delete --key-id 1
```

## Deploy a Module Game Server Version
```bash
# Check if server version with semantic version already exists
pixo mp serverVersions deploy \
    --pre-check \
    --module-id 1 \
    --server-version 1.00.00

# Deploy a new version - used in CI/CD pipelines or to test with a simple server like below
pixo mp serverVersions deploy \
    --module-id 1 \
    --server-version 1.00.00 \
    --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest
```

## Tail logs of a Game Server Build
```bash
pixo logs build --module-id 1
```

## Test Multiplayer Matchmaking
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

## Run Mock Matchmaking Server

Run a mock matchmaking server to test matchmaking functionality locally.
It has a single websocket endpoint, `/matchmaking/matchmake`, that accepts a message (which it ignores)
and sends a message back containing the IP and port of the game server to connect to.

```bash
# Run the server - Ctrl-C to stop
pixo mp mockserver

# Expected output:
Mon, 02 Jan 2006 15:04:05 MST INF Starting mock server serving endpoint matchmaking/matchmake on port 8080
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

You can even use the Pixo CLI to test the mock server (run [simple agones server](https://github.com/PixoVR/multiplayer-gameservers/tree/dev/simple-server) locally to test the example below):
```bash
# In one terminal, run the mock server
pixo mp mockserver

# In another terminal, request a match
pixo config set --lifecycle local
pixo mp matchmake \
    --module-id 1 \
    --server-version 1.00.00 \
    --connect
```
