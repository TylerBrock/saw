saw
====

`saw` is a command line tool for shredding Amazon CloudWatch logs.

Installation
------------

### Mac OS X (homebrew)

```sh
brew tap TylerBrock/saw
brew install saw
```

### Linux (binary)

```sh
wget <path_to_binary_release>
tar -zxvf <name_of_binary_release>
mv saw /usr/local/bin
```

Usage
-----

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

Features
--------

- Colorized output that can be formatted in various ways
    - `--expand` Explode JSON objects using indenting
    - `--rawString` Print JSON strings instead of escaping ("\n", ...)
    - `--invert` Invert white colors to black for light color schemes
    - `--no-color` Disable color output entirely

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
