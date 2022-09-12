// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/ory/oathkeeper/internal/httpclient/models"
)

// DecisionsReader is a Reader for the Decisions structure.
type DecisionsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DecisionsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDecisionsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewDecisionsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewDecisionsForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDecisionsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewDecisionsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewDecisionsOK creates a DecisionsOK with default headers values
func NewDecisionsOK() *DecisionsOK {
	return &DecisionsOK{}
}

/* DecisionsOK describes a response with status code 200, with default header values.

An empty response
*/
type DecisionsOK struct {
}

// IsSuccess returns true when this decisions o k response has a 2xx status code
func (o *DecisionsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this decisions o k response has a 3xx status code
func (o *DecisionsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this decisions o k response has a 4xx status code
func (o *DecisionsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this decisions o k response has a 5xx status code
func (o *DecisionsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this decisions o k response a status code equal to that given
func (o *DecisionsOK) IsCode(code int) bool {
	return code == 200
}

func (o *DecisionsOK) Error() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsOK ", 200)
}

func (o *DecisionsOK) String() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsOK ", 200)
}

func (o *DecisionsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDecisionsUnauthorized creates a DecisionsUnauthorized with default headers values
func NewDecisionsUnauthorized() *DecisionsUnauthorized {
	return &DecisionsUnauthorized{}
}

/* DecisionsUnauthorized describes a response with status code 401, with default header values.

genericError
*/
type DecisionsUnauthorized struct {
	Payload *models.GenericError
}

// IsSuccess returns true when this decisions unauthorized response has a 2xx status code
func (o *DecisionsUnauthorized) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this decisions unauthorized response has a 3xx status code
func (o *DecisionsUnauthorized) IsRedirect() bool {
	return false
}

// IsClientError returns true when this decisions unauthorized response has a 4xx status code
func (o *DecisionsUnauthorized) IsClientError() bool {
	return true
}

// IsServerError returns true when this decisions unauthorized response has a 5xx status code
func (o *DecisionsUnauthorized) IsServerError() bool {
	return false
}

// IsCode returns true when this decisions unauthorized response a status code equal to that given
func (o *DecisionsUnauthorized) IsCode(code int) bool {
	return code == 401
}

func (o *DecisionsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsUnauthorized  %+v", 401, o.Payload)
}

func (o *DecisionsUnauthorized) String() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsUnauthorized  %+v", 401, o.Payload)
}

func (o *DecisionsUnauthorized) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *DecisionsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDecisionsForbidden creates a DecisionsForbidden with default headers values
func NewDecisionsForbidden() *DecisionsForbidden {
	return &DecisionsForbidden{}
}

/* DecisionsForbidden describes a response with status code 403, with default header values.

genericError
*/
type DecisionsForbidden struct {
	Payload *models.GenericError
}

// IsSuccess returns true when this decisions forbidden response has a 2xx status code
func (o *DecisionsForbidden) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this decisions forbidden response has a 3xx status code
func (o *DecisionsForbidden) IsRedirect() bool {
	return false
}

// IsClientError returns true when this decisions forbidden response has a 4xx status code
func (o *DecisionsForbidden) IsClientError() bool {
	return true
}

// IsServerError returns true when this decisions forbidden response has a 5xx status code
func (o *DecisionsForbidden) IsServerError() bool {
	return false
}

// IsCode returns true when this decisions forbidden response a status code equal to that given
func (o *DecisionsForbidden) IsCode(code int) bool {
	return code == 403
}

func (o *DecisionsForbidden) Error() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsForbidden  %+v", 403, o.Payload)
}

func (o *DecisionsForbidden) String() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsForbidden  %+v", 403, o.Payload)
}

func (o *DecisionsForbidden) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *DecisionsForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDecisionsNotFound creates a DecisionsNotFound with default headers values
func NewDecisionsNotFound() *DecisionsNotFound {
	return &DecisionsNotFound{}
}

/* DecisionsNotFound describes a response with status code 404, with default header values.

genericError
*/
type DecisionsNotFound struct {
	Payload *models.GenericError
}

// IsSuccess returns true when this decisions not found response has a 2xx status code
func (o *DecisionsNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this decisions not found response has a 3xx status code
func (o *DecisionsNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this decisions not found response has a 4xx status code
func (o *DecisionsNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this decisions not found response has a 5xx status code
func (o *DecisionsNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this decisions not found response a status code equal to that given
func (o *DecisionsNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *DecisionsNotFound) Error() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsNotFound  %+v", 404, o.Payload)
}

func (o *DecisionsNotFound) String() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsNotFound  %+v", 404, o.Payload)
}

func (o *DecisionsNotFound) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *DecisionsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDecisionsInternalServerError creates a DecisionsInternalServerError with default headers values
func NewDecisionsInternalServerError() *DecisionsInternalServerError {
	return &DecisionsInternalServerError{}
}

/* DecisionsInternalServerError describes a response with status code 500, with default header values.

genericError
*/
type DecisionsInternalServerError struct {
	Payload *models.GenericError
}

// IsSuccess returns true when this decisions internal server error response has a 2xx status code
func (o *DecisionsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this decisions internal server error response has a 3xx status code
func (o *DecisionsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this decisions internal server error response has a 4xx status code
func (o *DecisionsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this decisions internal server error response has a 5xx status code
func (o *DecisionsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this decisions internal server error response a status code equal to that given
func (o *DecisionsInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *DecisionsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsInternalServerError  %+v", 500, o.Payload)
}

func (o *DecisionsInternalServerError) String() string {
	return fmt.Sprintf("[GET /decisions][%d] decisionsInternalServerError  %+v", 500, o.Payload)
}

func (o *DecisionsInternalServerError) GetPayload() *models.GenericError {
	return o.Payload
}

func (o *DecisionsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GenericError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
