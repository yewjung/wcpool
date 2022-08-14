module party

go 1.19

replace (
	wcpool/authorization => ../authorization
	wcpool/constants => ../constants
	wcpool/driver => ../driver
	wcpool/models => ../models
	wcpool/utils => ../utils
)

require (
	github.com/go-redis/redis/v9 v9.0.0-beta.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.6
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.10.1
	google.golang.org/grpc v1.48.0
	wcpool/authorization v0.0.0-00010101000000-000000000000
	wcpool/constants v0.0.0-00010101000000-000000000000
	wcpool/driver v0.0.0-00010101000000-000000000000
	wcpool/models v0.0.0-00010101000000-000000000000
	wcpool/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220422013727-9388b58f7150 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
