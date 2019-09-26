module f/server

go 1.12

require (
	f/proto v0.0.0
	github.com/Kagami/go-face v0.0.0-20190831182441-fab496201e78 // indirect
	google.golang.org/grpc v1.23.0
)

replace f/proto => ../proto
