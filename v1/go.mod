module learning

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/satori/go.uuid v1.2.0
	go.uber.org/ratelimit v0.1.0
	go.uber.org/zap v1.13.0
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.27.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
