// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/hanzezhenalex/socks5/src/agent/models"
)

// NewPostV1AuthUserCreateParams creates a new PostV1AuthUserCreateParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPostV1AuthUserCreateParams() *PostV1AuthUserCreateParams {
	return &PostV1AuthUserCreateParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPostV1AuthUserCreateParamsWithTimeout creates a new PostV1AuthUserCreateParams object
// with the ability to set a timeout on a request.
func NewPostV1AuthUserCreateParamsWithTimeout(timeout time.Duration) *PostV1AuthUserCreateParams {
	return &PostV1AuthUserCreateParams{
		timeout: timeout,
	}
}

// NewPostV1AuthUserCreateParamsWithContext creates a new PostV1AuthUserCreateParams object
// with the ability to set a context for a request.
func NewPostV1AuthUserCreateParamsWithContext(ctx context.Context) *PostV1AuthUserCreateParams {
	return &PostV1AuthUserCreateParams{
		Context: ctx,
	}
}

// NewPostV1AuthUserCreateParamsWithHTTPClient creates a new PostV1AuthUserCreateParams object
// with the ability to set a custom HTTPClient for a request.
func NewPostV1AuthUserCreateParamsWithHTTPClient(client *http.Client) *PostV1AuthUserCreateParams {
	return &PostV1AuthUserCreateParams{
		HTTPClient: client,
	}
}

/*
PostV1AuthUserCreateParams contains all the parameters to send to the API endpoint

	for the post v1 auth user create operation.

	Typically these are written to a http.Request.
*/
type PostV1AuthUserCreateParams struct {

	// Authorization.
	Authorization *string

	// User.
	User *models.User

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the post v1 auth user create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostV1AuthUserCreateParams) WithDefaults() *PostV1AuthUserCreateParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the post v1 auth user create params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PostV1AuthUserCreateParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) WithTimeout(timeout time.Duration) *PostV1AuthUserCreateParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) WithContext(ctx context.Context) *PostV1AuthUserCreateParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) WithHTTPClient(client *http.Client) *PostV1AuthUserCreateParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAuthorization adds the authorization to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) WithAuthorization(authorization *string) *PostV1AuthUserCreateParams {
	o.SetAuthorization(authorization)
	return o
}

// SetAuthorization adds the authorization to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) SetAuthorization(authorization *string) {
	o.Authorization = authorization
}

// WithUser adds the user to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) WithUser(user *models.User) *PostV1AuthUserCreateParams {
	o.SetUser(user)
	return o
}

// SetUser adds the user to the post v1 auth user create params
func (o *PostV1AuthUserCreateParams) SetUser(user *models.User) {
	o.User = user
}

// WriteToRequest writes these params to a swagger request
func (o *PostV1AuthUserCreateParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Authorization != nil {

		// header param Authorization
		if err := r.SetHeaderParam("Authorization", *o.Authorization); err != nil {
			return err
		}
	}
	if o.User != nil {
		if err := r.SetBodyParam(o.User); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
