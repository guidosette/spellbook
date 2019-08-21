package spellbook

import (
	"context"
	"decodica.com/flamel"
	"encoding/json"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/appengine/log"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	LanguageCookieKey string = "PAGE_LANG_CURRENT_ID"
)

//Reads a static file and outputs it as a string.
//It is usually used to print static html pages.
//If a template is needed use TemplatedPage instead
type StaticPage struct {
	FileName string
	flamel.Controller
}

func (page *StaticPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	fname := fmt.Sprintf("%s.html", page.FileName)
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		log.Errorf(ctx, "Can't find file %s", fname)
		return flamel.HttpResponse{Status: http.StatusNotFound}
	}

	str, err := ioutil.ReadFile(fname)

	if err != nil {
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	renderer := flamel.TextRenderer{}
	renderer.Data = string(str)
	out.Renderer = &renderer

	return flamel.HttpResponse{Status: http.StatusOK}
}

func (page *StaticPage) OnDestroy(ctx context.Context) {

}

/**
returns a 404 page with static page
*/
type FourOFourPage struct {
	StaticPage
}

func (page *FourOFourPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	if page.FileName != "" {
		redir := page.StaticPage.Process(ctx, out)
		out.AddHeader("Content-type", "text/html; charset=utf-8")
		switch redir.Status {
		case http.StatusOK:
			return flamel.HttpResponse{Status: http.StatusNotFound}
		case http.StatusInternalServerError:
			return redir
		}
	}

	return flamel.HttpResponse{Status: http.StatusNotFound}
}

/**
returns a 404 page with the given template
*/
type StatusTemplatedPage struct {
	TemplatedPage
	Status int
}

func (page *StatusTemplatedPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	if page.FileName != "" {
		redir := page.TemplatedPage.Process(ctx, out)
		out.AddHeader("Content-type", "text/html; charset=utf-8")
		switch redir.Status {
		case http.StatusOK:
			return flamel.HttpResponse{Status: page.Status}
		case http.StatusInternalServerError:
			return redir
		}
	}

	return flamel.HttpResponse{Status: page.Status}
}

/**
returns a 404 page with the given localized template
*/
type LocalizedStatusPage struct {
	LocalizedPage
	Status int
}

func (page *LocalizedStatusPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	if page.FileName != "" {
		redir := page.LocalizedPage.Process(ctx, out)
		out.AddHeader("Content-type", "text/html; charset=utf-8")
		switch redir.Status {
		case http.StatusOK:
			return flamel.HttpResponse{Status: page.Status}
		case http.StatusInternalServerError:
			return redir
		}
	}

	return flamel.HttpResponse{Status: page.Status}
}

//Reads a template and mixes it with a base template (useful for headers/footers)
//Base is the name of the base template if any
type TemplatedPage struct {
	Url         string
	FileName    string
	BaseName    string
	Bases       []string
	DataHandler TemplateDataHandler
	flamel.Controller
	FuncHandler TemplateFuncHandler
}

type TemplateDataHandler interface {
	AssignData(ctx context.Context) interface{}
}

type TemplateFuncHandler interface {
	AssignFuncMap(ctx context.Context) template.FuncMap
}

func NewTemplatedPage(url string, filename string, bases ...string) TemplatedPage {
	page := TemplatedPage{}
	page.Url = url
	page.FileName = filename
	page.Bases = bases
	return page
}

func (page *TemplatedPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	fname := fmt.Sprintf("%s.html", page.FileName)
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find file %s", fname)
		return flamel.HttpResponse{Status: http.StatusNotFound}
	}

	files := make([]string, 0, 0)
	files = append(files, page.Bases...)
	files = append(files, fname)

	tpl, err := template.ParseFiles(files...)

	if err != nil {
		log.Errorf(ctx, "Cant' parse template files: %v", err)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	renderer := flamel.TemplateRenderer{}
	if page.BaseName == "" {
		renderer.TemplateName = "base"
	} else {
		renderer.TemplateName = page.BaseName
	}

	renderer.Template = tpl
	if page.DataHandler != nil {
		renderer.Data = page.DataHandler.AssignData(ctx)
	}

	out.Renderer = &renderer

	return flamel.HttpResponse{Status: http.StatusOK}
}

func (page *TemplatedPage) OnDestroy(ctx context.Context) {

}

//Has a Templatedspellbook. Attaches to each templated page a corresponding json file that specifies translations
type LocalizedPage struct {
	TemplatedPage
	JsonBaseFile string
	JsonFile     string
	Language string
}

func (page *LocalizedPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	fname := fmt.Sprintf("%s.html", page.FileName)
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find file %s", fname)
		return flamel.HttpResponse{Status: http.StatusNotFound}
	}

	lang := page.Language
	if lang == "" {
		t := ctx.Value(KeyLanguageTag)
		tag := t.(language.Tag)
		lang = tag.String()
	}

	//create the link creator function
	funcMap := template.FuncMap{
		"LocalizedUrl": func(url string) string {
			return fmt.Sprintf("%s?hl=%s", url, lang)
		},
		"ToJson": func(data interface{}) template.HTML {
			j, _ := json.Marshal(data)
			return template.HTML(j)
		},
	}
	if page.FuncHandler != nil {
		customFuncMap := page.FuncHandler.AssignFuncMap(ctx)
		for k, v := range customFuncMap {
			funcMap[k] = v
		}
	}

	files := make([]string, 0, 0)
	files = append(files, page.Bases...)
	files = append(files, fname)

	tpl, err := template.New("").Funcs(funcMap).ParseFiles(files...)

	if err != nil {
		log.Errorf(ctx, "Cant' parse template files: %v", err)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	//get the base language file
	lbasename := page.JsonBaseFile

	if lbasename == "" {
		lbasename = fmt.Sprintf("i18n/%s", "base.json")
	}

	jbase, err := ioutil.ReadFile(lbasename)

	if err != nil {
		log.Errorf(ctx, "Error reading base language file %s: %v", lbasename, err)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	var base map[string]interface{}
	err = json.Unmarshal(jbase, &base)

	if err != nil {
		log.Errorf(ctx, "Invalid json for base file %s: %v", lbasename, err)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	_, bok := base[lang]

	if !bok {
		err := fmt.Errorf("base language file %s doesn't support language %s", lbasename, lang)
		log.Errorf(ctx, err.Error())
		//we get the default value if the user provides an invalid lang
		panic(err)
	}

	globals := base[lang]

	//---- get the specific language json file
	lfname := page.JsonFile

	if lfname == "" {
		lfname = fmt.Sprintf("i18n/%s.%s", page.FileName, "json")
	}

	_, err = os.Stat(lfname)
	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find json file %s", fname)
		return flamel.HttpResponse{Status: http.StatusNotFound}
	}

	//now that we have the locale, read the json language file and get the corresponding values
	jlang, err := ioutil.ReadFile(lfname)

	if err != nil {
		log.Errorf(ctx, "Error retrieving language file %s: %v", lfname, err)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	var contents map[string]interface{}
	err = json.Unmarshal(jlang, &contents)

	if err != nil {
		log.Errorf(ctx, "Invalid json for file %s: %v", lfname, err)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	_, dok := contents[lang]

	if !dok {
		log.Errorf(ctx, "File %s doesn't support language %s", lfname, lang)
		return flamel.HttpResponse{Status: http.StatusInternalServerError}
	}

	content := contents[lang]

	renderer := flamel.TemplateRenderer{}
	if page.BaseName == "" {
		renderer.TemplateName = "base"
	} else {
		renderer.TemplateName = page.BaseName
	}

	renderer.Template = tpl

	var data interface{}
	if page.DataHandler != nil {
		data = page.DataHandler.AssignData(ctx)
	}

	renderer.Data = struct {
		Url      string
		Language string
		Globals  interface{}
		Content  interface{}
		Data     interface{}
	}{
		page.Url,
		lang,
		globals,
		content,
		data,
	}

	out.Renderer = &renderer

	return flamel.HttpResponse{Status: http.StatusOK}
}

//sends an email with the specified message and sender
type SendMailPage struct {
	flamel.Controller
	Mailer
}

type Mailer interface {
	ValidateAndSend(ctx context.Context, inputs flamel.RequestInputs) error
}

func (page *SendMailPage) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {

	inputs := flamel.InputsFromContext(ctx)

	method := inputs[flamel.KeyRequestMethod].Value()

	if method != http.MethodPost {
		return flamel.HttpResponse{Status: http.StatusMethodNotAllowed}
	}

	err := page.Mailer.ValidateAndSend(ctx, inputs)

	if err != nil {
		//if we have a field error we handle it returning a 404
		renderer := flamel.JSONRenderer{}
		renderer.Data = err.Error()
		out.Renderer = &renderer
		return flamel.HttpResponse{Status: http.StatusBadRequest}
	}

	return flamel.HttpResponse{Status: http.StatusOK}
}

func (page *SendMailPage) OnDestroy(ctx context.Context) {

}

/**
REDIRECT StatusMovedPermanently
*/
type MovedController struct {
	flamel.Controller
	To string
}

func (controller *MovedController) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	return flamel.HttpResponse{Location: controller.To, Status: http.StatusMovedPermanently}
}

func (controller *MovedController) OnDestroy(ctx context.Context) {}

/**
REDIRECT StatusFound
*/
type FoundController struct {
	flamel.Controller
	To string
}

func (controller *FoundController) Process(ctx context.Context, out *flamel.ResponseOutput) flamel.HttpResponse {
	return flamel.HttpResponse{Location: controller.To, Status: http.StatusFound}
}

func (controller *FoundController) OnDestroy(ctx context.Context) {}
