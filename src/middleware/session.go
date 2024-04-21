package middleware

import (
	"dev.chaiyapluek.cloud.final.frontend/src/service"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type DefaultSessionMiddleware struct {
	serssionService service.SessionService
}

func NewDefaultSessionMiddleware(sessionService service.SessionService) *DefaultSessionMiddleware {
	return &DefaultSessionMiddleware{
		serssionService: sessionService,
	}
}

func (m *DefaultSessionMiddleware) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("sessionid", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400,
			HttpOnly: true,
		}
		sessionId, ok := sess.Values["id"].(string)

		//isLogin := c.Get("isLogin").(bool)

		if !ok {
			newSessionId, _ := m.serssionService.NewSessionDetail()
			sess.Values["id"] = newSessionId
		} else {
			detail := m.serssionService.GetSessionDetail(sessionId)
			if detail == nil {
				newSessionId, _ := m.serssionService.NewSessionDetail()
				sess.Values["id"] = newSessionId
			}
		}
		sess.Save(c.Request(), c.Response())
		return next(c)
	}
}
