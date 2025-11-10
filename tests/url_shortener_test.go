package tests

// Integration tests for URL shortener API.
// These tests require a running server on localhost:8080.
// Start the server before running these tests:
//   go run cmd/url-shortener/main.go

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"URL-shortener/internal/http-server/handlers/url/save"
	"URL-shortener/internal/lib/api"
	"URL-shortener/internal/lib/random"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid url",
			alias: gofakeit.Word(),
			error: "Field URL is not a valid URL",
		},
		{
			name:  "Invalid Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").WithJSON(save.Request{
				URL:   tc.url,
				Alias: tc.alias,
			}).
				WithBasicAuth("admin", "admin").
				Expect().Status(http.StatusOK).JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)

			// Удаляем URL
			reqDel := e.DELETE("/url/"+alias).
				WithBasicAuth("admin", "admin").
				Expect().Status(http.StatusOK).JSON().Object()

			reqDel.Value("status").String().IsEqual("OK")

			// Проверяем, что после удаления URL не найден
			testRedirectNotFound(t, e, alias)
		})
	}
}

func testRedirect(t *testing.T, alias, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}

func testRedirectNotFound(t *testing.T, e *httpexpect.Expect, alias string) {
	// После удаления URL, GET запрос должен вернуть JSON с ошибкой "URL not found"
	resp := e.GET("/" + alias).
		Expect().Status(http.StatusOK).JSON().Object()

	resp.Value("status").String().IsEqual("ERROR")
	resp.Value("error").String().IsEqual("URL not found")
}
