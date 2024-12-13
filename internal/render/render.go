package render

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n float64) string {
	i := n / 100
	return fmt.Sprintf("$%.2f", i)
}

// Embed templates directory
//
//go:embed templates
var templateFS embed.FS

type Renderer struct {
	TemplateCache map[string]*template.Template
	Env           string
	StripeKey     string
	API           string
	ErrorLog      *log.Logger
}

func NewRenderer(env, stripeKey, api string, errorLog *log.Logger) *Renderer {
	return &Renderer{
		TemplateCache: make(map[string]*template.Template),
		Env:           env,
		StripeKey:     stripeKey,
		API:           api,
		ErrorLog:      errorLog,
	}
}

func (r *Renderer) AddDefaultData(td *TemplateData, req *http.Request) *TemplateData {
	td.StripePublic = r.StripeKey
	td.API = r.API
	return td
}

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
