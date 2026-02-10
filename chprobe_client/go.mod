module github.com/ricky97gr/chprobe/chprobe_client

go 1.24.1

replace github.com/ricky97gr/chprobe/chprobe_common => ../chprobe_common

require (
	github.com/ricky97gr/chprobe/chprobe_common v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.78.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	go.yaml.in/yaml/v3 v3.0.4
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)
