![Pixo logo](assets/logo.png)

[//]: # ([![]&#40;https://img.shields.io/github/actions/workflow/status/spf13/cobra/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff&#41;]&#40;https://github.com/spf13/cobra/actions?query=workflow%3ATest&#41;)
[//]: # ([![Go Report Card]&#40;https://goreportcard.com/badge/github.com/spf13/cobra&#41;]&#40;https://goreportcard.com/report/github.com/spf13/cobra&#41;)

# Pixo Platform CLI
`pixo` is a CLI that allows you to interact with the Pixo Platform in various ways. It is likely to be most helpful to 
developers building modules to be deployed on the Pixo Platform. It can be used for things like authenticating 
with the platform, deploying gameserver versions, and simplifying the testing of multiplayer components.

## Prerequisites
- [Pixo Account](https://apex.pixovr.com) with API Key or Username/Password

## Table of Contents
- [Installation](#installation)
    - [MacOS - HomeBrew](#macos---homebrew)
    - [Windows](#windows)
    - [Build from Source](#build-from-source)
    - [Autocompletion](#autocompletion)
- [Configuration](#configuration)
    - [Set via Environment Variables](#set-via-environment-variables)
    - [Set via Command Line](#set-via-command-line)
    - [Show Configuration File](#show-configuration-file)
    - [Edit Configuration File](#edit-configuration-file)
    - [Get Platform Service URLs](#get-platform-service-urls)
- [Login to the Pixo Platform](#login-to-the-pixo-platform)
- [Users](#users)
    - [Create](#create)
- [API Keys](#api-keys)
    - [Create](#create)
    - [List](#list)
    - [Delete](#delete)
- [Modules](#modules)
  - [Create Module Version](#create-module-version)
- [Webhooks](#webhooks)
  - [Create](#create-webhook)
  - [List](#list-webhook)
  - [Delete](#delete-webhook)
- [Sessions](#sessions)
  - [Simulate a Session](#simulate-a-session)
  - [Using Legacy Headset API](#using-legacy-headset-api)
- [Running Mock Servers](#run-mock-servers)
  - [Platform](#platform)
  - [Matchmaking](#matchmaking)
- [Deploy a Module Game Server Version](#deploy-a-module-game-server-version)
    - [Gameserver Build Pipeline (e.g. Cloud Build)](#gameserver-build-pipeline-eg-cloud-build)
        - [Sample Ini Configuration](#sample-ini-configuration)
        - [Sample `cloudbuild.yaml`](#sample-cloudbuildyaml)
- [Test Multiplayer Matchmaking](#test-multiplayer-matchmaking)
    - [Request a Match](#request-a-match)
    - [Connect to the Game Server](#connect-to-the-game-server)
- [Load Testing](#load-testing)
  - [Sessions](#load-test-sessions)
  - [Matchmaking](#load-test-matchmaking)

## Installation
### Go - recommended
```bash
go install github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli@latest
````

### HomeBrew
```bash
brew tap PixoVR/pixo-golang-clients
brew install pixo-cli
```

### Windows
Unfortunately the Pixo CLI is not yet available on Windows via package manager.
The CLI can be installed by downloading the latest release from the [releases page](https://github.com/PixoVR/pixo-golang-clients/releases)
or building from source.

### Download from [Releases page](https://github.com/PixoVR/pixo-golang-clients/releases)

### Build from Source
```bash
git clone github.com/PixoVR/pixo-golang-clients.git
cd pixo-golang-clients/pixo-platform/platform-cli
make build
./bin/pixo help
```

### Autocompletion
```bash
# Zsh
echo "source <(pixo completion zsh)" >> ~/.zshrc
source ~/.zshrc

# Bash
echo "source <(pixo completion bash)" >> ~/.bashrc
source ~/.bashrc
```


## Configuration
Configurations can be set using flags, environment variables or a config file.

Configurations are prioritized in the following order:
1. Flags
2. Environment Variables
3. Local Configuration File `./.pixo/config.yaml`
4. Global Configuration File `~/.pixo/config.yaml`

### Set via Environment Variables:
```bash
# API Key
export PIXO_API_KEY=<api-key>

# Username/Password
export PIXO_USERNAME=<username>
export PIXO_PASSWORD=<password>

# Pixo Platform Environment Used - if not set, defaults to na
export PIXO_REGION=saudi
```

### Set via Command Line
```bash
pixo config set --region saudi
pixo config set --lifecycle dev
pixo config set --key module --val TST
```

### Show Configuration File
```bash
pixo config

# Example output:
📁  Config file: ~/.pixo/config.yaml
🌎  Region: na
⚙️   Lifecycle: dev

🆔  User ID: 1
👤  Username: <username>
🔒  Password: ********
🔑  API Key: ********
🪙  Token: ********


➡️  Module: TST
➡️  Server Version: 1.00.00
➡️  Gameserver: 127.0.0.1:7777
```

### Edit Configuration File
Editor can be set via the `EDITOR` environment variable. Defaults to `vim`.
```bash
pixo config --edit # or -e
```

### Get Platform Service URLs
```bash
pixo config urls

# Example output:
🌎  Region: na
⚙️   Lifecycle: prod

🔗  Web: https://apex.pixovr.com

🔗  Platform API: https://apex.pixovr.com/v2
🔗  Platform API Docs: https://apex.pixovr.com/v2/swagger/index.html

🔗  Matchmaking API: https://apex.pixovr.com/matchmaking
🔗  Matchmaking API Docs: https://apex.pixovr.com/matchmaking/swagger/index.html

🔗  Heartbeat API: https://apex.pixovr.com/heartbeat
🔗  Heartbeat API Docs: https://apex.pixovr.com/heartbeat/swagger/index.html
```


## Login to the Pixo Platform
```bash
pixo auth login

# Or with username and password
pixo auth login --username <username> --password <password>

# Example output:
🚀 Login successful. Here is your API token: 
<token>
```
Token redacted for security reasons.
![Made with VHS](https://vhs.charm.sh/vhs-6w7a14WEgRuEIG3400Zazf.gif)


## Users

### Create
```bash
pixo users create \
    --first-name Test \
    --last-name User \
    --user-email testuser@example.com \
    --user-username testuser \
    --user-password testpassword \
    --org "My Org" \
    --role developer
```
![Made with VHS](https://vhs.charm.sh/vhs-4IQGJes6OQQoWTN8dt3e8A.gif)

## API Keys

### Create
```bash
pixo keys create

# Or for a specific user
pixo keys create --username testuser
```
![Made with VHS](https://vhs.charm.sh/vhs-DBpsz1KVGCMzMHkEgF4Gg.gif)


### List
```bash
pixo keys list

# Or for a specific user
pixo keys list --username testuser
```

### Delete
```bash
pixo keys delete --key-ids 1
```
![Made with VHS](https://vhs.charm.sh/vhs-3vojsRNWUrNJH6lwC7ozQT.gif)

## Modules

### Create Module Version
```bash
pixo modules deploy \
    --module TST \
    --server-version "1.00.00" \
    --package "com.pixovr.test" \
    --platforms "android" \
    --controls "keyboard/mouse" \
    --zip-file "/path/to/zip"
```

## Webhooks

### Create Webhook
```bash
pixo webhooks create \
  --url https://example.com/webhook
  --description "Test Webhook"
```

### List Webhook
```bash
pixo webhooks list
```

### Delete Webhook
```bash
pixo webhooks delete --webhook-id 1
```

## Sessions

### Simulate a Session
```bash
pixo sessions simulate
```
![Made with VHS](https://vhs.charm.sh/vhs-424AzRmHaYF3Re5f6Yl6U.gif)

#### Using Legacy Headset API
```bash
pixo sessions simulate --legacy
```

## Run Mock Servers
### Platform
Run a mock server that mimics the Pixo Platform API to test functionality locally.
It has the following REST endpoints available. See the [Swagger API Docs](https://apex.pixovr.com/v2/swagger/index.html) for more details.

- GET - `/v2/assets`
- GET - `/v2/assets/download`
- POST - `/v2/assets`

### Matchmaking 

Run a mock server that mimics the Pixo Matchmaking Service to test matchmaking functionality locally.
It has a single websocket endpoint, `/matchmaking/matchmake`, that accepts a message (which it ignores)
and sends a message back containing the IP and port of the game server to connect to.

```bash
# Run the server - Ctrl-C to stop
pixo mp mockserver
```
![Made with VHS](https://vhs.charm.sh/vhs-2qOnNiisC4uI3ld6mwUewG.gif)


To customize the response, use the command line flags when starting the server
```bash
pixo mp mockserver --gameserver-port 7654
```

Defaults to the following values:
```yaml
matchmaker-port: 8080
gameserver-ip: 127.0.0.1
gameserver-port: 7777
map-name: Default
session-name: Test
session-id: FB0HIFBMY8NAME99IS7C3WALKERB4D76
owning-user-name: PixoServer
server-version: 1.00.00
module-id: 1
org: 1
```

You could even use the Pixo CLI to test the mock server (run [simple agones server](https://github.com/PixoVR/multiplayer-gameservers/tree/dev/simple-server) locally to test the example below):
```bash
# In one terminal, run the mock server
pixo mp mockserver

# In another terminal, request a match
pixo config set --lifecycle local
pixo mp matchmake \
    --module TST \
    --server-version 1.00.00 \
    --connect
```

## Deploy a Module Game Server Version
```bash
# Check if version with matching semantic version already exists
pixo mp servers deploy \
    --pre-check \
    --module TST \
    --server-version 1.00.00
```

```bash
# Deploy a new version with image
pixo mp servers deploy \
    --module TST \
    --server-version 1.00.00 \
    --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest
```

```bash
# Deploy a new version with zipfile
pixo mp servers deploy \
    --module TST \
    --server-version 1.00.00 \
    --zip-file /path/to/zipfile
```

```bash
# Update existing version
pixo mp servers deploy \
    --update \
    --module TST \
    --server-version 1.00.00 \
    --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:0.0.6
```


### Gameserver Build Pipeline (e.g. Cloud Build)
`server-version` can be configured via an ini file via the flag `--ini`

#### Sample Ini Configuration
```ini
[/Script/PixoConfig.PixoConfigSettings]
ServerMatchVersion=1.00.00
```

#### Sample `cloudbuild.yaml` - replace `latest` with the version you want to use
```yaml
steps:
  - name: "gcr.io/pixo-bootstrap/pixo-platform-cli:latest"
    id: "Version Pre-Check"
    args:
      - mp
      - servers
      - deploy
      - --module
      - ${_MODULE_ABBREVIATION}
      - --pre-check
    env:
      - "PIXO_REGION=${_PIXO_REGION}"
    secretEnv:
      - "PIXO_API_KEY"

  - name: "gcr.io/cloud-builders/docker"
    id: "Build and Push Game Server Image"
    args:
      - build
      - -t
      - gcr.io/${PROJECT_ID}/${_LIFECYCLE}/${_PROJECT_NAME}:latest
      - -t
      - gcr.io/${PROJECT_ID}/${_LIFECYCLE}/${_PROJECT_NAME}:${COMMIT_SHA}

  - name: "gcr.io/pixo-bootstrap/pixo-platform-cli:latest"
    id: "Deploy MP Server Version"
    args:
      - mp
      - servers
      - deploy
      - --module
      - ${_MODULE_ABBREVIATION}
      - --image
      - gcr.io/${PROJECT_ID}/${_LIFECYCLE}/${_PROJECT_NAME}:${COMMIT_SHA}
    env:
      - "PIXO_REGION=${_PIXO_REGION}"
    secretEnv:
      - "PIXO_API_KEY"

availableSecrets:
  secretManager:
    - versionName: projects/$PROJECT_ID/secrets/$SECRET_NAME/versions/latest
      env: PIXO_API_KEY
      
images:
  - gcr.io/${PROJECT_ID}/${_LIFECYCLE}/${_PROJECT_NAME}:latest
  - gcr.io/${PROJECT_ID}/${_LIFECYCLE}/${_PROJECT_NAME}:${COMMIT_SHA}


```


## Test Multiplayer Matchmaking

![Made with VHS](https://vhs.charm.sh/vhs-1vPC2fJWNNr9v9Smnmzshh.gif)

### Request a Match
```bash
# Request a match
pixo mp matchmake \
    --module TST \
    --server-version 1.00.00
```

### Connect to the Game Server
```bash
# Request a match and connect to the game server
pixo mp matchmake \
    --module TST \
    --server-version 1.00.00 \
    --connect
    
# If a match was previously found, the gameserver connection will be saved and can be used to reconnect
pixo mp --connect

# Enter message to send to gameserver: hello
# ACK: hello
```


## Load Testing

### Load Test Sessions
```bash
pixo cannon sessions \
    --module TST \
    --amount 5 \
    --concurrent 2
    
# With session details
pixo cannon sessions \
    --module TST \
    --version "1.00.00" \
    --mode "practice" \
    --scenario "warehouse" \
    --focus "packaging" \
    --specialization "machinery" \
    --score 5 \
    --max-score 10 \
    --passed
    
# With an event payload
pixo cannon sessions \
    --module TST \
    --payload '{"key": "value"}'
    
# Or with an event payload from a file
pixo cannon sessions \
    --module TST \
    --payload-file /path/to/payload.json
    
# Sample output
🚀  Starting load test with 5 requests and 2 concurrent workers

✅  2: session started for module TST
✅  1: session started for module TST
✅  2: event created for session 1
✅  1: event created for session 2
✅  2: session completed for module TST
✅  1: session completed for module TST
✅  3: session started for module TST
✅  3: event created for session 3
✅  4: session started for module TST
✅  4: event created for session 4
✅  3: session completed for module TST
✅  4: session completed for module TST
✅  5: session started for module TST
✅  5: event created for session 5
✅  5: session completed for module TST

Load Test Summary
===========================
Concurrent Workers:     2
Amount Requested:       5
Amount Completed:       5
Max Test Duration:      2m0s
Actual Test Duration:   4.25s

┌───────────────┬────────────┐
│ Stat          │ Value      │
├───────────────┼────────────┤
│ Avg Latency   │ 1.45s      │
│ Max Latency   │ 1.65s      │
│ Req / Sec     │ 1.18       │
└───────────────┴────────────┘

Start Session Errors:           0
Create Event Errors:            0
Complete Session Errors:        0
Unsuccessful Sessions:          0
Sessions Started:               5
Events Created:                 5
Sessions Completed:             5
```

### Legacy Sessions
```bash
pixo cannon sessions \
    --legacy \
    --module TST \
    --amount 5 \
    --concurrent 2
    
🚀  Starting load test with 5 requests and 2 concurrent workers

✅  2: session started for module 1
✅  1: session started for module 1
✅  1: event created for session 1
✅  2: event created for session 2
✅  1: session completed for module 1
✅  2: session completed for module 1
✅  3: session started for module 43
✅  3: event created for session 3
✅  4: session started for module 1
✅  4: event created for session 4
✅  3: session completed for module 1
✅  4: session completed for module 1
✅  5: session started for module 1
✅  5: event created for session 5
✅  5: session completed for module 1

Load Test Summary
===========================
Concurrent Workers:     2
Amount Requested:       5
Amount Completed:       5
Max Test Duration:      2m0s
Actual Test Duration:   6.75s

┌───────────────┬────────────┐
│ Stat          │ Value      │
├───────────────┼────────────┤
│ Avg Latency   │ 2.26s      │
│ Max Latency   │ 2.79s      │
│ Req / Sec     │ 0.74       │
└───────────────┴────────────┘

Start Session Errors:           0
Create Event Errors:            0
Complete Session Errors:        0
Unsuccessful Sessions:          0
Sessions Started:               5
Events Created:                 5
Sessions Completed:             5
```

### Load Test Matchmaking
```bash
pixo cannon matchmake \
    --module TST \
    --server-version 1.00.00 \
    --amount 5 \
    --concurrent 2
    
# Sample output

🚀  Starting load test with 5 requests and 2 concurrent workers

✅  2: Connection established
✅  1: Connection established
🏁  Match found - gameserver -> 34.1.2.3:7728
🏁  Match found - gameserver -> 34.1.2.3:7728
✅  3: Connection established
✅  4: Connection established
🏁  Match found - gameserver -> 34.1.2.3:7728
🏁  Match found - gameserver -> 34.1.2.3:7728
✅  5: Connection established
🏁  Match found - gameserver -> 34.1.2.3:7728

Load Test Summary
===========================
Concurrent Workers:     2
Amount Requested:       5
Amount Completed:       5
Max Test Duration:      2m0s
Actual Test Duration:   16.15s

┌───────────────┬────────────┐
│ Stat          │ Value      │
├───────────────┼────────────┤
│ Avg Latency   │ 6.06s      │
│ Max Latency   │ 12.57s     │
│ Req / Sec     │ 0.31       │
└───────────────┴────────────┘

Connection Errors:         0
Successful Connections:    5
Matching Errors:           0
Matches Received:          5
Gameservers Received:      1
```
