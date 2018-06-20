# Saw

`saw` is a multi-purpose tool for tailing AWS CloudWatch Logs

## Fake Praise

 - "It's like tail -f for CloudWatch!"
 - "Yesss, I don't have to install python just to use this"
 - "This guy gets it, pretty print that JSON son"

![Saw Gif](https://media.giphy.com/media/3fiohCfMJAKf7lhnPp/giphy.gif)

## Installation

### Mac OS X

```sh
brew tap TylerBrock/saw
brew install saw
```

### Linux

#### Arch Linux

```sh
# Using pacaur
pacaur -S saw

# Using trizen
trizen -S saw

# Using yaourt
yaourt -S saw

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

## Features

- Colorized output that can be formatted in various ways
    - `--expand` Explode JSON objects using indenting
    - `--rawString` Print JSON strings instead of escaping ("\n", ...)
    - `--invert` Invert white colors to black for light color schemes

- Filter logs using CloudWatch patterns
    - `--filter foo` Filter logs for the text "foo"

- Watch aggregated interleaved streams across a log group
    - `saw watch production` Stream logs from production log group
    - `saw watch production --prefix api` Stream logs from production log group with prefix "api"

TODO:

- Relative or Absolute start and end time specification
    - `saw dump --start 2017-01-01` Stream logs starting from the start of 2017
    - `saw dump --start -1m` Steam logs starting 1 minute ago
    - `saw dump --start -3h --end -2h` Stream logs from 3 - 2 hours ago
