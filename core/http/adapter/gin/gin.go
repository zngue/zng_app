package gin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/zngue/zng_app/core/http/core"
	"github.com/zngue/zng_app/core/http/server"
)

type Server struct {
	group    gin.IRoutes
	runtime  *server.Runtime
	services []core.ServiceDesc
}

func NewServer(group gin.IRoutes, runtime *server.Runtime) *Server {
	return &Server{group: group, runtime: runtime}
}

func (s *Server) RegisterService(serviceName string, register func(core.Router)) {
	collector := core.NewRouteCollector(serviceName)
	register(collector)
	desc := collector.ServiceDesc()
	s.services = append(s.services, desc)

	for _, route := range desc.Routes {
		currentRoute := route
		s.group.Handle(currentRoute.HTTPMethod, currentRoute.Path, func(ctx *gin.Context) {
			s.runtime.Handle(NewContext(ctx), currentRoute)
		})
	}
}

func (s *Server) Services() []core.ServiceDesc {
	services := make([]core.ServiceDesc, len(s.services))
	copy(services, s.services)
	return services
}

type Context struct {
	ctx *gin.Context
}

func NewContext(ctx *gin.Context) *Context {
	return &Context{ctx: ctx}
}

func (c *Context) Context() context.Context {
	return c.ctx.Request.Context()
}

func (c *Context) SetContext(ctx context.Context) {
	c.ctx.Request = c.ctx.Request.WithContext(ctx)
}

func (c *Context) Param(name string) string {
	return c.ctx.Param(name)
}

func (c *Context) Query(name string) string {
	return c.ctx.Query(name)
}

func (c *Context) QueryArray(name string) []string {
	return c.ctx.QueryArray(name)
}

func (c *Context) Header(name string) string {
	return c.ctx.GetHeader(name)
}

func (c *Context) SetHeader(name string, value string) {
	c.ctx.Header(name, value)
}

func (c *Context) BindJSON(obj any) error {
	return c.ctx.ShouldBindJSON(obj)
}

func (c *Context) JSON(statusCode int, obj any) {
	c.ctx.JSON(statusCode, obj)
}

func (c *Context) AbortWithJSON(statusCode int, obj any) {
	c.ctx.AbortWithStatusJSON(statusCode, obj)
}

func (c *Context) NewStream(ctx context.Context) core.Stream {
	c.ctx.Status(http.StatusOK)
	flusher, _ := c.ctx.Writer.(http.Flusher)
	return &streamWriter{
		writer:  c.ctx.Writer,
		flusher: flusher,
		ctx:     ctx,
	}
}

type streamWriter struct {
	writer  http.ResponseWriter
	flusher http.Flusher
	ctx     context.Context
}

func (s *streamWriter) Send(data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := s.writer.Write(append(payload, '\n')); err != nil {
		return err
	}
	if s.flusher != nil {
		s.flusher.Flush()
	}
	return nil
}

func (s *streamWriter) Close() error {
	return nil
}

func (s *streamWriter) Context() context.Context {
	return s.ctx
}
