# gatewayd

*gatewayd* is a neat server that lets you to casually set up a Linux-based remote desktop in the web browser.
This is a core part of a bigger project designed in a service-oriented manner.

## Status

Development is not active anymore (for now).
The project is more or less usable, however there are some missing parts here and there.
If you would like to use it production - contact me, I'll guide you through.

## Purpose

This server is used to manage desktop sessions to virtual desktops and make them available via the modern HTML5 web browser.
Originally it was developed for the needs of Moscow State Industrial University to ease the access to the familiar educational environment (GNU/Linux OS with a bunch of preinstalled software) for the freshmen students from home. (We have awesome diskless workstations in classes :P)

This server sits right in between the *remote access server* and *virtual desktop server*. Not a very good explanation, unless you know what these are, but I'll just leave it like that for now.
In serves session creation requests, manages virtual desktop life cycle on a virtual desktop server, and accepts connections from web-based VNC clients to the desktop via WebSockets.
A more in-depth documentation is available in russian, contact me if you have an interest in it.

## Installation

This project currently uses non-go-getable package names.
It is mostly due to legacy reasons. I know this is bad, sorry.

You have to install this project not like it's done usually with `go get`, but with symlinking the `/src` dir to your `$GOPATH/src/gatewayd`.

Use `goexpose` (https://github.com/MOZGIII/goexpose) for that:

```
$ go get github.com/MOZGIII/goexpose
$ git clone https://github.com/MOZGIII/gatewayd.git
$ cd gatewayd/src
$ go get ./...
$ goexpose .
```

Then you can just `go run main.go` and it should not complain.
See further usage instructions in the provided help message (`go run main.go -help`).

## Credits

Designed and implemented by MOZGIII (https://github.com/MOZGIII), review and assistance by Alexey Vereshchagin.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
