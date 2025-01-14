package handlers

import (
	"net/http"

	"github.com/captain-corp/captain/flash"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/utils"

	"github.com/gofiber/fiber/v2"
)

// ListPages handles the GET /admin/pages route
func (h *AdminHandlers) ListPages(c *fiber.Ctx) error {
	pages, err := h.repos.Pages.FindAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Render("admin_pages", fiber.Map{
		"title": "Pages",
		"pages": pages,
	})
}

// ShowCreatePage handles the GET /admin/pages/new route
func (h *AdminHandlers) ShowCreatePage(c *fiber.Ctx) error {
	return c.Render("admin_create_page", fiber.Map{
		"title": "Create Page",
		"page":  &models.Page{},
	})
}

// EditPage handles the GET /admin/pages/:id/edit route
func (h *AdminHandlers) EditPage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_edit_page", fiber.Map{
		"title": "Edit Page",
		"page":  page,
	})
}

// ConfirmDeletePage handles the GET /admin/pages/:id/delete route
func (h *AdminHandlers) ConfirmDeletePage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("admin_404", fiber.Map{})
	}

	return c.Render("admin_confirm_delete_page", fiber.Map{
		"title": "Confirm page deletion",
		"page":  page,
	})
}

// DeletePage handles the POST /admin/pages/:id/delete route
func (h *AdminHandlers) DeletePage(c *fiber.Ctx) error {
	id, err := utils.ParseUint(c.Params("id"))
	if err != nil {
		flash.Error(c, "Invalid page ID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid page ID",
			"redirect": "/admin/pages",
		})
	}

	page, err := h.repos.Pages.FindByID(id)
	if err != nil {
		flash.Error(c, "Page not found")
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error":    "Page not found",
			"redirect": "/admin/pages",
		})
	}

	// Check if page is referenced by menu items
	var menuItemCount int64
	err = h.repos.Pages.CountRelatedMenuItems(page.ID, &menuItemCount)

	if err != nil {
		flash.Error(c, "Failed to count related menu items")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to count related menu items",
			"redirect": "/admin/pages",
		})
	}

	if menuItemCount > 0 {
		flash.Error(c, "Page is referenced by menu items")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":    "Page is referenced by menu items",
			"redirect": "/admin/pages",
		})
	}

	// Delete page
	if err := h.repos.Pages.Delete(page); err != nil {
		flash.Error(c, "Failed to delete page")
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":    "Failed to delete page",
			"redirect": "/admin/pages",
		})
	}

	flash.Success(c, "Page deleted successfully")
	return c.JSON(fiber.Map{
		"message":  "Page deleted successfully",
		"redirect": "/admin/pages",
	})
}
