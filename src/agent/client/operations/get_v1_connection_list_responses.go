// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/hanzezhenalex/socks5/src/agent/models"
)

// GetV1ConnectionListReader is a Reader for the GetV1ConnectionList structure.
type GetV1ConnectionListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetV1ConnectionListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetV1ConnectionListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 500:
		result := NewGetV1ConnectionListInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetV1ConnectionListOK creates a GetV1ConnectionListOK with default headers values
func NewGetV1ConnectionListOK() *GetV1ConnectionListOK {
	return &GetV1ConnectionListOK{}
}

/* GetV1ConnectionListOK describes a response with status code 200, with default header values.

Get all the connections.
*/
type GetV1ConnectionListOK struct {
	Payload models.Connections
}

// IsSuccess returns true when this get v1 connection list o k response has a 2xx status code
func (o *GetV1ConnectionListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get v1 connection list o k response has a 3xx status code
func (o *GetV1ConnectionListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 connection list o k response has a 4xx status code
func (o *GetV1ConnectionListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 connection list o k response has a 5xx status code
func (o *GetV1ConnectionListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get v1 connection list o k response a status code equal to that given
func (o *GetV1ConnectionListOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetV1ConnectionListOK) Error() string {
	return fmt.Sprintf("[GET /v1/connection/list][%d] getV1ConnectionListOK  %+v", 200, o.Payload)
}

func (o *GetV1ConnectionListOK) String() string {
	return fmt.Sprintf("[GET /v1/connection/list][%d] getV1ConnectionListOK  %+v", 200, o.Payload)
}

func (o *GetV1ConnectionListOK) GetPayload() models.Connections {
	return o.Payload
}

func (o *GetV1ConnectionListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetV1ConnectionListInternalServerError creates a GetV1ConnectionListInternalServerError with default headers values
func NewGetV1ConnectionListInternalServerError() *GetV1ConnectionListInternalServerError {
	return &GetV1ConnectionListInternalServerError{}
}

/* GetV1ConnectionListInternalServerError describes a response with status code 500, with default header values.

Error.
*/
type GetV1ConnectionListInternalServerError struct {
	Payload string
}

// IsSuccess returns true when this get v1 connection list internal server error response has a 2xx status code
func (o *GetV1ConnectionListInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get v1 connection list internal server error response has a 3xx status code
func (o *GetV1ConnectionListInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get v1 connection list internal server error response has a 4xx status code
func (o *GetV1ConnectionListInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get v1 connection list internal server error response has a 5xx status code
func (o *GetV1ConnectionListInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get v1 connection list internal server error response a status code equal to that given
func (o *GetV1ConnectionListInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *GetV1ConnectionListInternalServerError) Error() string {
	return fmt.Sprintf("[GET /v1/connection/list][%d] getV1ConnectionListInternalServerError  %+v", 500, o.Payload)
}

func (o *GetV1ConnectionListInternalServerError) String() string {
	return fmt.Sprintf("[GET /v1/connection/list][%d] getV1ConnectionListInternalServerError  %+v", 500, o.Payload)
}

func (o *GetV1ConnectionListInternalServerError) GetPayload() string {
	return o.Payload
}

func (o *GetV1ConnectionListInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
