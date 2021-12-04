
# rfap-go-server

The reference server implementation of the rfap protocol.

See [here](#related-projects) for protocol specifications and related projects.

## Installation

### For Windows

Assuming you have docker-compose installed and in path.
```
docker-compose build
docker-compose up
```

### Stable release

**Coming soon**

### Bleeding-edge

```sh
git clone https://github.com/alexcoder04/rfap-go-server
cd rfap-go-server/src

make run     # start testing server
make linux   # compile executable, other possible arguments: windows/raspberry
go install . # compile and install executable to $GOPATH/bin
```

## Related projects

 - https://github.com/alexcoder04/rfap - general protocol specification
 - https://github.com/alexcoder04/librfap - Python library
 - https://github.com/BoettcherDasOriginal/rfap-cs-lib - C# library
 - https://github.com/alexcoder04/rfap-pycli - Python CLI client based on librfap

