package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/jgarrone/foaas-limiting/services/foaasapi"
	"github.com/jgarrone/foaas-limiting/services/foaasapi/mocks"
	"github.com/jgarrone/foaas-limiting/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	limitCount  = 5
	limitWindow = 100 * time.Millisecond
	svAddr      = "localhost:8888"
)

func TestServer_Run(t *testing.T) {
	// Define valid values for the request.
	protocol := "http"
	path := MessagePath
	method := "GET"
	aUser := "test-user-1"

	tearDown := setUp()
	defer tearDown()

	// Initialize the client. Trivial but it's advised to reuse it.
	client := &http.Client{}

	tests := []struct {
		name     string
		req      *http.Request
		wantCode int
	}{
		{
			name:     "invalid http method",
			req:      buildRequest(protocol, svAddr, path, "POST", aUser),
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "missing user",
			req:      buildRequest(protocol, svAddr, path, method, ""),
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "invalid user",
			req:      buildRequest(protocol, svAddr, path, method, "a\xc5z"),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "ok",
			req:      buildRequest(protocol, svAddr, path, method, aUser),
			wantCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.Do(test.req)
			assert.Nil(t, err)
			assert.Equal(t, test.wantCode, resp.StatusCode)
		})
	}
}

func TestServer_Run_limit(t *testing.T) {
	// Define required constants.
	protocol := "http"
	path := MessagePath
	method := "GET"
	aUser, anotherUser := "test-user-1", "test-user-2"

	tearDown := setUp()
	defer tearDown()

	// Initialize the client. Trivial but it's advised to reuse it.
	client := &http.Client{}

	// Let's do one request more than we are allowed to with each user, and assert that the last one returns
	// an error code for each one of them.
	for i := 0; i < limitCount+1; i++ {
		for _, user := range []string{aUser, anotherUser} {
			resp, err := client.Do(buildRequest(protocol, svAddr, path, method, user))
			assert.Nil(t, err)
			if i < limitCount {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}
	}

	// Let the server cool down for a bit. Use the same window as the limiter.
	time.Sleep(limitWindow)

	// So the following request should work again.
	resp, err := client.Do(buildRequest(protocol, svAddr, path, method, aUser))
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func buildRequest(protocol, address, path, method, userid string) *http.Request {
	url := utils.BuildURL(protocol, address, path)
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set(UserIdHeader, userid)

	return req
}

func setUp() func() {
	// Mock the FOAAS service to avoid contacting the outside world.
	serviceMock := &mocks.Service{}
	serviceMock.
		On("GetMessageFor", mock.AnythingOfType("string")).
		Return(&foaasapi.Response{Message: "mocked message"}, nil)

	limiter := NewTokenBucketLimiter(limitCount, limitWindow)
	sv := New(svAddr, limiter, serviceMock)
	go sv.Run()

	return sv.Stop
}
