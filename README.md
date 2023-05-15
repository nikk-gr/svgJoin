[![CI](https://github.com/nikk-gr/svgJoin/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/nikk-gr/svgJoin/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/nikk-gr/svgJoin/branch/master/graph/badge.svg?token=YD7Y2EIBZJ)](https://codecov.io/gh/nikk-gr/svgJoin)
![GitHub](https://img.shields.io/github/license/nikk-gr/svgJoin)
# svgJoin
The go library for combining svg images

Three steps required to join pictures
1. Parse original svg `svgJoin.Parse()`
2. Join parsed pictures `svgJoin.Join()`
3. Export resulting picture `svgJoin.Draw()`

# Types
## Chunk
**svgJoin.Chunk** is the result of parsing svg string. It contains body of svg file and size data. Can be joined and exported.
## Group
**svgJoin.Group** is the result of joining pictures. Contains pictures and auxiliary data. Can be joined and exported.

# Functions and methods
## Parse
`Parse(svg string)(result Chink, err error)` prepare picture to the following join. Error returns if svg have invalid <svg…> tag. Deeper structure is not checked.

