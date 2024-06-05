# CLI Login

Demo CLI application for logging via browser

## Install

1. Add Client ID and Client Secret of your Google application to file  `client/secret.go`

(See https://console.cloud.google.com/apis/credentials for details. Authorized redirect URIs should be set 
to http://localhost:3000/auth/google/callback)

2. Run `make`

or:

1. Get compiled binary from `dist` directory

## Usage

`clilogin status` - show application login status

`clilogin login` - perform login via browser



