package data

import (
	"errors"
	"github.com/zngue/zng_app/db/data/page"
	"github.com/zngue/zng_app/db/data/where"
	"gorm.io/gorm"
)

var ErrData = errors.New("data does not exist")

type Fn func(db *gorm.DB) *gorm.DB

type DB[T any] struct {
	Source *gorm.DB
	Model  *T
}
type ListRequest struct {
	*page.Page
	Where  map[string]any
	Order  []string
	Select any
	Fn     Fn
}

type OptionFn func(r *ListRequest)

func WhereStruct(v any) OptionFn {
	return func(r *ListRequest) {
		if v != nil {
			options := where.Where(v)
			r.Where = where.NewWhere(options...)
		}
	}
}
func WhereOption(fns ...where.Fn) OptionFn {
	return func(r *ListRequest) {
		if len(fns) > 0 {
			r.Where = where.NewWhereFn(fns...)
		}
	}
}
func OrderOption(v []string) OptionFn {
	return func(r *ListRequest) {
		if v != nil {
			r.Order = v
		}
	}
}
func SelectOption(v any) OptionFn {
	return func(r *ListRequest) {
		r.Select = v
	}
}
func FnWithData(fn Fn) OptionFn {
	return func(r *ListRequest) {
		r.Fn = fn
	}
}
func PageWithData(fns ...page.Fn) OptionFn {
	return func(r *ListRequest) {
		r.Page = page.NewPage(fns...)
	}
}

type ContentRequest struct {
	Where  map[string]any
	Select any
	Fn     Fn
	Order  []string
}

// ContentFn 获取单条数据
func (d *DB[T]) ContentFn(fns ...OptionFn) (resData *T, err error) {
	var data = &ListRequest{}
	for _, fn := range fns {
		fn(data)
	}
	db := d.Source.Model(d.Model)
	db = d.ListHelper(db, data)
	err = db.First(&resData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrData
	}
	return
}

// Content 获取单条数据
func (d *DB[T]) Content(data *ContentRequest) (resData *T, err error) {
	db := d.Source.Model(d.Model)
	db = d.ListHelper(db, &ListRequest{
		Where:  data.Where,
		Order:  data.Order,
		Select: data.Select,
		Fn:     data.Fn,
	})
	err = db.First(&resData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrData
	}
	return
}

// ListFn 获取列表带自定义Fn
func (d *DB[T]) ListFn(fns ...OptionFn) (list []*T, err error) {
	var data = &ListRequest{
		Page: page.NewPage(),
	}
	for _, fn := range fns {
		fn(data)
	}
	db := d.Source.Model(d.Model)
	db = d.ListHelper(db, data)
	if data.Page.Page != -1 {
		db = data.Page.PageHandle(db)
	}
	err = db.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// List 获取列表
func (d *DB[T]) List(data *ListRequest) (list []*T, err error) {
	if data.Page == nil {
		data.Page = page.NewPage()
	}
	db := d.Source.Model(d.Model)
	db = d.ListHelper(db, data)
	if data.Page.Page != -1 {
		db = data.Page.PageHandle(db)
	}
	err = db.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// ListHelper ListRequest
func (d *DB[T]) ListHelper(db *gorm.DB, data *ListRequest) *gorm.DB {

	if data.Where != nil && len(data.Where) > 0 {
		db = d.Where(data.Where, db)
	}
	if data.Order != nil && len(data.Order) > 0 {
		db = d.Order(data.Order, db)
	}
	if data.Select != nil {
		db = d.Select(db, data.Select)
	}
	if data.Fn != nil {
		db = data.Fn(db)
	}
	return db
}

// ListPageFn 获取列表带分页
func (d *DB[T]) ListPageFn(fns ...OptionFn) (list []*T, count int64, err error) {
	var data = &ListRequest{}
	for _, fn := range fns {
		fn(data)
	}
	db := d.Source.Model(d.Model)
	db = d.ListHelper(db, data)
	if data.Page != nil && data.Page.Page != -1 {
		err = db.Count(&count).Error
		if err != nil {
			return
		}
		db = data.Page.PageHandle(db)
	}
	err = db.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// ListPage 获取列表带分页
func (d *DB[T]) ListPage(data *ListRequest) (list []*T, count int64, err error) {
	db := d.Source.Model(d.Model)
	db = d.ListHelper(db, data)
	if data.Page != nil && data.Page.Page != -1 {
		err = db.Count(&count).Error
		if err != nil {
			return
		}
		db = data.Page.PageHandle(db)
	}
	err = db.Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// Select 设置查询字段
func (d *DB[T]) Select(db *gorm.DB, data any) *gorm.DB {
	if data != nil {
		db = db.Select(data)
	}
	return db
}

// Add 新增
func (d *DB[T]) Add(data *T) (err error) {
	db := d.Source.Model(d.Model)
	err = db.Create(data).Error
	return
}

// AddMore 批量新增
func (d *DB[T]) AddMore(data []*T) (err error) {
	db := d.Source.Model(d.Model)
	err = db.Create(data).Error
	return
}

// Where map[string]any where条件
func (d *DB[T]) Where(data map[string]any, db *gorm.DB) *gorm.DB {
	if data != nil && len(data) > 0 {
		for k, v := range data {
			db = db.Where(k, v)
		}
	}
	return db
}

// Order []string 排序条件
func (d *DB[T]) Order(data []string, db *gorm.DB) *gorm.DB {
	if data != nil && len(data) > 0 {
		for _, v := range data {
			db = db.Order(v)
		}
	}
	return db
}

// Update 更新 where map  data map
func (d *DB[T]) Update(where, data map[string]any) (err error) {
	db := d.Source.Model(d.Model)
	if where == nil || len(where) == 0 {
		err = errors.New("where条件不能为空")
		return
	}
	db = d.Where(where, db)
	err = db.Updates(data).Error
	return
}

// Delete 删除
func (d *DB[T]) Delete(where map[string]any) (err error) {
	db := d.Source.Model(d.Model)
	if where == nil || len(where) == 0 {
		err = errors.New("where条件不能为空")
		return
	}
	db = d.Where(where, db)
	err = db.Delete(d.Model).Error
	return
}

type DelRequest struct {
	Where map[string]any
	Fn    Fn
}

// DeleteFn 关联删除
func (d *DB[T]) DeleteFn(data *DelRequest) (err error) {
	db := d.Source.Model(d.Model)
	if data.Where == nil || len(data.Where) == 0 {
		err = errors.New("where条件不能为空")
		return
	}
	db = d.Where(data.Where, db)
	if data.Fn != nil {
		db = data.Fn(db)
	}
	err = db.Delete(d.Model).Error
	return
}

// Execute 自定义模型
func (d *DB[T]) Execute(fn DBFn) (err error) {
	if fn == nil {
		return
	}
	err = fn(d.Source)
	return
}

type DBFn func(db *gorm.DB) (err error)

// Transaction 事务提交
func (d *DB[T]) Transaction(fn DBFn) (err error) {
	begin := d.Source.Begin()
	defer begin.Rollback()
	if fn == nil {
		err = errors.New("事务函数不能为空")
		return
	}
	err = fn(begin)
	if err != nil {
		return
	}
	err = begin.Commit().Error
	return
}

// NewDB 实例化 DB
func NewDB[T any](db *gorm.DB) *DB[T] {
	model := new(T)
	return &DB[T]{
		Source: db,
		Model:  model,
	}
}
