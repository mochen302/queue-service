#！/bin/bash
govendor fetch github.com/emirpasic/gods@v1.12.0
govendor fetch github.com/gin-gonic/gin@v1.4.0
govendor fetch github.com/sirupsen/logrus@v1.4.2
govendor fetch gopkg.in/yaml.v3
go mod vendor
go build src/main/server.go
