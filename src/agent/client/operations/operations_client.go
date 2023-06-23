// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientOption is the option for Client methods
type ClientOption func(*runtime.ClientOperation)

// ClientService is the interface for Client methods
type ClientService interface {
	GetV1ConnectionList(params *GetV1ConnectionListParams, opts ...ClientOption) (*GetV1ConnectionListOK, error)

	PostLogin(params *PostLoginParams, opts ...ClientOption) (*PostLoginOK, error)

	PostV1AuthUserCreate(params *PostV1AuthUserCreateParams, opts ...ClientOption) (*PostV1AuthUserCreateOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetV1ConnectionList Get the proxy connections.
*/
func (a *Client) GetV1ConnectionList(params *GetV1ConnectionListParams, opts ...ClientOption) (*GetV1ConnectionListOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetV1ConnectionListParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "GetV1ConnectionList",
		Method:             "GET",
		PathPattern:        "/v1/connection/list",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetV1ConnectionListReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetV1ConnectionListOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for GetV1ConnectionList: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  PostLogin login in the socks system.
*/
func (a *Client) PostLogin(params *PostLoginParams, opts ...ClientOption) (*PostLoginOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostLoginParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostLogin",
		Method:             "POST",
		PathPattern:        "/login",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostLoginReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostLoginOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostLogin: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

/*
  PostV1AuthUserCreate create new user.
*/
func (a *Client) PostV1AuthUserCreate(params *PostV1AuthUserCreateParams, opts ...ClientOption) (*PostV1AuthUserCreateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostV1AuthUserCreateParams()
	}
	op := &runtime.ClientOperation{
		ID:                 "PostV1AuthUserCreate",
		Method:             "POST",
		PathPattern:        "/v1/auth/user/create",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostV1AuthUserCreateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	}
	for _, opt := range opts {
		opt(op)
	}

	result, err := a.transport.Submit(op)
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostV1AuthUserCreateOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	// safeguard: normally, absent a default response, unknown success responses return an error above: so this is a codegen issue
	msg := fmt.Sprintf("unexpected success response for PostV1AuthUserCreate: API contract not enforced by server. Client expected to get an error, but got: %T", result)
	panic(msg)
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
