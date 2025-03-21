package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/akamaaru/url-shortener/internal/http-server/handlers/redirect"
	"github.com/akamaaru/url-shortener/internal/http-server/handlers/redirect/mocks"
	"github.com/akamaaru/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/akamaaru/url-shortener/internal/lib/api"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		// {
		// 	name:  "Empty alias",
		// 	alias: "",
		// 	url:   "https://google.com",
		// },
		// {
		// 	name:      "Empty URL",
		// 	url:       "",
		// 	alias:     "some_alias",
		// 	respError: "field URL is a required field",
		// },
		// {
		// 	name:      "Invalid URL",
		// 	url:       "some invalid URL",
		// 	alias:     "some_alias",
		// 	respError: "field URL is not a valid URL",
		// },
		// {
		// 	name:      "SaveURL Error",
		// 	alias:     "test_alias",
		// 	url:       "https://google.com",
		// 	respError: "failed to add url",
		// 	mockError: errors.New("unexpected error"),
		// },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)
			assert.Equal(t, tc.url, redirectedToURL)
		})
	}
}