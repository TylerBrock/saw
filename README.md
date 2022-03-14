# Saw

`saw` is a multi-purpose tool for AWS CloudWatch Logs

![Saw Gif](https://media.giphy.com/media/3fiohCfMJAKf7lhnPp/giphy.gif)

[![Go Report Card](https://goreportcard.com/badge/github.com/TylerBrock/saw)](https://goreportcard.com/report/github.com/TylerBrock/saw)

## Features

- Colorized output that can be formatted in various ways
    - `--expand` Explode JSON objects using indenting
    - `--rawString` Print JSON strings instead of escaping ("\n", ...)
    - `--invert` Invert white colors to black for light color schemes
    - `--raw`, or `--pretty`, for `watch` and `get` commands respectively, toggles display of the timestamp and stream name prefix.

- Filter logs using CloudWatch patterns
    - `--filter foo` Filter logs for the text "foo"

- Watch aggregated interleaved streams across a log group
    - `saw watch production` Stream logs from production log group
    - `saw watch production --prefix api` Stream logs from production log group with prefix "api"

## Usage

- Basic
    ```sh
    # Get list of log groups
    saw groups

    # Get list of streams for production log group
    saw streams production
    ```

- Watch
    ```sh
    # Watch production log group
    saw watch production

    # Watch production log group streams for api
    saw watch production --prefix api

    # Watch production log group streams for api and filter for "error"
    saw watch production --prefix api --filter error
    ```

- Get
    ```sh
    # Get production log group for the last 2 hours
    saw get production --start -2h

    # Get production log group for the last 2 hours and filter for "error"
    saw get production --start -2h --filter error

    # Get production log group for api between 26th June 2018 and 28th June 2018
    saw get production --prefix api --start 2018-06-26 --stop 2018-06-28
    ```

### Profile and Region Support

By default Saw uses the region and credentials in your default profile. You can override these to your liking using the command line flags:

```sh
# Use personal profile
saw groups --profile personal

# Use us-west-1 region
saw groups --region us-west-1
```

## Installation

### Run from Docker

```sh
docker run --rm -it -v ~/.aws:$HOME/.aws tbrock/saw
```

### Mac OS X

```sh
brew tap TylerBrock/saw
brew install saw
```

### Linux

#### Arch Linux (source)

```sh
# Using yay
yay saw

# Using makepkg
git clone https://aur.archlinux.org/saw.git
cd saw
makepkg -sri
```

#### Red Hat Based Distributions (Fedora/RHEL/CentOS/Amazon Linux)
```sh
rpm -i <link_to_rpm_you_need_from_releases>
```

#### Debian Based Distributions (Debian/Ubuntu)
```sh
wget <link_to_deb_you_need_from_releases>
sudo dpkg -i <the_deb_name>
```

### Manual Install/Update

- [Install go](https://golang.org/doc/install)
- Configure your `GOPATH` and add `$GOPATH/bin` to your path
- Run `go get -u github.com/TylerBrock/saw`

#### Windows Specifics

- Add %GOPATH%/bin to your path (optional)
- Run from gopath/bin (If not in your path)
    ```DOS .bat
    cd %GOPATH%/bin
    saw ...
    ```

Alternatively you can hard code these in your shell's init scripts (bashrc, zshrc, etc...):

```sh
# Export profile and region that override the default
export AWS_PROFILE='work_profile'
export AWS_REGION='us-west-1'
```

## Run Tests
From root of repository: `go test -v ./...`

## TODO

- Bash + ZSH completion of log groups + (streams?)
- Create log streams and groups
- Delete log streams and groups
- Basic tests
