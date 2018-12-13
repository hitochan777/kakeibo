PROTO_FILES := $(shell find . -name '*.proto')
GO_PROTO_FILES := $(patsubst %.proto,%.pb.go,$(PROTO_FILES))
JS_PROTO_FILES := $(patsubst %.proto,%.pb.js,$(PROTO_FILES))

all: proto

debug:
	@echo $(GO_PROTO_FILES)

.PHONY:proto
proto: go_proto js_proto 

.PHONY:go_proto
go_proto: $(GO_PROTO_FILES)

.PHONY:js_proto
js_proto: $(JS_PROTO_FILES)

	
%.pb.go:%.proto
	protoc -I=. --go_out=plugins=grpc:. ./$<

%.pb.js:%.proto
	protoc -I=. --js_out=import_style=common.js:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. ./$<
