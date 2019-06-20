::protoc.exe --go_out=plugins=grpc:%outdir% --proto_path "proto" proto\*.proto
protoc --proto_path "proto" --go_out=../../pb proto/agent_path.proto
pause