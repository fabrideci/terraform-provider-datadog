package utils

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
)

// DatadogApiClient is a custom HTTP client used to interact with some Datadog APIs not available in datadog-api-client-go
type DatadogApiClient interface {
	SendRequest(method, path string, body any) ([]byte, *http.Response, error)
}

func NewDatadogApiClient(client *datadog.APIClient, auth context.Context) DatadogApiClient {
	return customClientImpl{client, auth}
}

type customClientImpl struct {
	client *datadog.APIClient
	auth   context.Context
}

func (c customClientImpl) SendRequest(method, path string, body any) ([]byte, *http.Response, error) {
	req, err := c.buildRequest(method, path, body)
	if err != nil {
		return nil, nil, err
	}

	httpRes, err := c.client.CallAPI(req)
	if err != nil {
		return nil, nil, err
	}

	var bodyResByte []byte
	bodyResByte, err = io.ReadAll(httpRes.Body)
	defer httpRes.Body.Close()
	if err != nil {
		return nil, httpRes, err
	}

	if httpRes.StatusCode >= 300 {
		newErr := CustomRequestAPIError{
			body:  bodyResByte,
			error: httpRes.Status,
		}
		return nil, httpRes, newErr
	}

	return bodyResByte, httpRes, nil
}

func (c customClientImpl) buildRequest(method, path string, body interface{}) (*http.Request, error) {
	var (
		localVarPostBody        interface{}
		localVarPath            string
		localVarQueryParams     url.Values
		localVarFormQueryParams url.Values
		localVarFormFile        *datadog.FormFile
	)

	localBasePath, err := c.client.GetConfig().ServerURLWithContext(c.auth, "")
	if err != nil {
		return nil, err
	}
	localVarPath = localBasePath + path

	localVarHeaderParams := make(map[string]string)
	localVarHeaderParams["Content-Type"] = "application/json"

	localVarHTTPHeaderAccepts := make(map[string]string)
	localVarHTTPHeaderAccepts["Accept"] = "application/json"

	if body != nil {
		localVarPostBody = body
	}

	if c.auth != nil {
		if auth, ok := c.auth.Value(datadog.ContextAPIKeys).(map[string]datadog.APIKey); ok {
			if apiKey, ok := auth["apiKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["DD-API-KEY"] = key
			}
		}
	}
	if c.auth != nil {
		if auth, ok := c.auth.Value(datadog.ContextAPIKeys).(map[string]datadog.APIKey); ok {
			if apiKey, ok := auth["appKeyAuth"]; ok {
				var key string
				if apiKey.Prefix != "" {
					key = apiKey.Prefix + " " + apiKey.Key
				} else {
					key = apiKey.Key
				}
				localVarHeaderParams["DD-APPLICATION-KEY"] = key
			}
		}
	}

	return c.client.PrepareRequest(c.auth, localVarPath, method, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormQueryParams, localVarFormFile)
}

// CustomRequestAPIError Provides access to the body, and error on returned errors.
type CustomRequestAPIError struct {
	body  []byte
	error string
}

// Error returns non-empty string if there was an error.
func (e CustomRequestAPIError) Error() string {
	return e.error
}

// Body returns the raw bytes of the response
func (e CustomRequestAPIError) Body() []byte {
	return e.body
}
