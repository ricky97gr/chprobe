GinLogMode=""


SystemName="chprobe"
GoVersion=$(shell go version | awk '{print $$3, $$4}')
BuildTime=$(shell date "+%Y-%m-%d %H:%M:%S")
GitCommitID=$(shell git rev-parse HEAD)
Version="0.0.1_base"


VER=debug


ifeq ($(VER),debug)
	GinLogMode="debug"
else
	GinLogMode="release"
endif


define start-build=
	@echo $(1)
	@make $(1)
endef


LDFLAGS="\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.GoVersion=$(GoVersion)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.SystemName=$(SystemName)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.CommitID=$(GitCommitID)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.BuildTime=$(BuildTime)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.Version=$(Version)'\
 -X 'github.com/ricky97gr/chprobe/pkg/bininfo.GinLogMode=${GinLogMode}'\
"



.PHONY:
prepare:
	@echo "\n\n\n" >> /etc/hosts
	@echo "192.168.0.200 mysql.test.com" >> /etc/hosts
	@echo "192.168.0.200 redis.test.com" >> /etc/hosts
	@echo "192.168.0.200 mongo.test.com" >> /etc/hosts
	
	@echo "10.182.34.112 mysql.test.com" >> /etc/hosts
	@echo "10.182.34.112 redis.test.com" >> /etc/hosts
	@echo "10.182.34.112 mongo.test.com" >> /etc/hosts
	@echo "system prepare ready"
	@cp ./internal/conf/config.yaml /root/tmp/.config.yaml
	@echo "prepare db env"
	@bash ./scripts/prepare.sh


.PHONY:
reportergrpc:
	@protoc --go_out=./proto/ ./proto/reporter.proto
	@protoc --go-grpc_out=./proto/ ./proto/reporter.proto
	@echo "reporter grpc protobuf generate successfully"

	

.PHONY:
release:
	@echo "start to build release version"
	@rm -rf bin
	@mkdir bin
	@go build -ldflags ${LDFLAGS} -o bin/family_webservice cmd/webservice/main.go
	@echo "build family_webservice successfully!"
	

.PHONY:
debug:
	@echo "start to build debug version"
	@rm -rf bin
	@mkdir bin
	@go build -ldflags ${LDFLAGS} -o bin/family_webservice cmd/webservice/main.go
	@echo "build family_webservice successfully!"


.PHONY:
run:
ifeq ($(VER), debug)
	@$(call start-build,debug)
else
	@$(call start-build,release)
endif


