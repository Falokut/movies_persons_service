// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: movies_persons_service_v1.proto

/*
Package protos is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package protos

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

var (
	filter_MoviesPersonsServiceV1_GetPersons_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_MoviesPersonsServiceV1_GetPersons_0(ctx context.Context, marshaler runtime.Marshaler, client MoviesPersonsServiceV1Client, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetMoviePersonsRequest
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_MoviesPersonsServiceV1_GetPersons_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GetPersons(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_MoviesPersonsServiceV1_GetPersons_0(ctx context.Context, marshaler runtime.Marshaler, server MoviesPersonsServiceV1Server, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq GetMoviePersonsRequest
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_MoviesPersonsServiceV1_GetPersons_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GetPersons(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterMoviesPersonsServiceV1HandlerServer registers the http handlers for service MoviesPersonsServiceV1 to "mux".
// UnaryRPC     :call MoviesPersonsServiceV1Server directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterMoviesPersonsServiceV1HandlerFromEndpoint instead.
func RegisterMoviesPersonsServiceV1HandlerServer(ctx context.Context, mux *runtime.ServeMux, server MoviesPersonsServiceV1Server) error {

	mux.Handle("GET", pattern_MoviesPersonsServiceV1_GetPersons_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/movies_persons_service.MoviesPersonsServiceV1/GetPersons", runtime.WithHTTPPathPattern("/v1/persons"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_MoviesPersonsServiceV1_GetPersons_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MoviesPersonsServiceV1_GetPersons_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterMoviesPersonsServiceV1HandlerFromEndpoint is same as RegisterMoviesPersonsServiceV1Handler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterMoviesPersonsServiceV1HandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterMoviesPersonsServiceV1Handler(ctx, mux, conn)
}

// RegisterMoviesPersonsServiceV1Handler registers the http handlers for service MoviesPersonsServiceV1 to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterMoviesPersonsServiceV1Handler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterMoviesPersonsServiceV1HandlerClient(ctx, mux, NewMoviesPersonsServiceV1Client(conn))
}

// RegisterMoviesPersonsServiceV1HandlerClient registers the http handlers for service MoviesPersonsServiceV1
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "MoviesPersonsServiceV1Client".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "MoviesPersonsServiceV1Client"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "MoviesPersonsServiceV1Client" to call the correct interceptors.
func RegisterMoviesPersonsServiceV1HandlerClient(ctx context.Context, mux *runtime.ServeMux, client MoviesPersonsServiceV1Client) error {

	mux.Handle("GET", pattern_MoviesPersonsServiceV1_GetPersons_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/movies_persons_service.MoviesPersonsServiceV1/GetPersons", runtime.WithHTTPPathPattern("/v1/persons"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_MoviesPersonsServiceV1_GetPersons_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_MoviesPersonsServiceV1_GetPersons_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_MoviesPersonsServiceV1_GetPersons_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "persons"}, ""))
)

var (
	forward_MoviesPersonsServiceV1_GetPersons_0 = runtime.ForwardResponseMessage
)
