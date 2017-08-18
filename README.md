SAWS
====

`saws` is a command line tool for shredding Amazon CloudWatch logs.

Features
--------

- Colorized output that can be formatted in various ways
    - `saws --expand` -- Explode JSON objects using indenting
    - `saws --rawString` -- Print JSON strings instead of escaping ("\n", ...)
    - `saws --invert` -- Invert white colors to black for light color schemes
    - `saws --no-color` -- Disable color output entirely

- Filter logs using CloudWatch patterns
    - `saws --filter foo` -- Filter logs for the text "foo"

- Aggregate interleaved streams across a log group
    - `saws --group production` -- Stream logs from production log group
    - `saws --group production --prefix api` -- Stream logs from production log group with prefix "api"

- Relative or Absolute start and end time specification
    - `saws --start 2017-01-01` -- Stream logs starting from the start of 2017
    - `saws --start -1m` -- Steam logs starting 1 minute ago
    - `saws --start -3h --end -2h` -- Stream logs from 3 - 2 hours ago
