// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.11.4
// source: gateway.proto

package moby_buildkit_v1_frontend

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LLBBridge_ResolveImageConfig_FullMethodName = "/moby.buildkit.v1.frontend.LLBBridge/ResolveImageConfig"
	LLBBridge_ResolveSourceMeta_FullMethodName  = "/moby.buildkit.v1.frontend.LLBBridge/ResolveSourceMeta"
	LLBBridge_Solve_FullMethodName              = "/moby.buildkit.v1.frontend.LLBBridge/Solve"
	LLBBridge_ReadFile_FullMethodName           = "/moby.buildkit.v1.frontend.LLBBridge/ReadFile"
	LLBBridge_ReadDir_FullMethodName            = "/moby.buildkit.v1.frontend.LLBBridge/ReadDir"
	LLBBridge_StatFile_FullMethodName           = "/moby.buildkit.v1.frontend.LLBBridge/StatFile"
	LLBBridge_Evaluate_FullMethodName           = "/moby.buildkit.v1.frontend.LLBBridge/Evaluate"
	LLBBridge_Ping_FullMethodName               = "/moby.buildkit.v1.frontend.LLBBridge/Ping"
	LLBBridge_Return_FullMethodName             = "/moby.buildkit.v1.frontend.LLBBridge/Return"
	LLBBridge_Inputs_FullMethodName             = "/moby.buildkit.v1.frontend.LLBBridge/Inputs"
	LLBBridge_NewContainer_FullMethodName       = "/moby.buildkit.v1.frontend.LLBBridge/NewContainer"
	LLBBridge_ReleaseContainer_FullMethodName   = "/moby.buildkit.v1.frontend.LLBBridge/ReleaseContainer"
	LLBBridge_ExecProcess_FullMethodName        = "/moby.buildkit.v1.frontend.LLBBridge/ExecProcess"
	LLBBridge_Warn_FullMethodName               = "/moby.buildkit.v1.frontend.LLBBridge/Warn"
)

// LLBBridgeClient is the client API for LLBBridge service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LLBBridgeClient interface {
	// apicaps:CapResolveImage
	ResolveImageConfig(ctx context.Context, in *ResolveImageConfigRequest, opts ...grpc.CallOption) (*ResolveImageConfigResponse, error)
	// apicaps:CapSourceMetaResolver
	ResolveSourceMeta(ctx context.Context, in *ResolveSourceMetaRequest, opts ...grpc.CallOption) (*ResolveSourceMetaResponse, error)
	// apicaps:CapSolveBase
	Solve(ctx context.Context, in *SolveRequest, opts ...grpc.CallOption) (*SolveResponse, error)
	// apicaps:CapReadFile
	ReadFile(ctx context.Context, in *ReadFileRequest, opts ...grpc.CallOption) (*ReadFileResponse, error)
	// apicaps:CapReadDir
	ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (*ReadDirResponse, error)
	// apicaps:CapStatFile
	StatFile(ctx context.Context, in *StatFileRequest, opts ...grpc.CallOption) (*StatFileResponse, error)
	// apicaps:CapGatewayEvaluate
	Evaluate(ctx context.Context, in *EvaluateRequest, opts ...grpc.CallOption) (*EvaluateResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error)
	Return(ctx context.Context, in *ReturnRequest, opts ...grpc.CallOption) (*ReturnResponse, error)
	// apicaps:CapFrontendInputs
	Inputs(ctx context.Context, in *InputsRequest, opts ...grpc.CallOption) (*InputsResponse, error)
	NewContainer(ctx context.Context, in *NewContainerRequest, opts ...grpc.CallOption) (*NewContainerResponse, error)
	ReleaseContainer(ctx context.Context, in *ReleaseContainerRequest, opts ...grpc.CallOption) (*ReleaseContainerResponse, error)
	ExecProcess(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ExecMessage, ExecMessage], error)
	// apicaps:CapGatewayWarnings
	Warn(ctx context.Context, in *WarnRequest, opts ...grpc.CallOption) (*WarnResponse, error)
}

type lLBBridgeClient struct {
	cc grpc.ClientConnInterface
}

func NewLLBBridgeClient(cc grpc.ClientConnInterface) LLBBridgeClient {
	return &lLBBridgeClient{cc}
}

func (c *lLBBridgeClient) ResolveImageConfig(ctx context.Context, in *ResolveImageConfigRequest, opts ...grpc.CallOption) (*ResolveImageConfigResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResolveImageConfigResponse)
	err := c.cc.Invoke(ctx, LLBBridge_ResolveImageConfig_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) ResolveSourceMeta(ctx context.Context, in *ResolveSourceMetaRequest, opts ...grpc.CallOption) (*ResolveSourceMetaResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResolveSourceMetaResponse)
	err := c.cc.Invoke(ctx, LLBBridge_ResolveSourceMeta_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) Solve(ctx context.Context, in *SolveRequest, opts ...grpc.CallOption) (*SolveResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SolveResponse)
	err := c.cc.Invoke(ctx, LLBBridge_Solve_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) ReadFile(ctx context.Context, in *ReadFileRequest, opts ...grpc.CallOption) (*ReadFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadFileResponse)
	err := c.cc.Invoke(ctx, LLBBridge_ReadFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) ReadDir(ctx context.Context, in *ReadDirRequest, opts ...grpc.CallOption) (*ReadDirResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadDirResponse)
	err := c.cc.Invoke(ctx, LLBBridge_ReadDir_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) StatFile(ctx context.Context, in *StatFileRequest, opts ...grpc.CallOption) (*StatFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(StatFileResponse)
	err := c.cc.Invoke(ctx, LLBBridge_StatFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) Evaluate(ctx context.Context, in *EvaluateRequest, opts ...grpc.CallOption) (*EvaluateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EvaluateResponse)
	err := c.cc.Invoke(ctx, LLBBridge_Evaluate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PongResponse)
	err := c.cc.Invoke(ctx, LLBBridge_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) Return(ctx context.Context, in *ReturnRequest, opts ...grpc.CallOption) (*ReturnResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReturnResponse)
	err := c.cc.Invoke(ctx, LLBBridge_Return_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) Inputs(ctx context.Context, in *InputsRequest, opts ...grpc.CallOption) (*InputsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InputsResponse)
	err := c.cc.Invoke(ctx, LLBBridge_Inputs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) NewContainer(ctx context.Context, in *NewContainerRequest, opts ...grpc.CallOption) (*NewContainerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NewContainerResponse)
	err := c.cc.Invoke(ctx, LLBBridge_NewContainer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) ReleaseContainer(ctx context.Context, in *ReleaseContainerRequest, opts ...grpc.CallOption) (*ReleaseContainerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReleaseContainerResponse)
	err := c.cc.Invoke(ctx, LLBBridge_ReleaseContainer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lLBBridgeClient) ExecProcess(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ExecMessage, ExecMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &LLBBridge_ServiceDesc.Streams[0], LLBBridge_ExecProcess_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ExecMessage, ExecMessage]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LLBBridge_ExecProcessClient = grpc.BidiStreamingClient[ExecMessage, ExecMessage]

func (c *lLBBridgeClient) Warn(ctx context.Context, in *WarnRequest, opts ...grpc.CallOption) (*WarnResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WarnResponse)
	err := c.cc.Invoke(ctx, LLBBridge_Warn_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LLBBridgeServer is the server API for LLBBridge service.
// All implementations should embed UnimplementedLLBBridgeServer
// for forward compatibility.
type LLBBridgeServer interface {
	// apicaps:CapResolveImage
	ResolveImageConfig(context.Context, *ResolveImageConfigRequest) (*ResolveImageConfigResponse, error)
	// apicaps:CapSourceMetaResolver
	ResolveSourceMeta(context.Context, *ResolveSourceMetaRequest) (*ResolveSourceMetaResponse, error)
	// apicaps:CapSolveBase
	Solve(context.Context, *SolveRequest) (*SolveResponse, error)
	// apicaps:CapReadFile
	ReadFile(context.Context, *ReadFileRequest) (*ReadFileResponse, error)
	// apicaps:CapReadDir
	ReadDir(context.Context, *ReadDirRequest) (*ReadDirResponse, error)
	// apicaps:CapStatFile
	StatFile(context.Context, *StatFileRequest) (*StatFileResponse, error)
	// apicaps:CapGatewayEvaluate
	Evaluate(context.Context, *EvaluateRequest) (*EvaluateResponse, error)
	Ping(context.Context, *PingRequest) (*PongResponse, error)
	Return(context.Context, *ReturnRequest) (*ReturnResponse, error)
	// apicaps:CapFrontendInputs
	Inputs(context.Context, *InputsRequest) (*InputsResponse, error)
	NewContainer(context.Context, *NewContainerRequest) (*NewContainerResponse, error)
	ReleaseContainer(context.Context, *ReleaseContainerRequest) (*ReleaseContainerResponse, error)
	ExecProcess(grpc.BidiStreamingServer[ExecMessage, ExecMessage]) error
	// apicaps:CapGatewayWarnings
	Warn(context.Context, *WarnRequest) (*WarnResponse, error)
}

// UnimplementedLLBBridgeServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLLBBridgeServer struct{}

func (UnimplementedLLBBridgeServer) ResolveImageConfig(context.Context, *ResolveImageConfigRequest) (*ResolveImageConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveImageConfig not implemented")
}
func (UnimplementedLLBBridgeServer) ResolveSourceMeta(context.Context, *ResolveSourceMetaRequest) (*ResolveSourceMetaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveSourceMeta not implemented")
}
func (UnimplementedLLBBridgeServer) Solve(context.Context, *SolveRequest) (*SolveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Solve not implemented")
}
func (UnimplementedLLBBridgeServer) ReadFile(context.Context, *ReadFileRequest) (*ReadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadFile not implemented")
}
func (UnimplementedLLBBridgeServer) ReadDir(context.Context, *ReadDirRequest) (*ReadDirResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadDir not implemented")
}
func (UnimplementedLLBBridgeServer) StatFile(context.Context, *StatFileRequest) (*StatFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StatFile not implemented")
}
func (UnimplementedLLBBridgeServer) Evaluate(context.Context, *EvaluateRequest) (*EvaluateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Evaluate not implemented")
}
func (UnimplementedLLBBridgeServer) Ping(context.Context, *PingRequest) (*PongResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedLLBBridgeServer) Return(context.Context, *ReturnRequest) (*ReturnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Return not implemented")
}
func (UnimplementedLLBBridgeServer) Inputs(context.Context, *InputsRequest) (*InputsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Inputs not implemented")
}
func (UnimplementedLLBBridgeServer) NewContainer(context.Context, *NewContainerRequest) (*NewContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewContainer not implemented")
}
func (UnimplementedLLBBridgeServer) ReleaseContainer(context.Context, *ReleaseContainerRequest) (*ReleaseContainerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseContainer not implemented")
}
func (UnimplementedLLBBridgeServer) ExecProcess(grpc.BidiStreamingServer[ExecMessage, ExecMessage]) error {
	return status.Errorf(codes.Unimplemented, "method ExecProcess not implemented")
}
func (UnimplementedLLBBridgeServer) Warn(context.Context, *WarnRequest) (*WarnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Warn not implemented")
}
func (UnimplementedLLBBridgeServer) testEmbeddedByValue() {}

// UnsafeLLBBridgeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LLBBridgeServer will
// result in compilation errors.
type UnsafeLLBBridgeServer interface {
	mustEmbedUnimplementedLLBBridgeServer()
}

func RegisterLLBBridgeServer(s grpc.ServiceRegistrar, srv LLBBridgeServer) {
	// If the following call pancis, it indicates UnimplementedLLBBridgeServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LLBBridge_ServiceDesc, srv)
}

func _LLBBridge_ResolveImageConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveImageConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).ResolveImageConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_ResolveImageConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).ResolveImageConfig(ctx, req.(*ResolveImageConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_ResolveSourceMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveSourceMetaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).ResolveSourceMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_ResolveSourceMeta_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).ResolveSourceMeta(ctx, req.(*ResolveSourceMetaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_Solve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SolveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).Solve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_Solve_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).Solve(ctx, req.(*SolveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_ReadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).ReadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_ReadFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).ReadFile(ctx, req.(*ReadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_ReadDir_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadDirRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).ReadDir(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_ReadDir_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).ReadDir(ctx, req.(*ReadDirRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_StatFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).StatFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_StatFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).StatFile(ctx, req.(*StatFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_Evaluate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EvaluateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).Evaluate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_Evaluate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).Evaluate(ctx, req.(*EvaluateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_Return_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReturnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).Return(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_Return_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).Return(ctx, req.(*ReturnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_Inputs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InputsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).Inputs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_Inputs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).Inputs(ctx, req.(*InputsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_NewContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).NewContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_NewContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).NewContainer(ctx, req.(*NewContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_ReleaseContainer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseContainerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).ReleaseContainer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_ReleaseContainer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).ReleaseContainer(ctx, req.(*ReleaseContainerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LLBBridge_ExecProcess_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(LLBBridgeServer).ExecProcess(&grpc.GenericServerStream[ExecMessage, ExecMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type LLBBridge_ExecProcessServer = grpc.BidiStreamingServer[ExecMessage, ExecMessage]

func _LLBBridge_Warn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WarnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LLBBridgeServer).Warn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LLBBridge_Warn_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LLBBridgeServer).Warn(ctx, req.(*WarnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LLBBridge_ServiceDesc is the grpc.ServiceDesc for LLBBridge service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LLBBridge_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "moby.buildkit.v1.frontend.LLBBridge",
	HandlerType: (*LLBBridgeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResolveImageConfig",
			Handler:    _LLBBridge_ResolveImageConfig_Handler,
		},
		{
			MethodName: "ResolveSourceMeta",
			Handler:    _LLBBridge_ResolveSourceMeta_Handler,
		},
		{
			MethodName: "Solve",
			Handler:    _LLBBridge_Solve_Handler,
		},
		{
			MethodName: "ReadFile",
			Handler:    _LLBBridge_ReadFile_Handler,
		},
		{
			MethodName: "ReadDir",
			Handler:    _LLBBridge_ReadDir_Handler,
		},
		{
			MethodName: "StatFile",
			Handler:    _LLBBridge_StatFile_Handler,
		},
		{
			MethodName: "Evaluate",
			Handler:    _LLBBridge_Evaluate_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _LLBBridge_Ping_Handler,
		},
		{
			MethodName: "Return",
			Handler:    _LLBBridge_Return_Handler,
		},
		{
			MethodName: "Inputs",
			Handler:    _LLBBridge_Inputs_Handler,
		},
		{
			MethodName: "NewContainer",
			Handler:    _LLBBridge_NewContainer_Handler,
		},
		{
			MethodName: "ReleaseContainer",
			Handler:    _LLBBridge_ReleaseContainer_Handler,
		},
		{
			MethodName: "Warn",
			Handler:    _LLBBridge_Warn_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExecProcess",
			Handler:       _LLBBridge_ExecProcess_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "gateway.proto",
}
