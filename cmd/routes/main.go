// Command routes prints every registered HTTP route, like Laravel's `route:list`.
//
//	go run ./cmd/routes
//	go run ./cmd/routes -method POST
//	go run ./cmd/routes -path auth
package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/gofiber/fiber/v3"

	"tripo-backend/internal/config"
	"tripo-backend/internal/routes"
)

func main() {
	method := flag.String("method", "", "only show routes with this HTTP method")
	path := flag.String("path", "", "only show routes whose path contains this text")
	flag.Parse()

	// Handlers are only registered, never called, so no database or real secret is needed.
	app := fiber.New()
	routes.Setup(app, nil, &config.Config{})

	// Group middleware is registered via Use, which GetRoutes(true) hides. Comparing
	// the two listings recovers it; it applies to every route under its path prefix.
	real := app.GetRoutes(true)
	isReal := make(map[string]bool, len(real))
	for _, r := range real {
		isReal[routeKey(r)] = true
	}
	var uses []row
	for _, r := range app.GetRoutes(false) {
		if !isReal[routeKey(r)] && r.Method == fiber.MethodGet { // Use is expanded once per method
			uses = append(uses, newRow(r))
		}
	}

	var rows []row
	for _, r := range real {
		if r.Method == fiber.MethodHead { // Fiber auto-adds HEAD for every GET
			continue
		}
		if *method != "" && !strings.EqualFold(r.Method, *method) {
			continue
		}
		if *path != "" && !strings.Contains(r.Path, *path) {
			continue
		}
		row := newRow(r)
		for _, u := range uses {
			if r.Path == u.path || strings.HasPrefix(r.Path, u.path+"/") {
				row.middleware = "group" // Fiber wraps Use handlers, so their real name is not recoverable
			}
		}
		rows = append(rows, row)
	}

	sort.SliceStable(rows, func(i, j int) bool { return rows[i].path < rows[j].path })

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "METHOD\tURI\tHANDLER\tMIDDLEWARE")
	for _, r := range rows {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.method, r.path, r.handler, r.middleware)
	}
	w.Flush()
	fmt.Printf("\nShowing %d routes\n", len(rows))
}

func routeKey(r fiber.Route) string {
	h := "-"
	if n := len(r.Handlers); n > 0 {
		h = funcName(r.Handlers[n-1])
	}
	return r.Method + " " + r.Path + " " + h
}

type row struct{ method, path, handler, middleware string }

func newRow(r fiber.Route) row {
	names := make([]string, len(r.Handlers))
	for i, h := range r.Handlers {
		names[i] = funcName(h)
	}

	out := row{method: r.Method, path: r.Path, handler: "-", middleware: "-"}
	if n := len(names); n > 0 {
		out.handler = names[n-1]
		if n > 1 {
			out.middleware = strings.Join(names[:n-1], ", ")
		}
	}
	return out
}

// funcName turns "tripo-backend/internal/controllers.(*AuthController).Login-fm"
// into "AuthController@Login".
func funcName(h any) string {
	name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	name = name[strings.LastIndex(name, "/")+1:] // drop module path
	name = strings.TrimSuffix(name, "-fm")
	name = strings.NewReplacer("(*", "", ")", "").Replace(name)

	parts := strings.Split(name, ".") // pkg, Type, Method  |  pkg, Func, funcN
	switch {
	case len(parts) == 3 && !strings.HasPrefix(parts[2], "func"):
		return parts[1] + "@" + parts[2]
	case len(parts) >= 2:
		return parts[0] + "." + parts[1] // closures, e.g. middleware.Authenticate
	}
	return name
}
