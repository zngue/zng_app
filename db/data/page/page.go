package page

import "gorm.io/gorm"

type Page struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

type Fn func(page *Page)

func DataWithPage(page int) Fn {
	return func(pageOption *Page) {
		pageOption.Page = page
	}
}
func DataWithPageSize(pageSize int) Fn {
	return func(page *Page) {
		page.PageSize = pageSize
	}
}
func (p *Page) PageHandle(db *gorm.DB) *gorm.DB {
	if p.Page == -1 { //不分页
		return db
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	offset := (p.Page - 1) * p.PageSize
	return db.Offset(offset).Limit(p.PageSize)
}
func NewPage(fns ...Fn) *Page {
	var page = new(Page)
	for _, fn := range fns {
		fn(page)
	}
	if page.Page == 0 {
		page.Page = -1
	}
	return page
}
