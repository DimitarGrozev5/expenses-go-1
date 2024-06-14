// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.0
// source: models.proto

package models

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DatabaseClient is the client API for Database service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseClient interface {
	// Simple ping method
	Ping(ctx context.Context, in *SimpleMessage, opts ...grpc.CallOption) (*SimpleMessage, error)
	// User
	GetUser(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GrpcUser, error)
	Authenticate(ctx context.Context, in *LoginCredentials, opts ...grpc.CallOption) (*LoginToken, error)
	Logout(ctx context.Context, in *LogoutParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	ModifyFreeFunds(ctx context.Context, in *ModifyFreeFundsParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	// Tags methods
	GetTags(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetTagsReturns, error)
	// Expenses methods
	GetExpenses(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetExpensesReturns, error)
	AddExpense(ctx context.Context, in *ExpensesParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	EditExpense(ctx context.Context, in *ExpensesParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	DeleteExpense(ctx context.Context, in *DeleteExpenseParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	// Accounts methods
	GetAccounts(ctx context.Context, in *GetAccountsParams, opts ...grpc.CallOption) (*GetAccountsReturns, error)
	AddAccount(ctx context.Context, in *AddAccountParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	EditAccountName(ctx context.Context, in *EditAccountNameParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	DeleteAccount(ctx context.Context, in *DeleteAccountParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	TransferFunds(ctx context.Context, in *TransferFundsParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	ReorderAccount(ctx context.Context, in *ReorderAccountParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	// Categories methods
	GetCategoriesCount(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetCategoriesCountReturns, error)
	GetCategories(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetCategoriesReturns, error)
	GetCategoriesOverview(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetCategoriesOverviewReturns, error)
	AddCategory(ctx context.Context, in *AddCategoryParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	ReorderCategory(ctx context.Context, in *ReorderCategoryParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	DeleteCategory(ctx context.Context, in *DeleteCategoryParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	ResetCategories(ctx context.Context, in *ResetCategoriesParams, opts ...grpc.CallOption) (*GrpcEmpty, error)
	// Time periods
	GetTimePeriods(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetTimePeriodsReturns, error)
}

type databaseClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseClient(cc grpc.ClientConnInterface) DatabaseClient {
	return &databaseClient{cc}
}

func (c *databaseClient) Ping(ctx context.Context, in *SimpleMessage, opts ...grpc.CallOption) (*SimpleMessage, error) {
	out := new(SimpleMessage)
	err := c.cc.Invoke(ctx, "/Database/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetUser(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GrpcUser, error) {
	out := new(GrpcUser)
	err := c.cc.Invoke(ctx, "/Database/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Authenticate(ctx context.Context, in *LoginCredentials, opts ...grpc.CallOption) (*LoginToken, error) {
	out := new(LoginToken)
	err := c.cc.Invoke(ctx, "/Database/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Logout(ctx context.Context, in *LogoutParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ModifyFreeFunds(ctx context.Context, in *ModifyFreeFundsParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/ModifyFreeFunds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetTags(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetTagsReturns, error) {
	out := new(GetTagsReturns)
	err := c.cc.Invoke(ctx, "/Database/GetTags", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetExpenses(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetExpensesReturns, error) {
	out := new(GetExpensesReturns)
	err := c.cc.Invoke(ctx, "/Database/GetExpenses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) AddExpense(ctx context.Context, in *ExpensesParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/AddExpense", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) EditExpense(ctx context.Context, in *ExpensesParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/EditExpense", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteExpense(ctx context.Context, in *DeleteExpenseParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/DeleteExpense", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetAccounts(ctx context.Context, in *GetAccountsParams, opts ...grpc.CallOption) (*GetAccountsReturns, error) {
	out := new(GetAccountsReturns)
	err := c.cc.Invoke(ctx, "/Database/GetAccounts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) AddAccount(ctx context.Context, in *AddAccountParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/AddAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) EditAccountName(ctx context.Context, in *EditAccountNameParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/EditAccountName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteAccount(ctx context.Context, in *DeleteAccountParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/DeleteAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) TransferFunds(ctx context.Context, in *TransferFundsParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/TransferFunds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ReorderAccount(ctx context.Context, in *ReorderAccountParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/ReorderAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetCategoriesCount(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetCategoriesCountReturns, error) {
	out := new(GetCategoriesCountReturns)
	err := c.cc.Invoke(ctx, "/Database/GetCategoriesCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetCategories(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetCategoriesReturns, error) {
	out := new(GetCategoriesReturns)
	err := c.cc.Invoke(ctx, "/Database/GetCategories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetCategoriesOverview(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetCategoriesOverviewReturns, error) {
	out := new(GetCategoriesOverviewReturns)
	err := c.cc.Invoke(ctx, "/Database/GetCategoriesOverview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) AddCategory(ctx context.Context, in *AddCategoryParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/AddCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ReorderCategory(ctx context.Context, in *ReorderCategoryParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/ReorderCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteCategory(ctx context.Context, in *DeleteCategoryParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/DeleteCategory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ResetCategories(ctx context.Context, in *ResetCategoriesParams, opts ...grpc.CallOption) (*GrpcEmpty, error) {
	out := new(GrpcEmpty)
	err := c.cc.Invoke(ctx, "/Database/ResetCategories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetTimePeriods(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GetTimePeriodsReturns, error) {
	out := new(GetTimePeriodsReturns)
	err := c.cc.Invoke(ctx, "/Database/GetTimePeriods", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseServer is the server API for Database service.
// All implementations must embed UnimplementedDatabaseServer
// for forward compatibility
type DatabaseServer interface {
	// Simple ping method
	Ping(context.Context, *SimpleMessage) (*SimpleMessage, error)
	// User
	GetUser(context.Context, *GrpcEmpty) (*GrpcUser, error)
	Authenticate(context.Context, *LoginCredentials) (*LoginToken, error)
	Logout(context.Context, *LogoutParams) (*GrpcEmpty, error)
	ModifyFreeFunds(context.Context, *ModifyFreeFundsParams) (*GrpcEmpty, error)
	// Tags methods
	GetTags(context.Context, *GrpcEmpty) (*GetTagsReturns, error)
	// Expenses methods
	GetExpenses(context.Context, *GrpcEmpty) (*GetExpensesReturns, error)
	AddExpense(context.Context, *ExpensesParams) (*GrpcEmpty, error)
	EditExpense(context.Context, *ExpensesParams) (*GrpcEmpty, error)
	DeleteExpense(context.Context, *DeleteExpenseParams) (*GrpcEmpty, error)
	// Accounts methods
	GetAccounts(context.Context, *GetAccountsParams) (*GetAccountsReturns, error)
	AddAccount(context.Context, *AddAccountParams) (*GrpcEmpty, error)
	EditAccountName(context.Context, *EditAccountNameParams) (*GrpcEmpty, error)
	DeleteAccount(context.Context, *DeleteAccountParams) (*GrpcEmpty, error)
	TransferFunds(context.Context, *TransferFundsParams) (*GrpcEmpty, error)
	ReorderAccount(context.Context, *ReorderAccountParams) (*GrpcEmpty, error)
	// Categories methods
	GetCategoriesCount(context.Context, *GrpcEmpty) (*GetCategoriesCountReturns, error)
	GetCategories(context.Context, *GrpcEmpty) (*GetCategoriesReturns, error)
	GetCategoriesOverview(context.Context, *GrpcEmpty) (*GetCategoriesOverviewReturns, error)
	AddCategory(context.Context, *AddCategoryParams) (*GrpcEmpty, error)
	ReorderCategory(context.Context, *ReorderCategoryParams) (*GrpcEmpty, error)
	DeleteCategory(context.Context, *DeleteCategoryParams) (*GrpcEmpty, error)
	ResetCategories(context.Context, *ResetCategoriesParams) (*GrpcEmpty, error)
	// Time periods
	GetTimePeriods(context.Context, *GrpcEmpty) (*GetTimePeriodsReturns, error)
	mustEmbedUnimplementedDatabaseServer()
}

// UnimplementedDatabaseServer must be embedded to have forward compatible implementations.
type UnimplementedDatabaseServer struct {
}

func (UnimplementedDatabaseServer) Ping(context.Context, *SimpleMessage) (*SimpleMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDatabaseServer) GetUser(context.Context, *GrpcEmpty) (*GrpcUser, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedDatabaseServer) Authenticate(context.Context, *LoginCredentials) (*LoginToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (UnimplementedDatabaseServer) Logout(context.Context, *LogoutParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedDatabaseServer) ModifyFreeFunds(context.Context, *ModifyFreeFundsParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyFreeFunds not implemented")
}
func (UnimplementedDatabaseServer) GetTags(context.Context, *GrpcEmpty) (*GetTagsReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTags not implemented")
}
func (UnimplementedDatabaseServer) GetExpenses(context.Context, *GrpcEmpty) (*GetExpensesReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetExpenses not implemented")
}
func (UnimplementedDatabaseServer) AddExpense(context.Context, *ExpensesParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddExpense not implemented")
}
func (UnimplementedDatabaseServer) EditExpense(context.Context, *ExpensesParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditExpense not implemented")
}
func (UnimplementedDatabaseServer) DeleteExpense(context.Context, *DeleteExpenseParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteExpense not implemented")
}
func (UnimplementedDatabaseServer) GetAccounts(context.Context, *GetAccountsParams) (*GetAccountsReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccounts not implemented")
}
func (UnimplementedDatabaseServer) AddAccount(context.Context, *AddAccountParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAccount not implemented")
}
func (UnimplementedDatabaseServer) EditAccountName(context.Context, *EditAccountNameParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditAccountName not implemented")
}
func (UnimplementedDatabaseServer) DeleteAccount(context.Context, *DeleteAccountParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccount not implemented")
}
func (UnimplementedDatabaseServer) TransferFunds(context.Context, *TransferFundsParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferFunds not implemented")
}
func (UnimplementedDatabaseServer) ReorderAccount(context.Context, *ReorderAccountParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReorderAccount not implemented")
}
func (UnimplementedDatabaseServer) GetCategoriesCount(context.Context, *GrpcEmpty) (*GetCategoriesCountReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategoriesCount not implemented")
}
func (UnimplementedDatabaseServer) GetCategories(context.Context, *GrpcEmpty) (*GetCategoriesReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategories not implemented")
}
func (UnimplementedDatabaseServer) GetCategoriesOverview(context.Context, *GrpcEmpty) (*GetCategoriesOverviewReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCategoriesOverview not implemented")
}
func (UnimplementedDatabaseServer) AddCategory(context.Context, *AddCategoryParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddCategory not implemented")
}
func (UnimplementedDatabaseServer) ReorderCategory(context.Context, *ReorderCategoryParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReorderCategory not implemented")
}
func (UnimplementedDatabaseServer) DeleteCategory(context.Context, *DeleteCategoryParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCategory not implemented")
}
func (UnimplementedDatabaseServer) ResetCategories(context.Context, *ResetCategoriesParams) (*GrpcEmpty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetCategories not implemented")
}
func (UnimplementedDatabaseServer) GetTimePeriods(context.Context, *GrpcEmpty) (*GetTimePeriodsReturns, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTimePeriods not implemented")
}
func (UnimplementedDatabaseServer) mustEmbedUnimplementedDatabaseServer() {}

// UnsafeDatabaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseServer will
// result in compilation errors.
type UnsafeDatabaseServer interface {
	mustEmbedUnimplementedDatabaseServer()
}

func RegisterDatabaseServer(s grpc.ServiceRegistrar, srv DatabaseServer) {
	s.RegisterService(&Database_ServiceDesc, srv)
}

func _Database_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SimpleMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Ping(ctx, req.(*SimpleMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetUser(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginCredentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Authenticate(ctx, req.(*LoginCredentials))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Logout(ctx, req.(*LogoutParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ModifyFreeFunds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyFreeFundsParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).ModifyFreeFunds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/ModifyFreeFunds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).ModifyFreeFunds(ctx, req.(*ModifyFreeFundsParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetTags",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetTags(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetExpenses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetExpenses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetExpenses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetExpenses(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_AddExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpensesParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).AddExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/AddExpense",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).AddExpense(ctx, req.(*ExpensesParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_EditExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExpensesParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).EditExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/EditExpense",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).EditExpense(ctx, req.(*ExpensesParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteExpense_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteExpenseParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).DeleteExpense(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/DeleteExpense",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).DeleteExpense(ctx, req.(*DeleteExpenseParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetAccounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountsParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetAccounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetAccounts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetAccounts(ctx, req.(*GetAccountsParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_AddAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAccountParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).AddAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/AddAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).AddAccount(ctx, req.(*AddAccountParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_EditAccountName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditAccountNameParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).EditAccountName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/EditAccountName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).EditAccountName(ctx, req.(*EditAccountNameParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccountParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).DeleteAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/DeleteAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).DeleteAccount(ctx, req.(*DeleteAccountParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_TransferFunds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferFundsParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).TransferFunds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/TransferFunds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).TransferFunds(ctx, req.(*TransferFundsParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ReorderAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReorderAccountParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).ReorderAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/ReorderAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).ReorderAccount(ctx, req.(*ReorderAccountParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetCategoriesCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetCategoriesCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetCategoriesCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetCategoriesCount(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetCategories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetCategories(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetCategoriesOverview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetCategoriesOverview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetCategoriesOverview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetCategoriesOverview(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_AddCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddCategoryParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).AddCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/AddCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).AddCategory(ctx, req.(*AddCategoryParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ReorderCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReorderCategoryParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).ReorderCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/ReorderCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).ReorderCategory(ctx, req.(*ReorderCategoryParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteCategory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCategoryParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).DeleteCategory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/DeleteCategory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).DeleteCategory(ctx, req.(*DeleteCategoryParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ResetCategories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetCategoriesParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).ResetCategories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/ResetCategories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).ResetCategories(ctx, req.(*ResetCategoriesParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetTimePeriods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetTimePeriods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Database/GetTimePeriods",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetTimePeriods(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

// Database_ServiceDesc is the grpc.ServiceDesc for Database service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Database_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Database",
	HandlerType: (*DatabaseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Database_Ping_Handler,
		},
		{
			MethodName: "GetUser",
			Handler:    _Database_GetUser_Handler,
		},
		{
			MethodName: "Authenticate",
			Handler:    _Database_Authenticate_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Database_Logout_Handler,
		},
		{
			MethodName: "ModifyFreeFunds",
			Handler:    _Database_ModifyFreeFunds_Handler,
		},
		{
			MethodName: "GetTags",
			Handler:    _Database_GetTags_Handler,
		},
		{
			MethodName: "GetExpenses",
			Handler:    _Database_GetExpenses_Handler,
		},
		{
			MethodName: "AddExpense",
			Handler:    _Database_AddExpense_Handler,
		},
		{
			MethodName: "EditExpense",
			Handler:    _Database_EditExpense_Handler,
		},
		{
			MethodName: "DeleteExpense",
			Handler:    _Database_DeleteExpense_Handler,
		},
		{
			MethodName: "GetAccounts",
			Handler:    _Database_GetAccounts_Handler,
		},
		{
			MethodName: "AddAccount",
			Handler:    _Database_AddAccount_Handler,
		},
		{
			MethodName: "EditAccountName",
			Handler:    _Database_EditAccountName_Handler,
		},
		{
			MethodName: "DeleteAccount",
			Handler:    _Database_DeleteAccount_Handler,
		},
		{
			MethodName: "TransferFunds",
			Handler:    _Database_TransferFunds_Handler,
		},
		{
			MethodName: "ReorderAccount",
			Handler:    _Database_ReorderAccount_Handler,
		},
		{
			MethodName: "GetCategoriesCount",
			Handler:    _Database_GetCategoriesCount_Handler,
		},
		{
			MethodName: "GetCategories",
			Handler:    _Database_GetCategories_Handler,
		},
		{
			MethodName: "GetCategoriesOverview",
			Handler:    _Database_GetCategoriesOverview_Handler,
		},
		{
			MethodName: "AddCategory",
			Handler:    _Database_AddCategory_Handler,
		},
		{
			MethodName: "ReorderCategory",
			Handler:    _Database_ReorderCategory_Handler,
		},
		{
			MethodName: "DeleteCategory",
			Handler:    _Database_DeleteCategory_Handler,
		},
		{
			MethodName: "ResetCategories",
			Handler:    _Database_ResetCategories_Handler,
		},
		{
			MethodName: "GetTimePeriods",
			Handler:    _Database_GetTimePeriods_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "models.proto",
}
