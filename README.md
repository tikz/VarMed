## Local build

Requires [go](https://golang.org/doc/install#install) >= 1.14 and [yarn](https://classic.yarnpkg.com/en/docs/install/) >= 1.22 in `PATH`

```
git clone git@gitlab.com:glyco1/respdb.git
cd respdb
make build
cp config-example.yaml config.yaml
./respdb
```

## Docker image

```
docker login registry.gitlab.com
docker pull registry.gitlab.com/glyco1/respdb:latest
docker run -p 8888:8888 --name respdb -dit registry.gitlab.com/glyco1/respdb:latest
```
