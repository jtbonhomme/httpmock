# httpmock

This lib aims at providing a way to test outbound http request.
It is forked from https://git.sr.ht/~ewintr/go-kit. I only renamed the package from `go-kit/test` into `httpmock`.

See the article written by [Erik Winter](https://erikwinter.nl/about/) here: https://erikwinter.nl/articles/2020/unit-test-outbound-http-requests-in-golang/

## Usage

`go get github.com/jtbonhomme/httpmock``

## Example

This code extract shows how to assert that an external HTTP endpoint has been called during a test:

```go
package mypackage_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/izysaas/backend/libs/httpmock"
)

const (
	path string = "/tenants/"
)

func TestCreateUserOK(t *testing.T) {
	tenantID := uuid.New()
    // Create a mock external http endpoint
	var record httpmock.MockAssertion
	mockServer := httpmock.NewMockServer(&record, httpmock.MockServerProcedure{
		URI:        path + tenantID.String(),
		HTTPMethod: http.MethodGet,
		Response: httpmock.MockResponse{
			StatusCode: http.StatusOK,
			Body:       nil,
		},
	})

	// mockServer URL to be called during the test is like
	// "http://127.0.0.1:63377"

	split := strings.Split(mockServer.URL, ":")
	mockProtocol := split[0]
	mockHost := strings.Split(split[1], "//")[1]
	mockPort := split[2]

    // Trigger the function you want to test
    // ...

    // Check an external http endpoint has been called
	httpmock.Equals(t, 1, record.Hits(path+tenantID.String(), http.MethodGet))
}
```