package models

import (
	"encoding/json"

	"github.com/captain-corp/captain/config"
	"gorm.io/gorm"
)

// Settings represents the site configuration
type Settings struct {
	gorm.Model
	Title        string `gorm:"not null" form:"title"`
	Subtitle     string `gorm:"not null" form:"subtitle"`
	ChromaStyle  string `gorm:"not null" form:"chroma_style"`
	Theme        string `gorm:"not null" form:"theme"`
	PostsPerPage int    `gorm:"not null" form:"posts_per_page"`
	LogoID       *uint  `gorm:"" form:"logo_id"`
	UseFavicon   bool   `gorm:"not null;default:false" form:"use_favicon"`
}

func (s *Settings) ToJSON() string {

	chromaStyleList := make([]map[string]string, len(config.GetChromaStyles()))
	for i, style := range config.GetChromaStyles() {
		chromaStyleList[i] = map[string]string{
			"name":  style,
			"value": style,
		}
	}

	b, _ := json.Marshal(map[string]interface{}{
		"title":              s.Title,
		"subtitle":           s.Subtitle,
		"codeHighlightTheme": s.ChromaStyle,
		"codeThemeOptions":   chromaStyleList,
		"theme":              s.Theme,
		"postsPerPage":       s.PostsPerPage,
		"logoId":             s.LogoID,
		"useFavicon":         s.UseFavicon,
	})
	return string(b)
}
