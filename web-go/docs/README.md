# documentation

swagger documentation is generated using swaggo

install swaggo:

`cd ..`
`go install github.com/swaggo/swag/cmd/swag@latest`

then generate json & yaml: 

`$HOME/g/bin/swag init`

generated documentation is checked in to git - it needs to be manually regenerated and checked in when interfaces, types and comments are changed