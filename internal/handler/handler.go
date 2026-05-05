package handler

import (
	"io"
	"miso/internal/config"
	"miso/internal/storage"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	Storage storage.Storage
	Config  config.S3
}

func NewHandler(storage storage.Storage, config config.S3) *Handler {
	return &Handler{
		Storage: storage,
		Config:  config,
	}
}

func (h *Handler) ListProviderVersions(c echo.Context) error {
	namespace := c.Param("namespace")
	typeName := c.Param("type")

	prefix := "providers/" + namespace + "/" + typeName + "/"
	keys, err := h.Storage.List(prefix)
	if err != nil {
		return err
	}

	versionSet := make(map[string]struct{})
	for _, key := range keys {
		version := strings.Split(strings.TrimPrefix(key, prefix), "/")[0]
		versionSet[version] = struct{}{}
	}

	versions := []map[string]interface{}{}
	for version := range versionSet {
		versions = append(versions, map[string]interface{}{"version": version})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"versions": versions,
	})
}

func (h *Handler) DownloadProviderVersion(c echo.Context) error {
	namespace := c.Param("namespace")
	typeName := c.Param("type")
	version := c.Param("version")
	os := c.Param("os")
	arch := c.Param("arch")

	key := "providers/" + namespace + "/" + typeName + "/" + version + "/" + os + "/" + arch + "/terraform-provider-" + typeName + "_v" + version

	if h.Config.DownloadMode == "proxy" {
		return h.proxyDownload(c, key)
	}

	downloadURL, err := h.Storage.GetPresignedURL(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"download_url": downloadURL,
	})
}

func (h *Handler) ListModuleVersions(c echo.Context) error {
	namespace := c.Param("namespace")
	name := c.Param("name")
	provider := c.Param("provider")

	prefix := "modules/" + namespace + "/" + name + "/" + provider + "/"
	keys, err := h.Storage.List(prefix)
	if err != nil {
		return err
	}

	versionSet := make(map[string]struct{})
	for _, key := range keys {
		version := strings.Split(strings.TrimPrefix(key, prefix), "/")[0]
		versionSet[version] = struct{}{}
	}

	versions := []map[string]interface{}{}
	for version := range versionSet {
		versions = append(versions, map[string]interface{}{"version": version})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"modules": []map[string]interface{}{
			{"versions": versions},
		},
	})
}

func (h *Handler) DownloadModuleVersion(c echo.Context) error {
	namespace := c.Param("namespace")
	name := c.Param("name")
	provider := c.Param("provider")
	version := c.Param("version")

	key := "modules/" + namespace + "/" + name + "/" + provider + "/" + version + "/module.zip"

	if h.Config.DownloadMode == "proxy" {
		return h.proxyDownload(c, key)
	}

	downloadURL, err := h.Storage.GetPresignedURL(key)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, downloadURL)
}

func (h *Handler) proxyDownload(c echo.Context, key string) error {
	stream, err := h.Storage.GetStream(key)
	if err != nil {
		return err
	}
	defer stream.Close()

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	c.Response().WriteHeader(http.StatusOK)
	_, err = io.Copy(c.Response().Writer, stream)
	return err
}
