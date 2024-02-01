# grpcme


Simple Go [gRPC](https://grpc.io/) server/client template

## API Definition

To define the API (payload shape and the methods that will be used to communicate between client and the server), edit the sample `proto` file is located in [api/v1/sample.proto](api/v1/sample.proto). 

> For more about [Protocol Buffers](https://developers.google.com/protocol-buffers/) (Protobuf), language-neutral mechanism for serializing structured data, see [here](https://github.com/golang/protobuf).

To generate the `Go` code from that `proto` run:

```shell
make api
```

The resulting Go code will be written to [pkg/api/v1](pkg/api/v1). You can review that file but don't edit it as it will be overwritten the next time you generate APIs.

## Container Images

```shell
bin/image
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.