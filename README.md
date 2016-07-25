# marc-mirror

This repo just holds some development-friendly MARC XML in case LC is unreachable or slow

## Use with Open ONI

Add this line to your `settings_local.py` file:

    MARC_RETRIEVAL_URLFORMAT = "https://raw.githubusercontent.com/open-oni/marc-mirror/master/marc/%s/marc.xml"

This will tell ONI to use this mirrored repo instead of the default URL, which relies on Chronicling America being up.

## getlc.go

For the adventurous, you can mirror the MARC from LC by running `getlc.go`:

    go run getlc.go

It can also be built into a binary if that's desired via `go build`, but that
isn't necessary since it's a simple one-file "script".
