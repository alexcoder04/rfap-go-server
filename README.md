
# rfap-go-server

[![GitHub release](https://img.shields.io/github/v/release/alexcoder04/rfap-go-server?include_prereleases)](https://github.com/alexcoder04/rfap-go-server/releases/latest)
[![GitHub top language](https://img.shields.io/github/languages/top/alexcoder04/rfap-go-server)](https://github.com/alexcoder04/rfap-go-server/search?l=go)
[![License](https://img.shields.io/github/license/alexcoder04/rfap-go-server)](https://github.com/alexcoder04/rfap-go-server/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/alexcoder04/rfap-go-server)](https://github.com/alexcoder04/rfap-go-server/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/alexcoder04/rfap-go-server)](https://github.com/alexcoder04/rfap-go-server/pulls)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/m/alexcoder04/rfap-go-server)](https://github.com/alexcoder04/rfap-go-server/commits/main)
[![GitHub contributors](https://img.shields.io/github/contributors-anon/alexcoder04/rfap-go-server)](https://github.com/alexcoder04/rfap-go-server/graphs/contributors)

The reference server implementation of the rfap protocol, written in GO.
It shares a local folder, which can be then accessed over the network using an
rfap client.

See [here](#related-projects) for protocol specifications and related projects.

## Installation

### Stable release

Simply download the binary for your OS from [the releases
page](https://github.com/alexcoder04/rfap-go-server/releases/latest).

### Bleeding-edge

Make sure you have `git`, `make` and `go` installed.

```sh
git clone https://github.com/alexcoder04/rfap-go-server
cd rfap-go-server

make run       # start testing server
make linux     # compile linux executable
make windows   # compile windows executable
make raspberry # compile linux arm executable
make mac-intel # compile mac intel executable
make install   # compile and install executable to $GOPATH/bin
```

Please use `make` to compile the server, because it tells `go` to inject build
information into the executable which is then useful for understanding logs.

## Related projects

 - https://github.com/alexcoder04/rfap - general protocol specification
 - https://github.com/alexcoder04/librfap - Python client library
 - https://github.com/BoettcherDasOriginal/rfap-cs-lib - C# client library
 - https://github.com/alexcoder04/rfap-pycli - Python CLI client based on librfap
 - https://github.com/alexcoder04/rfap-fuse - FUSE filesystem based on librfap

## Contributing

We appreciate any kind of contribution! Check out
[CONTRIBUTING.md](./CONTRIBUTING.md) for more info.

