package page

import (
	"distudio.com/mage"
	"golang.org/x/net/context"
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	netmail "net/mail"
	"google.golang.org/appengine/log"
	"html/template"
	"google.golang.org/appengine/mail"
)

//Reads a static file and outputs it as a string.
//It is usually used to print static html pages.
//If a template is needed use TemplatedPage instead
type StaticPage struct {
	FileName string
	mage.Page
}

func (page *StaticPage) Process(ctx context.Context, out *mage.RequestOutput) mage.Redirect {
	fname := fmt.Sprintf("%s.html", page.FileName);
	_, err := os.Stat(fname);

	if os.IsNotExist(err) {
		log.Errorf(ctx, "Can't find file %s", fname);
		return mage.Redirect{Status:http.StatusNotFound}
	}

	str, err := ioutil.ReadFile(fname);

	if err != nil {
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	renderer := mage.TextRenderer{};
	renderer.Data = string(str);
	out.Renderer = &renderer;

	return mage.Redirect{Status:http.StatusOK};
}

func (page *StaticPage) OnDestroy(ctx context.Context) {

}

//Reads a template and mixes it with a base template (useful for headers/footers)
//Base is the name of the base template if any
type TemplatedPage struct {
	Url string
	FileName string
	Bases []string
	mage.Page
}

func NewTemplatedPage(url string, filename string, bases ...string) TemplatedPage {
	page := TemplatedPage{};
	page.Url = url;
	page.FileName = filename;
	page.Bases = bases;
	return page;
}


func (page *TemplatedPage) Process(ctx context.Context, out *mage.RequestOutput) mage.Redirect {
	fname := fmt.Sprintf("%s.html", page.FileName);
	_, err := os.Stat(fname);

	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find file %s", fname);
		return mage.Redirect{Status:http.StatusNotFound}
	}

	files := make([]string, 0, 0);
	files = append(files, page.Bases...);
	files = append(files, fname);

	tpl, err := template.ParseFiles(files...);

	if err != nil {
		log.Errorf(ctx, "Cant' parse template files: %v", err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	renderer := mage.TemplateRenderer{};
	renderer.TemplateName = "base";
	renderer.Template = tpl;

	out.Renderer = &renderer;

	return mage.Redirect{Status:http.StatusOK}
}

func (page *TemplatedPage) OnDestroy(ctx context.Context) {

}

//Has a TemplatedPage. Attaches to each templated page a corresponding json file that specifies translations
type LocalizedPage struct {
	TemplatedPage
	Locale string
}

func (page *LocalizedPage) Process(ctx context.Context, out *mage.RequestOutput) mage.Redirect {
	fname := fmt.Sprintf("%s.html", page.FileName);
	_, err := os.Stat(fname);

	if os.IsNotExist(err) {
		log.Debugf(ctx, "Can't find file %s", fname);
		return mage.Redirect{Status:http.StatusNotFound}
	}

	files := make([]string, 0, 0);
	files = append(files, page.Bases...);
	files = append(files, fname);

	tpl, err := template.ParseFiles(files...);

	if err != nil {
		log.Errorf(ctx, "Cant' parse template files: %v", err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	//get the language hint
	inputs := mage.InputsFromContext(ctx);

	lang := page.Locale;

	_, lok := inputs["X-AppEngine-Country"];

	if lok {
		lang = inputs["X-AppEngine-Country"].Value();
	}

	_, lok = inputs["lang"];

	if lok {
		lang = inputs["lang"].Value();
	}

	//get the base language file
	lbasename := fmt.Sprintf("i18n/%s", "base.json");
	jbase, err := ioutil.ReadFile(lbasename);

	if err != nil {
		log.Errorf(ctx, "Error reading base language file %s: %v", lbasename, err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	var base map[string]interface{};
	err = json.Unmarshal(jbase, &base);

	if err != nil {
		log.Errorf(ctx, "Invalid json for base file %s: %v", lbasename, err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	_, bok := base[lang];

	if !bok {
		log.Errorf(ctx, "Base language file %s doesn't support language %s", lbasename, lang);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	globals := base[lang];

	//---- get the specific language json file
	lfname := fmt.Sprintf("i18n/%s.%s", page.FileName, "json")
	//now that we have the locale, read the json language file and get the corresponding values
	jlang, err := ioutil.ReadFile(lfname)

	if err != nil {
		log.Errorf(ctx, "Error retrieving language file %s: %v", lfname, err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	var contents map[string]interface{};
	err = json.Unmarshal(jlang, &contents);

	if err != nil {
		log.Errorf(ctx, "Invalid json for file %s: %v", lfname, err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	_, dok := contents[lang];

	if !dok {
		log.Errorf(ctx, "File %s doesn't support language %s", lfname, lang);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	data := contents[lang];

	renderer := mage.TemplateRenderer{};
	renderer.TemplateName = "base";
	renderer.Template = tpl;
	renderer.Data = struct {
		Url string
		Language string
		Globals interface{}
		Content interface{}
	}{
		page.Url,
		lang,
		globals,
		data,
	};

	out.Renderer = &renderer;

	return mage.Redirect{Status:http.StatusOK}
}

//sends an email with the specified message and sender
type SendMailPage struct {
	mage.Page
	Sender string
	Recipients []string
	Message string
}

func NewSendMailPage(from string, to []string, message string) {
	p := SendMailPage{};
	p.Sender = from;
	p.Recipients = to;
	p.Message = message;
}

func (page *SendMailPage) Process(ctx context.Context, out *mage.RequestOutput) mage.Redirect {

	inputs := mage.InputsFromContext(ctx);

	method := inputs[mage.REQUEST_METHOD].Value();

	if method != http.MethodPost {
		return mage.Redirect{Status:http.StatusMethodNotAllowed};
	}

	renderer := mage.JSONRenderer{};
	out.Renderer = &renderer;

	_, cok := inputs["cname"];

	if !cok {
		renderer.Data = struct {
			Field string
		} {"cname"}
		return mage.Redirect{Status:http.StatusBadRequest};
	}

	cname := inputs["cname"].Value();

	_, cok = inputs["cemail"];

	if !cok {
		renderer.Data = struct {
			Field string
		} {"cemail"}
		return mage.Redirect{Status:http.StatusBadRequest};
	}

	cemail := inputs["cemail"].Value();
	_, err := netmail.ParseAddress(cemail);

	if err != nil {
		renderer.Data = struct {
			Field string
		} {"cemail"}
		return mage.Redirect{Status:http.StatusBadRequest};
	}

	subj := fmt.Sprintf("Nuovo contatto da %s", cname);

	cmsg := "";

	if imsg, cok := inputs["cmessage"]; cok {
		cmsg = imsg.Value();
	}

	cphone := "";
	if iphone, cok := inputs["cphone"];cok {
		cphone = iphone.Value();
	}

	msg := mail.Message{
		Sender: page.Sender,
		To: page.Recipients,
		ReplyTo: cemail,
		Subject: subj,
		Body: fmt.Sprintf(page.Message, cname, cemail, cphone, cmsg),
	};

	err = mail.Send(ctx, &msg);
	if err != nil {
		log.Errorf(ctx, "Couldn't send email message from sender %s to recipient %s: %v", page.Sender, cemail, err);
		return mage.Redirect{Status:http.StatusInternalServerError};
	}

	renderer.Data = struct {

	}{};
	return mage.Redirect{Status:http.StatusOK};
}

func (page *SendMailPage) OnDestroy(ctx context.Context) {

}