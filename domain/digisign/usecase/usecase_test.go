package usecase

import (
	"errors"
	"los-int-digisign/domain/digisign/repository"
	"los-int-digisign/model/request"
	"los-int-digisign/model/response"
	"los-int-digisign/shared/constant"
	"los-int-digisign/shared/httpclient"
	"net/http"
	"os"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestActivation(t *testing.T) {
	// mocks
	mockRepo := new(repository.MockRepository)
	mockHttp := new(httpclient.MockHttpClient)

	testcases := []struct {
		label          string
		request        request.ActivationRequest
		response       response.ActivationResponse
		expected_error error

		err_http error
		code     int
		body     string
	}{
		{
			label: "TEST_PASS",
			request: request.ActivationRequest{
				JsonFile: request.JsonFileActivation{
					UserID: "test",
					Email:  "test",
				},
			},
			code: http.StatusOK,
			body: `{"JSONFile": {"UserID": "test", "Email": "test"}}`,
		},
		{
			label: "TEST_FAILED_UNMARSHAL",
			request: request.ActivationRequest{
				JsonFile: request.JsonFileActivation{
					UserID: "test",
					Email:  "test",
				},
			},
			code:           http.StatusOK,
			expected_error: errors.New("error while unmarshal activation response: invalid character 'a' looking for beginning of value"),
			body:           `asdasdasd'asdasdasd`,
		},
		{
			label: "TEST_FAILED_REQUEST_API_ACTIVATION",
			request: request.ActivationRequest{
				JsonFile: request.JsonFileActivation{
					UserID: "test",
					Email:  "test",
				},
			},
			code:           http.StatusInternalServerError,
			expected_error: errors.New("error while do activation: error"),
			err_http:       errors.New("error"),
		},
	}

	for _, test := range testcases {
		t.Run(test.label, func(t *testing.T) {
			url := os.Getenv("DIGISIGN_BASE_URL") + os.Getenv("DIGISIGN_ACTIVATION_URL")

			// mock
			rst := resty.New()
			httpmock.ActivateNonDefault(rst.GetClient())
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(constant.METHOD_POST, url, httpmock.NewStringResponder(test.code, test.body))
			resp, _ := rst.R().Post(url)

			header := map[string]string{
				"Content-Type":  "multipart/form-data",
				"Authorization": os.Getenv("Bearer ") + os.Getenv("DIGISIGN_TOKEN"),
			}

			mockHttp.On("DigiAPI", url, http.MethodPost, mock.Anything, "", header, 30, "").Return(resp, test.err_http).Once()

			// usecase
			usecase := NewUsecase(mockRepo, mockHttp)
			res, err := usecase.Activation(test.request)

			// assert
			assert.Equal(t, test.expected_error, err)
			assert.Equal(t, test.response, res)
		})
	}
}
