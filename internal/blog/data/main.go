package data

import (
	"fmt"
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

func (post *BlogPost) Handler() func(e echo.Context) error {
	return func(e echo.Context) error {
		_, err := InitSession(e)
		if err != nil {
			return err
		}
		return Render(e, http.StatusOK, blogPostView(post))
	}
}

func (post *BlogPost) Url() templ.SafeURL {
	return templ.SafeURL(fmt.Sprintf("/blog/%s", post.SafeName))
}

func IndexHandler(posts []BlogPost) func(e echo.Context) error {
	return func(e echo.Context) error {
		_, err := InitSession(e)
		if err != nil {
			return err
		}
		return Render(e, http.StatusOK, blogRoot(posts))
	}
}
