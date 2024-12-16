package render

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// functions defines custom template functions.
var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
	"concat":         concat,
}

// formatCurrency formats a float as a currency string.
func formatCurrency(n float64) string {
	i := n / 100
	return fmt.Sprintf("$%.2f", i)
}

// concat Concat two strings
func concat(x, y, sep string) string {
	if sep == "" {
		sep = " "
	}
	return x + sep + y
}

// Embed templates directory
//
//go:embed templates
var templateFS embed.FS

// Renderer manages template rendering and caching.
type Renderer struct {
	TemplateCache map[string]*template.Template
	Env           string
	StripeKey     string
	API           string
	ErrorLog      *log.Logger
}

// NewRenderer initializes a Renderer with caching and configuration.
func NewRenderer(env, stripeKey, api string, errorLog *log.Logger) *Renderer {
	return &Renderer{
		TemplateCache: make(map[string]*template.Template),
		Env:           env,
		StripeKey:     stripeKey,
		API:           api,
		ErrorLog:      errorLog,
	}
}

// AddDefaultData adds default data like Stripe key and API URL to templates.
func (r *Renderer) AddDefaultData(td *TemplateData, req *http.Request) *TemplateData {
	td.StripePublic = r.StripeKey
	td.API = r.API
	return td
}

// RenderTemplate renders a template with the provided data and optional partials.
func (r *Renderer) RenderTemplate(w http.ResponseWriter, req *http.Request, page string, td *TemplateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.html", page)

	_, templateInMap := r.TemplateCache[templateToRender]

	if r.Env == "production" && templateInMap {
		t = r.TemplateCache[templateToRender]
	} else {
		t, err = r.parseTemplate(partials, page, templateToRender)
		if err != nil {
			r.ErrorLog.Println("Error parsing template:", err)
			return err
		}
	}

	if td == nil {
		td = &TemplateData{}
	}

	td = r.AddDefaultData(td, req)

	err = t.Execute(w, td)
	if err != nil {
		r.ErrorLog.Println("Error executing template:", err)
		return err
	}

	return nil
}

// parseTemplate parses and caches a template with optional partials.
func (r *Renderer) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error
	baseTemplate := "templates/base.layout.html"

	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partials.html", x)
		}

		templateName := fmt.Sprintf("%s.page.html", page)
		partialsList := strings.Join(partials, ",")
		t, err = template.New(templateName).
			Funcs(functions).
			ParseFS(templateFS, baseTemplate, partialsList, templateToRender)
	} else {
		templateName := fmt.Sprintf("%s.page.html", page)
		t, err = template.New(templateName).
			Funcs(functions).
			ParseFS(templateFS, baseTemplate, templateToRender)
	}

	if err != nil {
		r.ErrorLog.Println("Error parsing template:", err)
		return nil, err
	}

	r.TemplateCache[templateToRender] = t

	return t, nil
}
