// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: proto/v1/qrurl.proto

package qrurlv1connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v1 "github.com/emahiro/qrurl/server/gen/proto/v1"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// QrUrlServiceName is the fully-qualified name of the QrUrlService service.
	QrUrlServiceName = "qrurl.v1.QrUrlService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// QrUrlServicePostCodeProcedure is the fully-qualified name of the QrUrlService's PostCode RPC.
	QrUrlServicePostCodeProcedure = "/qrurl.v1.QrUrlService/PostCode"
)

// QrUrlServiceClient is a client for the qrurl.v1.QrUrlService service.
type QrUrlServiceClient interface {
	PostCode(context.Context, *connect_go.Request[v1.PostCodeRequest]) (*connect_go.Response[v1.PostCodeResponse], error)
}

// NewQrUrlServiceClient constructs a client for the qrurl.v1.QrUrlService service. By default, it
// uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses, and sends
// uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the connect.WithGRPC() or
// connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewQrUrlServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) QrUrlServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &qrUrlServiceClient{
		postCode: connect_go.NewClient[v1.PostCodeRequest, v1.PostCodeResponse](
			httpClient,
			baseURL+QrUrlServicePostCodeProcedure,
			opts...,
		),
	}
}

// qrUrlServiceClient implements QrUrlServiceClient.
type qrUrlServiceClient struct {
	postCode *connect_go.Client[v1.PostCodeRequest, v1.PostCodeResponse]
}

// PostCode calls qrurl.v1.QrUrlService.PostCode.
func (c *qrUrlServiceClient) PostCode(ctx context.Context, req *connect_go.Request[v1.PostCodeRequest]) (*connect_go.Response[v1.PostCodeResponse], error) {
	return c.postCode.CallUnary(ctx, req)
}

// QrUrlServiceHandler is an implementation of the qrurl.v1.QrUrlService service.
type QrUrlServiceHandler interface {
	PostCode(context.Context, *connect_go.Request[v1.PostCodeRequest]) (*connect_go.Response[v1.PostCodeResponse], error)
}

// NewQrUrlServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewQrUrlServiceHandler(svc QrUrlServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle(QrUrlServicePostCodeProcedure, connect_go.NewUnaryHandler(
		QrUrlServicePostCodeProcedure,
		svc.PostCode,
		opts...,
	))
	return "/qrurl.v1.QrUrlService/", mux
}

// UnimplementedQrUrlServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedQrUrlServiceHandler struct{}

func (UnimplementedQrUrlServiceHandler) PostCode(context.Context, *connect_go.Request[v1.PostCodeRequest]) (*connect_go.Response[v1.PostCodeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("qrurl.v1.QrUrlService.PostCode is not implemented"))
}
