```ps
$Env:SOURCE_DIR=".\"
$Env:DEST_DIR=".\"

protoc --proto_path=$Env:SOURCE_DIR --go_out=$Env:DEST_DIR --go_opt=paths=source_relative --go-grpc_out=$Env:DEST_DIR --go-grpc_opt=paths=source_relative $Env:SOURCE_DIR/keyvalue.proto
```