package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/captain-corp/captain/models"
)

// ShowSettings handles the GET /admin/settings route
func (h *AdminHandlers) ShowSettings(c *fiber.Ctx) error {
	var logo *models.Media
	var logoUrl string

	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	if settings.LogoID != nil {
		logo, _ = h.repos.Media.FindByID(*settings.LogoID)
		logoUrl = "/media/" + logo.Path
	}

	data := fiber.Map{
		"title":    "Site Settings",
		"settings": settings,
		"logoUrl":  logoUrl,
	}

	return c.Render("admin_settings", data)
}
