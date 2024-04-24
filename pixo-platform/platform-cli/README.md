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
  - [Start a Session](#start-a-session)
  - [End a Session](#end-a-session)
- [Run Mock Matchmaking Server](#run-mock-matchmaking-server)
- [Deploy a Module Game Server Version](#deploy-a-module-game-server-version)
    - [Gameserver Build Pipeline (e.g. Cloud Build)](#gameserver-build-pipeline-eg-cloud-build)
        - [Sample Ini Configuration](#sample-ini-configuration)
        - [Sample `cloudbuild.yaml`](#sample-cloudbuildyaml)
- [Test Multiplayer Matchmaking](#test-multiplayer-matchmaking)
    - [Request a Match](#request-a-match)
    - [Connect to the Game Server](#connect-to-the-game-server)
    - [Load Testing](#load-testing)


## Installation
### MacOS - HomeBrew
```bash
brew tap PixoVR/pixo-golang-clients
brew install pixo-cli

pixo help
```

### Windows
Unfortunately the Pixo CLI is not yet available on Windows via package manager.
The CLI can be installed by downloading the latest release from the [releases page](https://github.com/PixoVR/pixo-golang-clients/releases)
or building from source.
```
git clone git@github.com:PixoVR/pixo-golang-clients.git
cd pixo-golang-clients/pixo-platform/platform-cli
make build
./bin/pixo help
```

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
pixo config set --key module-id --val 1
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


➡️  Module ID: 1
➡️  Server Version: 1.00.00
➡️  Gameserver: 127.0.0.1:7777
```

### Edit Configuration File
Editor can be set via the `EDITOR` environment variable. Defaults to `vim`.
```bash
pixo config --edit
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


## Users

### Create
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

### Create
```bash
pixo keys create

# Or for a specific user
pixo keys create --user-id 1
```

### List
```bash
pixo keys list

# Or for a specific user
pixo keys list --user-id 1
```

### Delete
```bash
pixo keys delete --key-id 1
```

## Modules

### Create Module Version
```bash
pixo modules deploy \
    --module-id 1 \
    --server-version 1.00.00 \
    --package com.pixovr.test \
    --zip-file /path/to/zip
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

### Start a Session
```bash
pixo sessions start --module-id 1
```

### End a Session
```bash
# Using current session ID
pixo sessions end \
  --score 1 \
  --max-score 2
```

```bash
# Or with session ID as input
pixo sessions end \
  --session-id 123 \
  --score 1 \
  --max-score 2
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

To customize the response, use the command line flags when starting the server
```bash
pixo mp mockserver \
    --gameserver-port 7654
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
org-id: 1
```

You could even use the Pixo CLI to test the mock server (run [simple agones server](https://github.com/PixoVR/multiplayer-gameservers/tree/dev/simple-server) locally to test the example below):
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

## Deploy a Module Game Server Version
```bash
# Check if version with matching semantic version already exists
pixo mp servers deploy \
    --pre-check \
    --module-id 1 \
    --server-version 1.00.00
```

```bash
# Deploy a new version with image
pixo mp servers deploy \
    --module-id 1 \
    --server-version 1.00.00 \
    --image gcr.io/pixo-bootstrap/multiplayer/gameservers/simple-server:latest
```

```bash
# Deploy a new version with zipfile
pixo mp servers deploy \
    --module-id 1 \
    --server-version 1.00.00 \
    --zip-file /path/to/zipfile
```


### Gameserver Build Pipeline (e.g. Cloud Build)
If no `server-version` configuration value is found, it will search for an ini file  
The ini used file can be set with the flag `--ini` and defaults to `Config/DefaultGame.ini`

#### Sample Ini Configuration
```ini
[/Script/PixoConfig.PixoConfigSettings]
ServerMatchVersion=1.00.00
```

#### Sample `cloudbuild.yaml`
```yaml
steps:
  - name: "gcr.io/pixo-bootstrap/pixo-platform-cli:0.0.177"
    id: "Version Pre-Check"
    args:
      - mp
      - servers
      - deploy
      - --module-id
      - ${_MODULE_ID}
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

  - name: "gcr.io/pixo-bootstrap/pixo-platform-cli:0.0.177"
    id: "Deploy MP Server Version"
    args:
      - mp
      - servers
      - deploy
      - --module-id
      - ${_MP_MODULE_ID}
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
```


## Test Multiplayer Matchmaking

### Request a Match
```bash
# Request a match
pixo mp matchmake \
    --module-id 1 \
    --server-version 1.00.00
```

### Connect to the Game Server
```bash
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


### Load Testing
```
pixo mp matchmake \
    --module-id 1 \
    --server-version 1.00.00 \
    --load 5
    
# Example output
🚀  Starting load test with 5 connections to wss://apex.dev.pixovr.com/matchmaking...

✅  Connection 5: established
✅  Connection 4: established
✅  Connection 2: established
✅  Connection 3: established
✅  Connection 1: established
🏁  Connection 1: Match found - gameserver -> 34.1.2.3:7566
🏁  Connection 2: Match found - gameserver -> 34.1.2.3:7566
🏁  Connection 4: Match found - gameserver -> 34.1.2.3:7566
🏁  Connection 5: Match found - gameserver -> 34.1.2.3:7756
🏁  Connection 3: Match found - gameserver -> 34.1.2.3:7756

Matchmaking Load Test Summary
==============================
Max Test Duration:       10m0s
Actual Test Duration:    13.6s
Connections:             5
Total Messages Sent:     5

Total Messages Received: 5
Connection Errors:       0
Matching Errors:         0
Matches Received:        5
Gameservers Received:    2

┌─────────────┬────────────┐
│ Stat        │ Value      │
├─────────────┼────────────┤
│ Avg Latency │ 6.58 s     │
│ Max Latency │ 13.52 s    │
│ Msgs/Sec    │ 0.37       │
└─────────────┴────────────┘
```
