package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// FiberContextToHTTP converts a Fiber context to standard net/http types
func FiberContextToHTTP(c *fiber.Ctx) (*http.Request, http.ResponseWriter) {
	// Convert fasthttp request to net/http request
	var req http.Request
	fasthttpadaptor.ConvertRequest(c.Context(), &req, true)
	
	// Create response writer wrapper
	w := &fiberResponseWriter{c}
	return &req, w
}

type fiberResponseWriter struct {
	c *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	h := make(http.Header)
	w.c.Response().Header.VisitAll(func(k, v []byte) {
		h.Add(utils.UnsafeString(k), utils.UnsafeString(v))
	})
	return h
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	return w.c.Write(b)
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.c.Status(statusCode)
} 