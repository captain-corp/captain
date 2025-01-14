package handlers

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/captain-corp/captain/config"
	"github.com/captain-corp/captain/models"
	"github.com/captain-corp/captain/repository"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// PublicHandlers handles all public routes
type PublicHandlers struct {
	*BaseHandlers
}

// NewPublicHandlers creates a new public handlers instance
func NewPublicHandlers(repos *repository.Repositories, cfg *config.Config) *PublicHandlers {
	return &PublicHandlers{
		BaseHandlers: NewBaseHandlers(repos, cfg),
	}
}

func (h *PublicHandlers) GetPostBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	post, err := h.repos.Posts.FindBySlug(slug)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	// Render markdown content
	post.Content = renderMarkdown(post.Content)

	return c.Render("post", fiber.Map{
		"title": post.Title,
		"post":  post,
	})
}

func (h *PublicHandlers) ListPosts(c *fiber.Ctx) error {
	settings := c.Locals("settings").(*models.Settings)

	// Get page number from query params
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	// Get user from context
	user := c.Locals("user")
	var posts []models.Post
	var total int64

	if user != nil {
		// Logged-in users can see all posts
		posts, total, err = h.repos.Posts.FindAllPaginated(page, settings.PostsPerPage)
	} else {
		// Anonymous users can only see visible posts
		posts, total, err = h.repos.Posts.FindVisiblePaginated(page, settings.PostsPerPage)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(settings.PostsPerPage)))

	processPostsContent(posts)
	processPostsPublishedAt(posts)

	return c.Render("posts", fiber.Map{
		"title":       "Latest Articles",
		"posts":       posts,
		"currentPage": page,
		"totalPages":  totalPages,
		"user":        user,
		"settings":    settings,
	})
}

func (h *PublicHandlers) ListPostsByTag(c *fiber.Ctx) error {
	tagSlug := c.Params("slug")
	settings := c.Locals("settings").(*models.Settings)

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	tag, err := h.repos.Tags.FindBySlug(tagSlug)
	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("404", fiber.Map{
			"title": "Tag not found",
		})
	}

	// Get user from context
	user := c.Locals("user")
	var posts []models.Post
	var total int64

	if user != nil {
		// Logged-in users can see all posts with tag
		posts, total, err = h.repos.Posts.FindAllByTag(tag.ID, page, settings.PostsPerPage)
	} else {
		// Anonymous users can only see visible posts with tag
		posts, total, err = h.repos.Posts.FindVisibleByTag(tag.ID, page, settings.PostsPerPage)
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).Render("admin_500", fiber.Map{
			"error": err.Error(),
		})
	}

	// Process posts
	processPostsPublishedAt(posts)
	processPostsContent(posts)

	totalPages := int(math.Ceil(float64(total) / float64(settings.PostsPerPage)))

	return c.Render("tag_posts", h.addCommonData(c, fiber.Map{
		"title":      fmt.Sprintf("Posts tagged with %s", tag.Name),
		"posts":      posts,
		"tag":        tag,
		"page":       page,
		"totalPages": totalPages,
		"user":       user,
		"settings":   settings,
	}))
}

func (h *PublicHandlers) GetPageBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	page, err := h.repos.Pages.FindBySlug(slug)
	if err != nil {
		return c.Status(http.StatusNotFound).Render("404", fiber.Map{})
	}

	// Render content based on type
	if page.ContentType == "markdown" {
		page.Content = renderMarkdown(page.Content)
	}

	return c.Render("page", fiber.Map{
		"title": page.Title,
		"page":  page,
	})
}

func processPostsPublishedAt(posts []models.Post) {
	for i := range posts {
		posts[i].PublishedAt = posts[i].PublishedAt.In(time.UTC)
	}
}

func processPostsContent(posts []models.Post) {
	for i := range posts {
		if posts[i].Excerpt != nil && *posts[i].Excerpt != "" {
			continue
		}
		// Truncate content to ~200 chars if no excerpt
		content := posts[i].Content
		if len(content) > 200 {
			content = content[:200]
			if idx := strings.LastIndex(content, " "); idx > 0 {
				content = content[:idx]
			}
			content += "..."
		}
		posts[i].Excerpt = &content
		// Render markdown for excerpt
		if posts[i].Excerpt != nil {
			rendered := renderMarkdown(*posts[i].Excerpt)
			posts[i].Excerpt = &rendered
		}
	}
}

// renderMarkdown converts markdown content to HTML
func renderMarkdown(content string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
	opts := mdhtml.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: renderHook,
	}
	renderer := mdhtml.NewRenderer(opts)

	return string(markdown.Render(doc, renderer))
}

func renderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok && entering {
		language := string(code.Info)
		content := string(code.Literal)

		// Default to plain text if language not specified
		if language == "" {
			language = "text"
		}

		lexer := lexers.Get(language)
		if lexer == nil {
			lexer = lexers.Get("text")
		}

		// Use a light theme by default
		style := styles.Get("github")
		if style == nil {
			style = styles.Fallback
		}

		formatter := html.New(html.WithClasses(true))
		iterator, err := lexer.Tokenise(nil, content)
		if err != nil {
			// Fallback to plain text on error
			if _, err := io.WriteString(w, "<pre><code>"+content+"</code></pre>"); err != nil {
				return ast.GoToNext, true
			}
			return ast.GoToNext, true
		}

		buf := new(bytes.Buffer)
		err = formatter.Format(buf, style, iterator)
		if err != nil {
			if _, err := io.WriteString(w, "<pre><code>"+content+"</code></pre>"); err != nil {
				return ast.GoToNext, true
			}
			return ast.GoToNext, true
		}

		if _, err := io.WriteString(w, buf.String()); err != nil {
			return ast.GoToNext, true
		}
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}
