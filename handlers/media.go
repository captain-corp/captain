package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"codeinstyle.io/captain/config"
	"codeinstyle.io/captain/repository"
	"codeinstyle.io/captain/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ServeMedia serves media files from the configured storage provider
func ServeMedia(repositories *repository.Repositories, cfg *config.Config) gin.HandlerFunc {
	// Initialize storage provider
	var provider storage.Provider
	var err error

	switch cfg.Storage.Provider {
	case "s3":
		provider, err = storage.NewS3Provider(cfg.Storage.S3.Bucket, cfg.Storage.S3.Region, cfg.Storage.S3.Endpoint, cfg.Storage.S3.AccessKey, cfg.Storage.S3.SecretKey)
	default: // "local"
		provider, err = storage.NewLocalProvider(cfg.Storage.LocalPath)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize storage provider: %v", err))
	}

	return func(c *gin.Context) {
		// Get path and trim leading slash if present
		path := c.Param("path")
		if path == "" {
			c.String(http.StatusBadRequest, "No path provided")
			return
		}
		// Trim the leading slash as paths are stored without it in the database
		path = strings.TrimPrefix(path, "/")

		// Query the media from the database
		media, err := repositories.Media.FindByPath(path)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// If the media doesn't exist, return a 404
				c.String(http.StatusNotFound, "Media not found")
				return
			}
			c.String(http.StatusInternalServerError, "Error retrieving media")
			return
		}

		// Generate ETag based on last modified time and size
		etag := fmt.Sprintf(`"%x-%x"`, media.UpdatedAt.Unix(), media.Size)

		// Check If-None-Match header
		if match := c.GetHeader("If-None-Match"); match != "" {
			if match == etag {
				c.Status(http.StatusNotModified)
				return
			}
		}

		// Get file from storage provider
		file, err := provider.Get(path)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error retrieving media file")
			return
		}
		defer file.Close()

		// Set content type header
		c.Header("Content-Type", media.MimeType)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", path))
		c.Header("ETag", etag)
		c.Header("Last-Modified", media.UpdatedAt.Format(http.TimeFormat))
		c.Header("Content-Length", fmt.Sprintf("%d", media.Size))
		c.Header("Accept-Ranges", "bytes")
		c.Header("Cache-Control", "public, max-age=31536000")

		// Stream the file to the response
		if _, err := io.Copy(c.Writer, file); err != nil {
			c.String(http.StatusInternalServerError, "Error streaming media file")
			return
		}
	}
}
