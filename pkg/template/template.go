package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"reflect"
)

type manager struct {
	dir       string
	data      map[string]interface{}
	templates map[string]*template.Template
}

type IManager interface {
	LoadTemplates() (IManager, error)
	Render(w http.ResponseWriter, status int, template string, data map[string]interface{}) error
	RenderRaw(w http.ResponseWriter, status int, content interface{}) error
	RenderJson(w http.ResponseWriter, status int, value interface{}) error
	AddData(key string, value interface{}) IManager
	InjectData(key string, value interface{}) IManager
}

func NewTemplates(dir string, defaultData map[string]interface{}) IManager {
	return &manager{
		dir:       dir,
		data:      defaultData,
		templates: make(map[string]*template.Template, 0),
	}
}

func (t *manager) LoadTemplates() (IManager, error) {
	layouts, err := filepath.Glob(fmt.Sprintf("%s/layouts/*.html", t.dir))
	if err != nil {
		return nil, errors.Unwrap(fmt.Errorf("error %w when parsing views template at: %s/layouts/*.html", err, t.dir))
	}
	views, err2 := filepath.Glob(fmt.Sprintf("%s/views/*.html", t.dir))
	if err2 != nil {
		return nil, errors.Unwrap(fmt.Errorf("error %w when parsing views template at: %s/views/*.html", err, t.dir))
	}
	for _, view := range views {
		files := append(layouts, view)
		t.templates[filepath.Base(view)] = template.Must(template.ParseFiles(files...))
	}
	return t, nil
}

func (t *manager) render(w http.ResponseWriter, status int, tpl string, data map[string]interface{}, f interface{}) error {
	tmpl, err := t.getTemplate(tpl)
	if err != nil {
		return err
	}
	if f != nil {
		if ff, ok := f.(template.FuncMap); ok {
			tmpl = tmpl.Funcs(ff)
		}
	}
	_ = t.mergeData(data)
	// Render the template 'name' with data
	w.WriteHeader(status)
	if err = tmpl.ExecuteTemplate(w, tpl, t.getData()); err != nil {
		return err
	}
	return nil
}

func (t *manager) Render(w http.ResponseWriter, status int, tpl string, data map[string]interface{}) error {
	return t.render(w, status, tpl, data, nil)
}

func (t *manager) RenderRaw(w http.ResponseWriter, status int, content interface{}) error {
	w.WriteHeader(status)
	_, err := fmt.Fprintln(w, content)
	return err
}

func (t *manager) RenderJson(w http.ResponseWriter, status int, value interface{}) error {
	v, err := json.Marshal(value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err == nil {
		_, err = fmt.Fprintln(w, string(v))
		return err
	}
	_, err = fmt.Fprintln(w, "{}")
	return err
}

func (t *manager) getTemplate(name string) (*template.Template, error) {
	if v, ok := t.templates[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("template %s not found", name)
}

func (t *manager) AddData(key string, value interface{}) IManager {
	if _, ok := t.data[key]; !ok {
		t.data[key] = value
	}
	return t
}

func (t *manager) InjectData(key string, value interface{}) IManager {
	t.data[key] = value
	return t
}

func (t *manager) getData() map[string]interface{} {
	return t.data
}

func (t *manager) mergeData(data map[string]interface{}) error {
	if len(data) < 1 {
		return fmt.Errorf("data arguments are empty")
	}
	t.data = t.merge(t.data, data).(map[string]interface{})
	return nil
}

func (t *manager) merge(dst interface{}, src interface{}) (rs interface{}) {
	dstValue := reflect.TypeOf(dst)
	dstType := dstValue.Kind()
	srcValue := reflect.TypeOf(src)
	srcType := srcValue.Kind()
	// when different then replace the value
	if srcType != dstType {
		return src
	}
	switch dstType {
	case reflect.Map:
		if dstValue.Elem().Kind() != srcValue.Elem().Kind() {
			return src
		}
		dstMap := dst.(map[string]interface{})
		srcMap := src.(map[string]interface{})
	srcLoop:
		for k, v := range srcMap {
			for kk, vv := range dstMap {
				if kk == k {
					dstMap[kk] = t.merge(vv, v)
					continue srcLoop
				}
			}
			dstMap[k] = v
		}
		rs = dstMap
	default:
		rs = src
	}
	return
}
