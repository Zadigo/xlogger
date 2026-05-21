gen:
	@protoc \
		--proto_path=proto "proto/logs.proto" \
		--go_out=services/common/genproto/logs \
		--go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/logs \
		--go-grpc_opt=paths=source_relative
