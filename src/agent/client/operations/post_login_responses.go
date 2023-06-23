// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// PostLoginReader is a Reader for the PostLogin structure.
type PostLoginReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostLoginReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostLoginOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPostLoginBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewPostLoginInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[POST /login] PostLogin", response, response.Code())
	}
}

// NewPostLoginOK creates a PostLoginOK with default headers values
func NewPostLoginOK() *PostLoginOK {
	return &PostLoginOK{}
}

/*
PostLoginOK describes a response with status code 200, with default header values.

success.
*/
type PostLoginOK struct {
}

// IsSuccess returns true when this post login o k response has a 2xx status code
func (o *PostLoginOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this post login o k response has a 3xx status code
func (o *PostLoginOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post login o k response has a 4xx status code
func (o *PostLoginOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this post login o k response has a 5xx status code
func (o *PostLoginOK) IsServerError() bool {
	return false
}

// IsCode returns true when this post login o k response a status code equal to that given
func (o *PostLoginOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the post login o k response
func (o *PostLoginOK) Code() int {
	return 200
}

func (o *PostLoginOK) Error() string {
	return fmt.Sprintf("[POST /login][%d] postLoginOK ", 200)
}

func (o *PostLoginOK) String() string {
	return fmt.Sprintf("[POST /login][%d] postLoginOK ", 200)
}

func (o *PostLoginOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewPostLoginBadRequest creates a PostLoginBadRequest with default headers values
func NewPostLoginBadRequest() *PostLoginBadRequest {
	return &PostLoginBadRequest{}
}

/*
PostLoginBadRequest describes a response with status code 400, with default header values.

illegal request.
*/
type PostLoginBadRequest struct {
	Payload string
}

// IsSuccess returns true when this post login bad request response has a 2xx status code
func (o *PostLoginBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post login bad request response has a 3xx status code
func (o *PostLoginBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post login bad request response has a 4xx status code
func (o *PostLoginBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this post login bad request response has a 5xx status code
func (o *PostLoginBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this post login bad request response a status code equal to that given
func (o *PostLoginBadRequest) IsCode(code int) bool {
	return code == 400
}

// Code gets the status code for the post login bad request response
func (o *PostLoginBadRequest) Code() int {
	return 400
}

func (o *PostLoginBadRequest) Error() string {
	return fmt.Sprintf("[POST /login][%d] postLoginBadRequest  %+v", 400, o.Payload)
}

func (o *PostLoginBadRequest) String() string {
	return fmt.Sprintf("[POST /login][%d] postLoginBadRequest  %+v", 400, o.Payload)
}

func (o *PostLoginBadRequest) GetPayload() string {
	return o.Payload
}

func (o *PostLoginBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPostLoginInternalServerError creates a PostLoginInternalServerError with default headers values
func NewPostLoginInternalServerError() *PostLoginInternalServerError {
	return &PostLoginInternalServerError{}
}

/*
PostLoginInternalServerError describes a response with status code 500, with default header values.

Error.
*/
type PostLoginInternalServerError struct {
	Payload string
}

// IsSuccess returns true when this post login internal server error response has a 2xx status code
func (o *PostLoginInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this post login internal server error response has a 3xx status code
func (o *PostLoginInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this post login internal server error response has a 4xx status code
func (o *PostLoginInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this post login internal server error response has a 5xx status code
func (o *PostLoginInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this post login internal server error response a status code equal to that given
func (o *PostLoginInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the post login internal server error response
func (o *PostLoginInternalServerError) Code() int {
	return 500
}

func (o *PostLoginInternalServerError) Error() string {
	return fmt.Sprintf("[POST /login][%d] postLoginInternalServerError  %+v", 500, o.Payload)
}

func (o *PostLoginInternalServerError) String() string {
	return fmt.Sprintf("[POST /login][%d] postLoginInternalServerError  %+v", 500, o.Payload)
}

func (o *PostLoginInternalServerError) GetPayload() string {
	return o.Payload
}

func (o *PostLoginInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
