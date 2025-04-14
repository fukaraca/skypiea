package templater

import (
	"html/template"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
)

type Templates map[string]*render.HTMLProduction

func New() Templates {
	return make(Templates)
}

func (t Templates) Instance(name string, data any) render.Render {
	if v, ok := t[name]; ok {
		return v.Instance(name, data)
	}
	return nil
}

func (t Templates) LoadHTMLGlob(pattern string) {
	pattern, err := filepath.Abs(pattern)
	if err != nil {
		panic(err.Error())
	}
	pages, err := filepath.Glob(filepath.Join(pattern, "pages/*"))
	if err != nil {
		panic(err.Error())
	}
	layouts, err := filepath.Glob(filepath.Join(pattern, "layouts/*"))
	if err != nil {
		panic(err.Error())
	}
	for _, page := range pages {
		filename := filepath.Base(page)
		filename, _, _ = strings.Cut(filename, ".")
		filenames := make([]string, len(layouts))
		copy(filenames, layouts)
		filenames = append(filenames, page) //nolint: makezero
		temp := template.Must(template.New(filename).ParseFiles(filenames...))
		t[filename] = &render.HTMLProduction{
			Template: temp,
		}
	}
}
