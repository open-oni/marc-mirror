# marc-mirror

This repo just holds some development-friendly MARC XML in case LC is unreachable or slow

For the adventurous, you can mirror the MARC from LC by running `getlc.go`:

    go run getlc.go

It can also be built into a binary if that's desired via `go build`, but that
isn't necessary since it's a simple one-file "script".
