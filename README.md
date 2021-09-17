env_scanner
===========

Find environment reports from [env_reporter](https://github.com/mayth/env_reporter) in BLE advertisements.

NOTICE: This does not works with Darwin (macOS).

## Prerequisites

* [Go](https://golang.org) 1.17+
* BLE broadcaster running [env_reporter](https://github.com/mayth/env_reporter)

## Build

```
$ go build
```

## Testing

```
$ go test .
```

## Run

To access Bluetooth device, you may need to be root or use `sudo`.

```
$ go build -o env_scanner
$ sudo ./env_scanner [-prefix NAME]
```

This app searches advertisement packet that contains the environmental sensing service data.

`-prefix` is optional. If specified, this app filters out the advertisement packets whose local name (short name) does not start with the given prefix.

This app accepts `SIGUSR1` signal. When received the signal, this app shows the current (latest) reports from reporters.
