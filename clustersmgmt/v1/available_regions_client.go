/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/openshift-online/ocm-sdk-go/errors"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// AvailableRegionsClient is the client of the 'available_regions' resource.
//
// Manages collection of cloud provider regions available to a particular cloud provider account
type AvailableRegionsClient struct {
	transport http.RoundTripper
	path      string
}

// NewAvailableRegionsClient creates a new client for the 'available_regions'
// resource using the given transport to send the requests and receive the
// responses.
func NewAvailableRegionsClient(transport http.RoundTripper, path string) *AvailableRegionsClient {
	return &AvailableRegionsClient{
		transport: transport,
		path:      path,
	}
}

// Search creates a request for the 'search' method.
//
// Retrieves the list of available regions of the cloud provider.
//
// IMPORTANT: This collection doesn't currently support paging or searching, so the returned
// `page` will always be 1 and `size` and `total` will always be the total number of available regions
// of the provider.
func (c *AvailableRegionsClient) Search() *AvailableRegionsSearchRequest {
	return &AvailableRegionsSearchRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AvailableRegionsSearchRequest is the request for the 'search' method.
type AvailableRegionsSearchRequest struct {
	transport http.RoundTripper
	path      string
	query     url.Values
	header    http.Header
	body      *AWS
	page      *int
	size      *int
}

// Parameter adds a query parameter.
func (r *AvailableRegionsSearchRequest) Parameter(name string, value interface{}) *AvailableRegionsSearchRequest {
	helpers.AddValue(&r.query, name, value)
	return r
}

// Header adds a request header.
func (r *AvailableRegionsSearchRequest) Header(name string, value interface{}) *AvailableRegionsSearchRequest {
	helpers.AddHeader(&r.header, name, value)
	return r
}

// Body sets the value of the 'body' parameter.
//
// AWS account details
func (r *AvailableRegionsSearchRequest) Body(value *AWS) *AvailableRegionsSearchRequest {
	r.body = value
	return r
}

// Page sets the value of the 'page' parameter.
//
// Index of the returned page, where one corresponds to the first page. As this
// collection doesn't support paging the result will always be `1`.
func (r *AvailableRegionsSearchRequest) Page(value int) *AvailableRegionsSearchRequest {
	r.page = &value
	return r
}

// Size sets the value of the 'size' parameter.
//
// Number of items that will be contained in the returned page. As this collection
// doesn't support paging or searching the result will always be the total number of
// regions of the provider.
func (r *AvailableRegionsSearchRequest) Size(value int) *AvailableRegionsSearchRequest {
	r.size = &value
	return r
}

// Send sends this request, waits for the response, and returns it.
//
// This is a potentially lengthy operation, as it requires network communication.
// Consider using a context and the SendContext method.
func (r *AvailableRegionsSearchRequest) Send() (result *AvailableRegionsSearchResponse, err error) {
	return r.SendContext(context.Background())
}

// SendContext sends this request, waits for the response, and returns it.
func (r *AvailableRegionsSearchRequest) SendContext(ctx context.Context) (result *AvailableRegionsSearchResponse, err error) {
	query := helpers.CopyQuery(r.query)
	if r.page != nil {
		helpers.AddValue(&query, "page", *r.page)
	}
	if r.size != nil {
		helpers.AddValue(&query, "size", *r.size)
	}
	header := helpers.CopyHeader(r.header)
	buffer := &bytes.Buffer{}
	err = writeAvailableRegionsSearchRequest(r, buffer)
	if err != nil {
		return
	}
	uri := &url.URL{
		Path:     r.path,
		RawQuery: query.Encode(),
	}
	request := &http.Request{
		Method: "POST",
		URL:    uri,
		Header: header,
		Body:   ioutil.NopCloser(buffer),
	}
	if ctx != nil {
		request = request.WithContext(ctx)
	}
	response, err := r.transport.RoundTrip(request)
	if err != nil {
		return
	}
	defer response.Body.Close()
	result = &AvailableRegionsSearchResponse{}
	result.status = response.StatusCode
	result.header = response.Header
	if result.status >= 400 {
		result.err, err = errors.UnmarshalErrorStatus(response.Body, result.status)
		if err != nil {
			return
		}
		err = result.err
		return
	}
	err = readAvailableRegionsSearchResponse(result, response.Body)
	if err != nil {
		return
	}
	return
}

// AvailableRegionsSearchResponse is the response for the 'search' method.
type AvailableRegionsSearchResponse struct {
	status int
	header http.Header
	err    *errors.Error
	items  *CloudRegionList
	page   *int
	size   *int
	total  *int
}

// Status returns the response status code.
func (r *AvailableRegionsSearchResponse) Status() int {
	if r == nil {
		return 0
	}
	return r.status
}

// Header returns header of the response.
func (r *AvailableRegionsSearchResponse) Header() http.Header {
	if r == nil {
		return nil
	}
	return r.header
}

// Error returns the response error.
func (r *AvailableRegionsSearchResponse) Error() *errors.Error {
	if r == nil {
		return nil
	}
	return r.err
}

// Items returns the value of the 'items' parameter.
//
// Retrieved list of cloud regions.
func (r *AvailableRegionsSearchResponse) Items() *CloudRegionList {
	if r == nil {
		return nil
	}
	return r.items
}

// GetItems returns the value of the 'items' parameter and
// a flag indicating if the parameter has a value.
//
// Retrieved list of cloud regions.
func (r *AvailableRegionsSearchResponse) GetItems() (value *CloudRegionList, ok bool) {
	ok = r != nil && r.items != nil
	if ok {
		value = r.items
	}
	return
}

// Page returns the value of the 'page' parameter.
//
// Index of the returned page, where one corresponds to the first page. As this
// collection doesn't support paging the result will always be `1`.
func (r *AvailableRegionsSearchResponse) Page() int {
	if r != nil && r.page != nil {
		return *r.page
	}
	return 0
}

// GetPage returns the value of the 'page' parameter and
// a flag indicating if the parameter has a value.
//
// Index of the returned page, where one corresponds to the first page. As this
// collection doesn't support paging the result will always be `1`.
func (r *AvailableRegionsSearchResponse) GetPage() (value int, ok bool) {
	ok = r != nil && r.page != nil
	if ok {
		value = *r.page
	}
	return
}

// Size returns the value of the 'size' parameter.
//
// Number of items that will be contained in the returned page. As this collection
// doesn't support paging or searching the result will always be the total number of
// regions of the provider.
func (r *AvailableRegionsSearchResponse) Size() int {
	if r != nil && r.size != nil {
		return *r.size
	}
	return 0
}

// GetSize returns the value of the 'size' parameter and
// a flag indicating if the parameter has a value.
//
// Number of items that will be contained in the returned page. As this collection
// doesn't support paging or searching the result will always be the total number of
// regions of the provider.
func (r *AvailableRegionsSearchResponse) GetSize() (value int, ok bool) {
	ok = r != nil && r.size != nil
	if ok {
		value = *r.size
	}
	return
}

// Total returns the value of the 'total' parameter.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page. As this collection doesn't support paging or
// searching the result will always be the total number of available regions of the provider.
func (r *AvailableRegionsSearchResponse) Total() int {
	if r != nil && r.total != nil {
		return *r.total
	}
	return 0
}

// GetTotal returns the value of the 'total' parameter and
// a flag indicating if the parameter has a value.
//
// Total number of items of the collection that match the search criteria,
// regardless of the size of the page. As this collection doesn't support paging or
// searching the result will always be the total number of available regions of the provider.
func (r *AvailableRegionsSearchResponse) GetTotal() (value int, ok bool) {
	ok = r != nil && r.total != nil
	if ok {
		value = *r.total
	}
	return
}
