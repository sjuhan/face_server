module f/client

go 1.12

require (
	f/proto v0.0.0
	google.golang.org/grpc v1.23.0
)

replace f/proto => ../proto
