Slackmac
=========

_A tiny HTTP proxy that validates Slack payloads._

SlackMac is a tiny HTTP proxy that validates Slack payloads.
I wrote this because Spring Boot has an issue where it is impossible
to get the raw request body parameters in the correct order if a POST
request is sent with Content-Type `application/x-www-form-urlencoded`.

By offloading this work to a proxy, SlackMac can be dropped in front
of any service that needs to validate Slack payloads without the developer
ever having to worry about calculating HMACs; it's already done!

See [Verifying Requests From Slack](https://api.slack.com/docs/verifying-requests-from-slack)
for more information about the general process for calculating Slack's HMAC implementation. 

## Usage
```
SlackMac is a tiny HTTP proxy that validates Slack payloads.
I wrote this because Spring Boot has an issue where it is impossible
to get the raw request body parameters in the correct order if a POST
request is sent with Content-Type application/x-www-form-urlencoded.
By offloading this work to a proxy, SlackMac can be dropped in front
of any service that needs to validate Slack payloads without the developer
ever having to worry about calculating HMACs. It's already done!

Usage:
  slackmac [flags]

Flags:
  -c, --config string   config file
  -h, --help            help for slackmac
  -v, --verbose         verbose level logging
```

## Config
See [config.go](./config/config.go) for the defaults. Slackmac accepts both
`json` and `toml` config files.

### Stores
Slackmac uses the concept of a `Store` which describes the backend used to 
retrieve the Slack token. The following backends are currently supported:

* Config: The token is stored in your config file. This is the least secure 
option but the easiest. 
```toml
[store]
type = "config"
key = "store.secret"
secret = "THIS IS MY SLACK SIGNING SECRET"
```
* Propsd: The token is retrieved using [Propsd](https://github.com/rapid7/propsd).
```toml
[store]
type = "propsd"
key = "slack.secret"
```
* SecretsManager: The token is retrieved and decrypted from [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
```toml
[store]
type = "secretsmanager"
region = "us-east-1"
id = "the secret ID"
```
* KMS: The token is retrieved and decrypted from [AWS KMS](https://aws.amazon.com/kms/)
```toml
[store]
type = "kms"
region = "us-east-1"
ciphertext = "The CiphertextBlob value that KMS returns when encrypting"
```

Stores are easy to build. They must comply with the `store.Store` interface
which implements a `Get() string` signature.

They should then be registered with the `StoreFactory` by updating the 
`store.init()` function.

See the [store](./store) package for more info.

## Development
We use [dep](https://github.com/golang/dep) to manage dependencies.
You can install it via

```bash
$ go get -u github.com/golang/dep/cmd/dep
```

or, on macOS

```bash
$ brew install dep
$ brew upgrade dep
```
Once you clone the repo, make sure to run `dep ensure` to pull down
the project's (minimal) dependencies.

## Building
You can build Slackmac with any Golang build tool. We prefer
using [Gox](https://github.com/mitchellh/gox). It's simple to use:

```bash
$ go get github.com/mitchellh/gox
$ gox -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

Number of parallel builds: 7

-->   freebsd/amd64: github.com/davepgreene/slackmac
-->       linux/arm: github.com/davepgreene/slackmac
-->      netbsd/arm: github.com/davepgreene/slackmac
-->    darwin/amd64: github.com/davepgreene/slackmac
-->     freebsd/386: github.com/davepgreene/slackmac
-->       linux/386: github.com/davepgreene/slackmac
-->     linux/amd64: github.com/davepgreene/slackmac
-->   windows/amd64: github.com/davepgreene/slackmac
-->     openbsd/386: github.com/davepgreene/slackmac
-->   openbsd/amd64: github.com/davepgreene/slackmac
-->     windows/386: github.com/davepgreene/slackmac
-->      netbsd/386: github.com/davepgreene/slackmac
-->     freebsd/arm: github.com/davepgreene/slackmac
-->    netbsd/amd64: github.com/davepgreene/slackmac
-->      darwin/386: github.com/davepgreene/slackmac
