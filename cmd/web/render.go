package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type templateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]any
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated int
	API             string
	CSSVersion      string
	StripePublic    string
}

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n float64) string {
	i := n / 100
	return fmt.Sprintf("$%.2f", i)
}

//go:embed templates
var templateFS embed.FS

func (app *application) AddDefaultData(td *templateData, r *http.Request) *templateData {
	td.StripePublic = app.config.stripe.key
	td.API = app.config.api
	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.html", page)

	_, templateInMap := app.templateCache[templateToRender]

	if app.config.env == "production" && templateInMap {
		t = app.templateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = app.AddDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
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
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t

	return t, nil
}
