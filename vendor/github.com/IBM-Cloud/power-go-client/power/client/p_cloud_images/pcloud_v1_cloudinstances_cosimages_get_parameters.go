// Code generated by go-swagger; DO NOT EDIT.

package p_cloud_images

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
)

// NewPcloudV1CloudinstancesCosimagesGetParams creates a new PcloudV1CloudinstancesCosimagesGetParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPcloudV1CloudinstancesCosimagesGetParams() *PcloudV1CloudinstancesCosimagesGetParams {
	return &PcloudV1CloudinstancesCosimagesGetParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPcloudV1CloudinstancesCosimagesGetParamsWithTimeout creates a new PcloudV1CloudinstancesCosimagesGetParams object
// with the ability to set a timeout on a request.
func NewPcloudV1CloudinstancesCosimagesGetParamsWithTimeout(timeout time.Duration) *PcloudV1CloudinstancesCosimagesGetParams {
	return &PcloudV1CloudinstancesCosimagesGetParams{
		timeout: timeout,
	}
}

// NewPcloudV1CloudinstancesCosimagesGetParamsWithContext creates a new PcloudV1CloudinstancesCosimagesGetParams object
// with the ability to set a context for a request.
func NewPcloudV1CloudinstancesCosimagesGetParamsWithContext(ctx context.Context) *PcloudV1CloudinstancesCosimagesGetParams {
	return &PcloudV1CloudinstancesCosimagesGetParams{
		Context: ctx,
	}
}

// NewPcloudV1CloudinstancesCosimagesGetParamsWithHTTPClient creates a new PcloudV1CloudinstancesCosimagesGetParams object
// with the ability to set a custom HTTPClient for a request.
func NewPcloudV1CloudinstancesCosimagesGetParamsWithHTTPClient(client *http.Client) *PcloudV1CloudinstancesCosimagesGetParams {
	return &PcloudV1CloudinstancesCosimagesGetParams{
		HTTPClient: client,
	}
}

/*
PcloudV1CloudinstancesCosimagesGetParams contains all the parameters to send to the API endpoint

	for the pcloud v1 cloudinstances cosimages get operation.

	Typically these are written to a http.Request.
*/
type PcloudV1CloudinstancesCosimagesGetParams struct {

	/* CloudInstanceID.

	   Cloud Instance ID of a PCloud Instance
	*/
	CloudInstanceID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the pcloud v1 cloudinstances cosimages get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PcloudV1CloudinstancesCosimagesGetParams) WithDefaults() *PcloudV1CloudinstancesCosimagesGetParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the pcloud v1 cloudinstances cosimages get params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PcloudV1CloudinstancesCosimagesGetParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) WithTimeout(timeout time.Duration) *PcloudV1CloudinstancesCosimagesGetParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) WithContext(ctx context.Context) *PcloudV1CloudinstancesCosimagesGetParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) WithHTTPClient(client *http.Client) *PcloudV1CloudinstancesCosimagesGetParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithCloudInstanceID adds the cloudInstanceID to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) WithCloudInstanceID(cloudInstanceID string) *PcloudV1CloudinstancesCosimagesGetParams {
	o.SetCloudInstanceID(cloudInstanceID)
	return o
}

// SetCloudInstanceID adds the cloudInstanceId to the pcloud v1 cloudinstances cosimages get params
func (o *PcloudV1CloudinstancesCosimagesGetParams) SetCloudInstanceID(cloudInstanceID string) {
	o.CloudInstanceID = cloudInstanceID
}

// WriteToRequest writes these params to a swagger request
func (o *PcloudV1CloudinstancesCosimagesGetParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param cloud_instance_id
	if err := r.SetPathParam("cloud_instance_id", o.CloudInstanceID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}