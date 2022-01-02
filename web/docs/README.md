# documentation

swagger documentation is generated using swaggo

install swaggo:

`go install github.com/swaggo/swag/cmd/swag@latest`

then generate json & yaml: 

`cd ..`
`swag init -g main.go -d .,./handlers/ -o ./docs`

generated documentation is checked in to git - it needs to be manually regenerated and checked in when interfaces are changed