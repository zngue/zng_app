package stdhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/zngue/zng_app/core/http/core"
	"github.com/zngue/zng_app/core/http/server"
)

var colonParamRe = regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)

func toServeMuxPath(path string) string {
	return colonParamRe.ReplaceAllString(path, "{$1}")
}

type Server struct {
	mux      *http.ServeMux
	runtime  *server.Runtime
	prefix   string
	services []core.ServiceDesc
}

func NewServer(mux *http.ServeMux, runtime *server.Runtime, prefix string) *Server {
	prefix = "/" + strings.Trim(prefix, "/")
	if prefix != "/" {
		prefix += "/"
	}
	return &Server{mux: mux, runtime: runtime, prefix: prefix}
}

func (s *Server) RegisterService(serviceName string, register func(core.Router)) {
	collector := core.NewRouteCollector(serviceName)
	register(collector)
	desc := collector.ServiceDesc()
	s.services = append(s.services, desc)

	for _, route := range desc.Routes {
		currentRoute := route
		pattern := currentRoute.HTTPMethod + " " + s.prefix + strings.TrimLeft(toServeMuxPath(currentRoute.Path), "/")
		s.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			s.runtime.Handle(NewContext(w, r), currentRoute)
		})
	}
}

func (s *Server) Services() []core.ServiceDesc {
	services := make([]core.ServiceDesc, len(s.services))
	copy(services, s.services)
	return services
}

type Context struct {
	w       http.ResponseWriter
	r       *http.Request
	ctx     context.Context
	written bool
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w: w, r: r, ctx: r.Context()}
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *Context) Param(name string) string {
	return c.r.PathValue(name)
}

func (c *Context) Query(name string) string {
	return c.r.URL.Query().Get(name)
}

func (c *Context) QueryArray(name string) []string {
	return c.r.URL.Query()[name]
}

func (c *Context) Header(name string) string {
	return c.r.Header.Get(name)
}

func (c *Context) SetHeader(name string, value string) {
	c.w.Header().Set(name, value)
}

func (c *Context) BindJSON(obj any) error {
	return json.NewDecoder(c.r.Body).Decode(obj)
}

func (c *Context) JSON(statusCode int, obj any) {
	if c.written {
		return
	}
	c.written = true
	c.w.Header().Set("Content-Type", "application/json")
	c.w.WriteHeader(statusCode)
	_ = json.NewEncoder(c.w).Encode(obj)
}

func (c *Context) AbortWithJSON(statusCode int, obj any) {
	c.JSON(statusCode, obj)
}

func (c *Context) NewStream(ctx context.Context) core.Stream {
	c.written = true
	c.w.WriteHeader(http.StatusOK)
	flusher, _ := c.w.(http.Flusher)
	return &streamWriter{
		writer:  c.w,
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
