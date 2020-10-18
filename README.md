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
