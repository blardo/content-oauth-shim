# Content OAuth Shim

A Google Oauth shim in front of an asset root.

## Go must be installed

```
brew install go
mkdir -p ~/go/{src,bin,pkg}
export GOPATH=~/go
# Append GOPATH to profile
echo 'export GOPATH=~/go' | tee -a ~/.profile
```

## Installation

```bash
mkdir -p $GOPATH/src/github.com/skuid/content-oauth-shim
git clone git@github.com:skuid/content-oauth-shim.git $GOPATH/src/github.com/skuid/content-oauth-shim
cd $GOPATH/src/github.com/skuid/content-oauth-shim

# Install dependencies

go get ./...
```

## Build

```
go build
```

## Usage

By default, the server listens on port 3000.

```
Usage of content-oauth-shim:

  --app-name string
    	The application name to use (default "Oauth Application")
  --assets string
    	The asset root to use (default "./assets")
  --client-id string
    	The Client ID to use
  --client-secret string
    	The Client Secret to use
  --domain value
    	Email domain that authenticated users must match. Can be specified multiple times, or once separating domains with a comma.
        Not listing domains will result in any google user being allowed to authenticate.
  --port int
    	The port to listen on (default 3000)
  --redirect-host string
    	The scheme://host[:port] for Google to return to
  --redis string
    	The redis connection param to use (default "192.168.99.100:6379")
  --secret string
    	The session secret to use
  --tls-cert string
    	The path to the TLS cert to use
  --tls-key string
    	The path to the TLS key to use. tls-key and tls-cert must be used to serve TLS
  --version
    	print version and exit

```

## TODO

* Tests/CI
* Default to FileSystem cookie store, switch to redis if specified
* Have a 401 "Not authorized" for invalid login
* Redirect to asset root on successful login
* Make the build process bundle a `login.html` template into the binary

## License

MIT (See [License](/LICENSE))
