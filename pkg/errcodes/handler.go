package errcodes

import (
	"github.com/dmlittle/discoverrewind/pkg/errutils"
	"github.com/dmlittle/discoverrewind/pkg/logger"
	"github.com/foolin/goview"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Handle is an Echo error handler that uses HTTP errors accordingly, and any
// generic error will be interpreted as an internal server error.
func (h *Handler) Handle(err error, c echo.Context) {
	if errutils.IsIgnorableErr(err) {
		logger.FromEchoContext(c).Err(err).Warn("broken pipe")
		return
	}

	httpCode, msg := h.generatePayload(err)

	// Internal server errors
	if httpCode == http.StatusInternalServerError {
		logger.FromEchoContext(c).Err(err).Error("server error")
	}

	if err := c.Render(httpCode, "error", goview.M{"code": httpCode, "message": msg}); err != nil {
		logger.FromEchoContext(c).Err(errors.WithStack(err)).Error("error handler json error")
	}
}

func (h *Handler) generatePayload(err error) (int, string) {
	msg := ""
	httpCode := http.StatusInternalServerError

	// Echo errors
	var he *echo.HTTPError
	if ok := errors.As(err, &he); ok {
		httpCode = he.Code
		msg = he.Message.(string)
	}

	// Custom errors
	var e *Error
	if ok := errors.As(err, &e); ok {
		httpCode = e.HTTPCode
		msg = e.Message
	}

	// Internal server errors that aren't Echo errors or custom errors
	if httpCode == http.StatusInternalServerError && msg == "" {
		msg = "Internal Server Error"
	}

	return httpCode, msg
}
