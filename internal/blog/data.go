package blog

import (
	"net/http"
	"time"

	. "github.com/Alan-Daniels/web/internal"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type BlogPost struct {
	PublishTime time.Time
	SafeName    string
	PostTitle   string
	HeroImage   string
	ReadTimeAvg string
	PostHead    func() templ.Component
	PostBody    func() templ.Component
}

func (post *BlogPost) Handler(e echo.Context) error {
	return Render(e, http.StatusOK, BlogPostView(post))
}
