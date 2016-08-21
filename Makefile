all: protogen

protogen:
	protowrap -I $${GOPATH}/src \
		--gogo_out=plugins=grpc:$${GOPATH}/src \
		--proto_path $${GOPATH}/src \
		--print_structure \
		--only_specified_files \
		$$(pwd)/**/*.proto
	go install github.com/fuserobotics/rethinkts/proto
