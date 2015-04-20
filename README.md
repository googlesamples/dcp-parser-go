Go DCP Parser
=============

Digital Cinema Package parser written in [Go](http://wwww.golang.org)

For more info on DCPs, have a look on [Wikipedia](http://en.wikipedia.org/wiki/Digital_Cinema_Package)

Serves as an example of how Go's core XML parsers work.

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
dcp, err := dcp.New(<path to root of a DCP folder>)
```

To parse individual XML docs, use ParseXXX() or ParseXXXFile():

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
go run cmd/main.go <path to root of a DCP folde>
```

Support
-------

If you've found an error in this sample, please file an issue:
https://github.com/googlesamples/dcp-parser-go/issues

Patches are encouraged, and may be submitted by forking this project and
submitting a pull request through GitHub.

License
-------

Copyright 2015 Google, Inc.

Licensed to the Apache Software Foundation (ASF) under one or more contributor
license agreements.  See the NOTICE file distributed with this work for
additional information regarding copyright ownership.  The ASF licenses this
file to you under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License.  You may obtain a copy of
the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  See the
License for the specific language governing permissions and limitations under
the License.
