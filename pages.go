package page

import (
	"distudio.com/mage"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
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
	mage.Controller
}

func (page *StaticPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	fname := fmt.Sprintf("%s.html", page.FileName)
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		log.Errorf(ctx, "Can't find file %s", fname)
		return mage.Redirect{Status: http.StatusNotFound}
	}

	str, err := ioutil.ReadFile(fname)

	if err != nil {
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	renderer := mage.TextRenderer{}
	renderer.Data = string(str)
	out.Renderer = &renderer

	return mage.Redirect{Status: http.StatusOK}
}

func (page *StaticPage) OnDestroy(ctx context.Context) {

}

/**
returns a 404 page with static page
*/
type FourOFourPage struct {
	StaticPage
}

func (page *FourOFourPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	if page.FileName != "" {
		redir := page.StaticPage.Process(ctx, out)
		out.AddHeader("Content-type", "text/html; charset=utf-8")
		switch redir.Status {
		case http.StatusOK:
			return mage.Redirect{Status: http.StatusNotFound}
		case http.StatusInternalServerError:
			return redir
		}
	}

	return mage.Redirect{Status: http.StatusNotFound}
}

/**
returns a 404 page with the given template
*/
type StatusTemplatedPage struct {
	TemplatedPage
	Status int
}

func (page *StatusTemplatedPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	if page.FileName != "" {
		redir := page.TemplatedPage.Process(ctx, out)
		out.AddHeader("Content-type", "text/html; charset=utf-8")
		switch redir.Status {
		case http.StatusOK:
			return mage.Redirect{Status: page.Status}
		case http.StatusInternalServerError:
			return redir
		}
	}

	return mage.Redirect{Status: page.Status}
}

/**
returns a 404 page with the given localized template
*/
type LocalizedStatusPage struct {
	LocalizedPage
	Status int
}

func (page *LocalizedStatusPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	if page.FileName != "" {
		redir := page.LocalizedPage.Process(ctx, out)
		out.AddHeader("Content-type", "text/html; charset=utf-8")
		switch redir.Status {
		case http.StatusOK:
			return mage.Redirect{Status: page.Status}
		case http.StatusInternalServerError:
			return redir
		}
	}

	return mage.Redirect{Status: page.Status}
}

//Reads a template and mixes it with a base template (useful for headers/footers)
//Base is the name of the base template if any
type TemplatedPage struct {
	Url         string
	FileName    string
	BaseName    string
	Bases       []string
	DataHandler TemplateDataHandler
	mage.Controller
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

func (page *TemplatedPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	fname := fmt.Sprintf("%s.html", page.FileName)
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find file %s", fname)
		return mage.Redirect{Status: http.StatusNotFound}
	}

	files := make([]string, 0, 0)
	files = append(files, page.Bases...)
	files = append(files, fname)

	tpl, err := template.ParseFiles(files...)

	if err != nil {
		log.Errorf(ctx, "Cant' parse template files: %v", err)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	renderer := mage.TemplateRenderer{}
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

	return mage.Redirect{Status: http.StatusOK}
}

func (page *TemplatedPage) OnDestroy(ctx context.Context) {

}

//Has a TemplatedPage. Attaches to each templated page a corresponding json file that specifies translations
type LocalizedPage struct {
	TemplatedPage
	JsonBaseFile string
	JsonFile     string
	Locale       string
}

func (page *LocalizedPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	fname := fmt.Sprintf("%s.html", page.FileName)
	_, err := os.Stat(fname)

	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find file %s", fname)
		return mage.Redirect{Status: http.StatusNotFound}
	}

	//get the language hint
	inputs := mage.InputsFromContext(ctx)

	lang := page.Locale

	_, lok := inputs["X-AppEngine-Country"]

	if lok {
		lang = inputs["X-AppEngine-Country"].Value()
	}

	_, lok = inputs[LanguageCookieKey]
	if lok {
		lang = inputs[LanguageCookieKey].Value()
	}

	_, lok = inputs["hl"]

	if lok {
		lang = inputs["hl"].Value()
	}

	lcookie := http.Cookie{}
	lcookie.Name = LanguageCookieKey
	lcookie.Value = lang
	lcookie.Path = "/"
	out.AddCookie(lcookie)

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

	files := make([]string, 0, 0)
	files = append(files, page.Bases...)
	files = append(files, fname)

	tpl, err := template.New("").Funcs(funcMap).ParseFiles(files...)

	if err != nil {
		log.Errorf(ctx, "Cant' parse template files: %v", err)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	//get the base language file
	lbasename := page.JsonBaseFile

	if lbasename == "" {
		lbasename = fmt.Sprintf("i18n/%s", "base.json")
	}

	jbase, err := ioutil.ReadFile(lbasename)

	if err != nil {
		log.Errorf(ctx, "Error reading base language file %s: %v", lbasename, err)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	var base map[string]interface{}
	err = json.Unmarshal(jbase, &base)

	if err != nil {
		log.Errorf(ctx, "Invalid json for base file %s: %v", lbasename, err)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	_, bok := base[lang]

	if !bok {
		log.Errorf(ctx, "Base language file %s doesn't support language %s", lbasename, lang)
		//we get the default value if the user provides an invalid lang
		lang = page.Locale
	}

	globals := base[lang]

	//---- get the specific language json file
	lfname := page.JsonFile

	if lfname == "" {
		lfname = fmt.Sprintf("i18n/%s.%s", page.FileName, "json")
	}

	//now that we have the locale, read the json language file and get the corresponding values
	jlang, err := ioutil.ReadFile(lfname)

	if err != nil {
		log.Errorf(ctx, "Error retrieving language file %s: %v", lfname, err)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	var contents map[string]interface{}
	err = json.Unmarshal(jlang, &contents)

	if err != nil {
		log.Errorf(ctx, "Invalid json for file %s: %v", lfname, err)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	_, dok := contents[lang]

	if !dok {
		log.Errorf(ctx, "File %s doesn't support language %s", lfname, lang)
		return mage.Redirect{Status: http.StatusInternalServerError}
	}

	content := contents[lang]

	renderer := mage.TemplateRenderer{}
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

	return mage.Redirect{Status: http.StatusOK}
}

//sends an email with the specified message and sender
type SendMailPage struct {
	mage.Controller
	Mailer
}

type Mailer interface {
	ValidateAndSend(ctx context.Context, inputs mage.RequestInputs) error
}

func (page *SendMailPage) Process(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {

	inputs := mage.InputsFromContext(ctx)

	method := inputs[mage.KeyRequestMethod].Value()

	if method != http.MethodPost {
		return mage.Redirect{Status: http.StatusMethodNotAllowed}
	}

	err := page.Mailer.ValidateAndSend(ctx, inputs)

	if err != nil {
		//if we have a field error we handle it returning a 404
		renderer := mage.JSONRenderer{}
		renderer.Data = err.Error()
		out.Renderer = &renderer
		return mage.Redirect{Status: http.StatusBadRequest}
	}

	return mage.Redirect{Status: http.StatusOK}
}

func (page *SendMailPage) OnDestroy(ctx context.Context) {

}
