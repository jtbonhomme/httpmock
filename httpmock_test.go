//nolint
package httpmock_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"testing"

	"github.com/izysaas/go-kit/httpmock"
)

func TestHTTPMock(t *testing.T) {

	procs := []httpmock.MockServerProcedure{
		httpmock.MockServerProcedure{
			URI:        "/",
			HTTPMethod: "GET",
			Response: httpmock.MockResponse{
				Body: []byte("getRoot"),
			},
		},
		httpmock.MockServerProcedure{
			URI:        "/",
			HTTPMethod: "POST",
			Response: httpmock.MockResponse{
				Body: []byte("postRoot"),
			},
		},
		httpmock.MockServerProcedure{
			URI:        "/get/header",
			HTTPMethod: "GET",
			Response: httpmock.MockResponse{
				StatusCode: http.StatusAccepted,
				Headers: http.Header{
					"some-key": []string{"some-value"},
				},
				Body: []byte("getResponseHeader"),
			},
		},
		httpmock.MockServerProcedure{
			URI:        "/get/auth",
			HTTPMethod: "GET",
			Response: httpmock.MockResponse{
				Body: []byte("getRootAuth"),
			},
		},
		httpmock.MockServerProcedure{
			URI:        "/my_account",
			HTTPMethod: "GET",
			Response: httpmock.MockResponse{
				Body: []byte("getAccount"),
			},
		},
		httpmock.MockServerProcedure{
			URI:        "/my_account.json",
			HTTPMethod: "GET",
			Response: httpmock.MockResponse{
				Body: []byte("getAccountJSON"),
			},
		},
	}

	var record httpmock.MockAssertion
	testMockServer := httpmock.NewMockServer(&record, procs...)

	type mockRequest struct {
		uri            string
		method         string
		user, password string
		header         http.Header
		body           []byte
		hits           int
	}

	canonical := textproto.CanonicalMIMEHeaderKey
	tcs := []struct {
		m        string
		request  mockRequest
		response httpmock.MockResponse
	}{
		{
			m: "method get root path",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				hits:   2,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method get root path with headers",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				header: http.Header{
					canonical("input-header-key"): []string{"Just the Value"},
				},
				hits: 2,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method get root path with body",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				body:   []byte("input"),
				hits:   2,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method get root path with headers and body",
			request: mockRequest{
				uri:    "/",
				method: http.MethodGet,
				header: http.Header{
					canonical("input-header-key"): []string{"Just the Value"},
				},
				body: []byte("input"),
				hits: 2,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("getRoot"),
			},
		},
		{
			m: "method post root path",
			request: mockRequest{
				uri:    "/",
				method: http.MethodPost,
				hits:   2,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("postRoot"),
			},
		},
		{
			m: "method post root path with basic authentication",
			request: mockRequest{
				uri:      "/",
				method:   http.MethodPost,
				user:     "my-user",
				password: "my-password",
				hits:     1,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusOK,
				Body:       []byte("postRoot"),
			},
		},
		{
			m: "unmatched uri path",
			request: mockRequest{
				uri:    "/unmatched",
				method: http.MethodGet,
				hits:   0,
			},
			response: httpmock.MockResponse{
				StatusCode: http.StatusNotFound,
				Body:       []byte{},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.m, func(t *testing.T) {
			httpmock.OK(t, record.Reset())

			for _ = range make([]int, tc.request.hits) {
				_url, errU := url.Parse(testMockServer.URL + tc.request.uri)
				httpmock.OK(t, errU)

				req, errReq := http.NewRequest(
					tc.request.method,
					_url.String(),
					bytes.NewReader(tc.request.body),
				)
				httpmock.OK(t, errReq)

				for k, v := range tc.request.header {
					req.Header[k] = v
				}

				// testing authentication in the request
				if len(tc.request.user) > 0 || len(tc.request.password) > 0 {
					req.SetBasicAuth(tc.request.user, tc.request.password)

					if tc.request.header == nil {
						tc.request.header = make(http.Header)
					}

					auth := tc.request.user + ":" + tc.request.password
					tc.request.header["Authorization"] = []string{
						fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))}
				}

				client := new(http.Client)
				resp, errResp := client.Do(req)
				httpmock.OK(t, errResp)

				actualBody, err := ioutil.ReadAll(resp.Body)
				httpmock.OK(t, err)
				defer resp.Body.Close()

				httpmock.Equals(t, tc.response.StatusCode, resp.StatusCode)
				httpmock.Equals(t, tc.response.Body, actualBody)
			}
			httpmock.Equals(t, tc.request.hits, record.Hits(tc.request.uri, tc.request.method))

			// assert if all request had the correct header
			for _, h := range record.Headers(tc.request.uri, tc.request.method) {
				httpmock.Equals(t, tc.request.header, h)
			}

			// assert if all request had the correct body
			for _, b := range record.Body(tc.request.uri, tc.request.method) {
				httpmock.Equals(t, tc.request.body, b)
			}
		})
	}
}
