package page

import (
	"distudio.com/mage"
	"distudio.com/mage/model"
	"distudio.com/page/resource/content"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type BaseController struct {
	mage.Controller
}

type Paging struct {
	page        int
	size        int
	order       model.Order
	orderField  string
	filterField string
	filterValue string
}

const orderAscKey = "asc"
const orderDescKey = "desc"

func (controller *BaseController) OnDestroy(ctx context.Context) {}

func (controller *BaseController) GetPaging(ins mage.RequestInputs) (*Paging, error) {
	page := 0
	size := 20
	var order model.Order
	orderField := ""
	filterField := ""
	filterValue := ""

	if pin, ok := ins["page"]; ok {
		if num, err := strconv.Atoi(pin.Value()); err == nil {
			page = num
		} else {
			return nil, errors.New("error paging page")
		}
	}

	if sin, ok := ins["results"]; ok {
		if num, err := strconv.Atoi(sin.Value()); err == nil {
			size = num
			// cap the size to 100
			if size > 100 {
				size = 100
			}
		} else {
			return nil, errors.New("error paging results")
		}
	}

	if oin, ok := ins["order"]; ok {
		oins := oin.Value()
		if len(oins) > 0 {
			if strings.Compare(oins, orderAscKey) == 0 {
				order = model.ASC
			} else if strings.Compare(oin.Value(), orderDescKey) == 0 {
				order = model.DESC
			} else {
				return nil, errors.New("error order")
			}
		}
	}

	if ofin, ok := ins["orderField"]; ok {
		orderField = ofin.Value()
	}

	if ffin, ok := ins["filterField"]; ok {
		filterField = ffin.Value()
	}
	if fvin, ok := ins["filterValue"]; ok {
		filterValue = fvin.Value()
	}

	var paging Paging
	paging.page = page
	paging.size = size
	paging.order = order
	paging.orderField = orderField
	paging.filterField = filterField
	paging.filterValue = filterValue

	return &paging, nil
}

func (controller *BaseController) HandlerPagingResult(ctx context.Context, out *mage.ResponseOutput) mage.Redirect {
	// todo generalize not used
	// handle query params for page data:
	ins := mage.InputsFromContext(ctx)
	paging, err := controller.GetPaging(ins)
	if err != nil {
		return mage.Redirect{Status: http.StatusBadRequest}
	}
	page := paging.page
	size := paging.size

	var result interface{}
	l := 0
	// check property
	property, ok := ins["property"]
	if ok {
		// property
		properties, err := controller.HandleResourceProperties(ctx, property.Value(), page, size)
		if err != nil {
			log.Errorf(ctx, "Error retrieving posts property: %v %+v", property, err)
			return mage.Redirect{Status: http.StatusInternalServerError}
		}
		l = len(properties)
		result = properties[:controller.GetCorrectCountForPaging(size, l)]
	} else {
		// list contents
		var conts []*content.Content
		q := model.NewQuery(&content.Content{})
		q = q.OffsetBy(page * size)
		// get one more so we know if we are done
		q = q.Limit(size + 1)
		err := q.GetMulti(ctx, &conts)
		if err != nil {
			log.Errorf(ctx, "Error retrieving list posts %+v", err)
			return mage.Redirect{Status: http.StatusInternalServerError}
		}
		l = len(conts)
		result = conts[:controller.GetCorrectCountForPaging(size, l)]
	}

	// todo: generalize list handling and responses
	response := struct {
		Items interface{} `json:"items"`
		More  bool        `json:"more"`
	}{result, l > size}
	renderer := mage.JSONRenderer{}
	renderer.Data = response
	out.Renderer = &renderer
	return mage.Redirect{Status: http.StatusOK}
}

func (controller *BaseController) GetCorrectCountForPaging(size int, l int) int {
	count := size
	if l < size {
		count = l
	}
	return count
}

func (controller *BaseController) UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func (controller *BaseController) HandleResourceProperties(ctx context.Context, property string, page int, size int) ([]interface{}, error) {
	// todo: generalize not used
	//name := ""
	//switch property {
	//// content
	//case "category":
	//	name = "Category"
	//case "topic":
	//	name = "Topic"
	//case "name":
	//	name = "Name"
	//	// attachment
	//case "group":
	//	name = "Group"
	//default:
	//	return nil, errors.New("no property found")
	//}
	name := controller.UcFirst(property)

	var posts []*content.Content
	q := model.NewQuery(&content.Content{})
	q = q.OffsetBy(page * size)
	q = q.Distinct(name)
	// get one more so we know if we are done
	q = q.Limit(size + 1)
	err := q.GetAll(ctx, &posts)
	if err != nil {
		log.Errorf(ctx, "Error retrieving result: %+v", err)
		return nil, err
	}
	var result []interface{}
	for _, p := range posts {
		value := reflect.ValueOf(p).Elem().FieldByName(name).String()
		result = append(result, &value)
	}
	return result, nil

}
