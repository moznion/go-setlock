go-setlock
==

[![Build Status](https://travis-ci.org/moznion/go-setlock.svg?branch=master)](https://travis-ci.org/moznion/go-setlock)
[![wercker status](https://app.wercker.com/status/96120abee397cccab2b78f61a91f8051/s/master "wercker status")](https://app.wercker.com/project/bykey/96120abee397cccab2b78f61a91f8051)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg)](http://godoc.org/github.com/moznion/go-setlock)
[![GitHub release](http://img.shields.io/github/release/moznion/go-setlock.svg)](https://github.com/moznion/go-setlock/releases)

go-setlock is a go port of [setlock](http://cr.yp.to/daemontools/setlock.html) (an utility of daemontools),
and accompanying library.

Command Usage
--

```
setlock [ -nNxXvV ] file program [ arg ... ]
```

Command Features
--

- Functions of original setlock (See: [http://cr.yp.to/daemontools/setlock.html](http://cr.yp.to/daemontools/setlock.html))
- Support multiple environments
    - Linux
    - OS X
    - Windows

Command Installation
--

- Download built archive from [GitHub Releases](https://github.com/moznion/go-setlock/releases) and extract it
- Or install by `go get` command: `go get github.com/moznion/go-setlock/cmd/setlock`

Library Usage
--

go-setlock provides file based exclusive lock functions by `setlock` package.  
Please refer to the godoc: [![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg)](http://godoc.org/github.com/moznion/go-setlock)

Author
--

moznion (<moznion@gmail.com>)

Contributor
--

lestrrat

License
--

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/moznion/go-setlock/blob/master/LICENSE)

