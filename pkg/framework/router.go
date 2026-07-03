package framework

import (
	"errors"
	"net/http"
	"strings"
)

var errNotFound = errors.New("not found")

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	params   map[string]string
}

type Handler func(*Context) (Page, error)

type Router struct {
	routes []route
}

type route struct {
	method  string
	pattern string
	parts   []string
	page    Handler
}

func NewRouter() *Router {
	return &Router{
		routes: []route{},
	}
}

func (r *Router) Page(method, pattern string, handler Handler) {
	r.routes = append(r.routes, route{method: method, pattern: pattern, parts: splitPattern(pattern), page: handler})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		params, ok := matchRoute(route, req.Method, req.URL.Path)
		if !ok {
			continue
		}

		ctx := &Context{Response: w, Request: req, params: params}
		page, err := route.page(ctx)
		if err != nil {
			r.handleError(w, err)
			return
		}
		page.render(w)
		return
	}
	r.handleError(w, errNotFound)
}

func (r *Router) handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, errNotFound) {
		Page{Status: http.StatusNotFound, Title: "Not found", Template: "error.html", Data: "The page was not found."}.render(w)
		return
	}
	Page{Status: http.StatusInternalServerError, Title: "Server error", Template: "error.html", Data: err.Error()}.render(w)
}

func splitPattern(pattern string) []string {
	trimmed := strings.Trim(pattern, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func matchRoute(route route, method, path string) (map[string]string, bool) {
	if route.method != method {
		return nil, false
	}
	parts := splitPattern(path)
	if len(parts) != len(route.parts) {
		return nil, false
	}

	params := map[string]string{}
	for i, routePart := range route.parts {
		pathPart := parts[i]
		if strings.HasPrefix(routePart, "{") && strings.HasSuffix(routePart, "}") {
			params[strings.TrimSuffix(strings.TrimPrefix(routePart, "{"), "}")] = pathPart
			continue
		}
		if routePart != pathPart {
			return nil, false
		}
	}
	return params, true
}
