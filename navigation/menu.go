package navigation

import (
	"decodica.com/spellbook"
)

type MenuByOrder Menu

/**
ByOrder start
*/
func (menu MenuByOrder) Len() int {
	return len(menu)
}

func (menu MenuByOrder) Swap(i, j int) {
	menu[i], menu[j] = menu[j], menu[i]
}

func (menu MenuByOrder) Less(i, j int) bool {
	return menu[i].Order < menu[j].Order
}

type Menu []MenuItem

func (menu Menu) ItemByCode(code spellbook.StaticPageCode) *MenuItem {
	for _, v := range menu {
		if v.Code == code {
			return &v
		}
	}
	return nil
}

// menu is derived from the page resource
type MenuItem struct {
	Url string
	Locale string
	Code spellbook.StaticPageCode
	Label string
	Order int
}

func NewMenuItemFromPage(page *Page) MenuItem {
	return MenuItem {
		Url: page.Url,
		Locale: page.Locale,
		Code: page.Code,
		Label: page.Label,
		Order: page.Order,
	}
}

func (item MenuItem) LocalizedUrl() string {
	return "/" + item.Locale + "/" + item.Url
}
