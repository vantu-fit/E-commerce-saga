protoc:
	rm -rf ./src/pb/*.d.ts
	rm -rf ./src/pb/*.ts
	rm -rf ./src/pb/*.js
	protoc -I=../proto ../proto/service_media.proto \
	  --js_out=import_style=commonjs,binary:./src/pb \
	  --grpc-web_out=import_style=typescript,mode=grpcweb:./src/pb

	protoc -I=../proto ../proto/google/api/*.proto \
	  --js_out=import_style=commonjs,binary:./src/pb \
	  --grpc-web_out=import_style=typescript,mode=grpcweb:./src/pb

	protoc -I=../proto ../proto/protoc-gen-openapiv2/options/*.proto \
	  --js_out=import_style=commonjs,binary:./src/pb \
	  --grpc-web_out=import_style=typescript,mode=grpcweb:./src/pb
	



