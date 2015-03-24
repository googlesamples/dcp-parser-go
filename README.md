DCP
===

Digital Cinema Package parser written in [Go](http://wwww.golang.org)

For more info on DCPs, have a look on [Wikipedia](http://en.wikipedia.org/wiki/Digital_Cinema_Package)

Serves as a good example of how Go's core XML parsers work.

The package includes parsers for the following DCP XML documents:

* Asset Maps
* CPLs
* PKLs

The parsers support both SMPTE and INTEROP packages.

It also includes an example command line program that prints out basic DCP info.

Parsers
-------

To parse a DCP, use:

```go
dcp, err := dcp.New(os.Args[1])
```

To parse individual XML docs, use ParseXXX or ParseXXXFile:

```go
cpl, error := cpl.ParseCPLFile(filepath string)
```

or

```go
cpl, error := cpl.ParseCPL(xmlDoc []byte)
```

Run from the Command Line
-------------------------
To run the DCP inspector from the command line use:

```bash
go run cmd/main.go <path to a dcp root directory>
```
