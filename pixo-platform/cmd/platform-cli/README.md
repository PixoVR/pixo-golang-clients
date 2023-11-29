# Pixo Platform CLI

## Prerequisites
- [Pixo Account](https://apex.pixovr.com)

API Key or Username/Password are needed to authenticate with the Pixo Platform API.

## Configuration
Credentials can be configured using environment variables or a configuration file.  

Configurations are loaded in the following order:
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
platform-api-url : https://primary.apex.dev.pixovr.com
legacy-api-url: https://api.apex.dev.pixovr.com
matchmaking-api-url : wss://match.apex.dev.pixovr.com

module-id : 271
```

## Installation
```bash
brew tap PixoVR/pixo-golang-clients
brew install pixo-cli
pixo help
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
pixo mp matchmake --module-id 1 --server-version 1.00.00

# Request a match and connect to the game server
pixo mp matchmake --connect --module-id 1 --server-version 1.00.00

# If a match was previously found, the gameserver connection will be saved and can be used to reconnect
pixo mp --connect

# Enter message to send to gameserver: hello
# ACK: hello
```
