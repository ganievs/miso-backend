package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"miso/internal/config"
	"miso/internal/handler"
	"miso/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListProviderVersions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/providers/:namespace/:type/versions")
		c.SetParamNames("namespace", "type")
		c.SetParamValues("my-namespace", "my-type")

		storage := &storage.MockStorage{
			ListFunc: func(prefix string) ([]string, error) {
				return []string{
					"providers/my-namespace/my-type/1.0.0/terraform-provider-my-type_v1.0.0",
					"providers/my-namespace/my-type/1.1.0/terraform-provider-my-type_v1.1.0",
				}, nil
			},
		}

		h := handler.NewHandler(storage, config.S3{})

		if assert.NoError(t, h.ListProviderVersions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"versions":[{"version":"1.0.0"},{"version":"1.1.0"}]}`, rec.Body.String())
		}
	})

	t.Run("empty", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/providers/:namespace/:type/versions")
		c.SetParamNames("namespace", "type")
		c.SetParamValues("my-namespace", "my-type")

		storage := &storage.MockStorage{
			ListFunc: func(prefix string) ([]string, error) {
				return []string{}, nil
			},
		}

		h := handler.NewHandler(storage, config.S3{})

		if assert.NoError(t, h.ListProviderVersions(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"versions":[]}`, rec.Body.String())
		}
	})

	t.Run("error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/providers/:namespace/:type/versions")
		c.SetParamNames("namespace", "type")
		c.SetParamValues("my-namespace", "my-type")

		storage := &storage.MockStorage{
			ListFunc: func(prefix string) ([]string, error) {
				return nil, echo.ErrInternalServerError
			},
		}

		h := handler.NewHandler(storage, config.S3{})

		assert.Error(t, h.ListProviderVersions(c))
	})
}

func TestDownloadProviderVersion(t *testing.T) {
	t.Run("presigned-url", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/providers/:namespace/:type/:version/download/:os/:arch")
		c.SetParamNames("namespace", "type", "version", "os", "arch")
		c.SetParamValues("my-namespace", "my-type", "1.0.0", "linux", "amd64")

		storage := &storage.MockStorage{
			GetPresignedURLFunc: func(key string) (string, error) {
				return "https://example.com/download", nil
			},
		}

		cfg := config.S3{DownloadMode: "presigned-url"}
		h := handler.NewHandler(storage, cfg)

		if assert.NoError(t, h.DownloadProviderVersion(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"download_url":"https://example.com/download"}`, rec.Body.String())
		}
	})

	t.Run("proxy", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/providers/:namespace/:type/:version/download/:os/:arch")
		c.SetParamNames("namespace", "type", "version", "os", "arch")
		c.SetParamValues("my-namespace", "my-type", "1.0.0", "linux", "amd64")

		storage := &storage.MockStorage{
			GetStreamFunc: func(key string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("file content")), nil
			},
		}

		cfg := config.S3{DownloadMode: "proxy"}
		h := handler.NewHandler(storage, cfg)

		if assert.NoError(t, h.DownloadProviderVersion(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "file content", rec.Body.String())
		}
	})
}

func TestDownloadModuleVersion(t *testing.T) {
	t.Run("presigned-url", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/modules/:namespace/:name/:provider/:version/download")
		c.SetParamNames("namespace", "name", "provider", "version")
		c.SetParamValues("my-namespace", "my-module", "my-provider", "1.0.0")

		storage := &storage.MockStorage{
			GetPresignedURLFunc: func(key string) (string, error) {
				return "https://example.com/download", nil
			},
		}

		cfg := config.S3{DownloadMode: "presigned-url"}
		h := handler.NewHandler(storage, cfg)

		if assert.NoError(t, h.DownloadModuleVersion(c)) {
			assert.Equal(t, http.StatusFound, rec.Code)
			assert.Equal(t, "https://example.com/download", rec.Header().Get(echo.HeaderLocation))
		}
	})

	t.Run("proxy", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/modules/:namespace/:name/:provider/:version/download")
		c.SetParamNames("namespace", "name", "provider", "version")
		c.SetParamValues("my-namespace", "my-module", "my-provider", "1.0.0")

		storage := &storage.MockStorage{
			GetStreamFunc: func(key string) (io.ReadCloser, error) {
				return io.NopCloser(strings.NewReader("file content")), nil
			},
		}

		cfg := config.S3{DownloadMode: "proxy"}
		h := handler.NewHandler(storage, cfg)

		if assert.NoError(t, h.DownloadModuleVersion(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "file content", rec.Body.String())
		}
	})
}
