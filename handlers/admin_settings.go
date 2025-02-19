package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/captain-corp/captain/models"
)

// ShowSettings handles the GET /admin/settings route
func (h *AdminHandlers) ShowSettings(c *fiber.Ctx) error {
	var logo *models.Media

	settings, err := h.repos.Settings.Get()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	if settings.LogoID != nil {
		logo, _ = h.repos.Media.FindByID(*settings.LogoID)
	}

	// TODO: Add logo
	// TODO: Use a favicon?

	data := fiber.Map{
		"title":    "Site Settings",
		"settings": settings,
		"logo":     logo,
	}

	return c.Render("admin_settings", data)
}
