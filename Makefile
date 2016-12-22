all: protogen

protogen:
	export CWD=$$(pwd) && \
	cd $${GOPATH}/src && \
	protowrap \
		-I $${GOPATH}/src \
		-I $${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:$${GOPATH}/src \
		--grpc-gateway_out=logtostderr=true:. \
		--proto_path $${GOPATH}/src \
		--print_structure \
		--only_specified_files \
		$${CWD}/**/*.proto
	go install -v github.com/paralin/rethinkts/metric
