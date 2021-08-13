# VarMed [![Go](https://github.com/tikz/VarMed/actions/workflows/go.yml/badge.svg)](https://github.com/tikz/VarMed/actions/workflows/go.yml) [![Webpack](https://github.com/tikz/VarMed/actions/workflows/webpack.yml/badge.svg)](https://github.com/tikz/VarMed/actions/workflows/webpack.yml) [![Docker](https://github.com/tikz/VarMed/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/tikz/VarMed/actions/workflows/docker-publish.yml)
![Logo](http://varmed.qb.fcen.uba.ar/assets/varmed.svg)

http://varmed.qb.fcen.uba.ar

## Build

Requires [go](https://golang.org/doc/install#install) >= 1.14 and [yarn](https://classic.yarnpkg.com/en/docs/install/) >= 1.22 in `PATH`

Requires **propietary binaries not included** (FoldX, abSwitch, Tango) in `bin/` for some of the pipeline steps.

```
git clone https://github.com/tikz/VarMed.git
cd varmed
make build
cp config-example.yaml config.yaml
./varmed
```
