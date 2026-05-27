package db

import (
	"context"
	"errors"
	"strings"

	"github.com/zngue/zng_app/pkg/errors_ez"
	"gorm.io/gorm"
)

const (
	PageNoPagination int32 = -1
)

type Pages struct {
	Page     int32
	PageSize int32
}

type ListResponse[T any] struct {
	List     []*T   // 列表数据
	Count    uint32 // 总数
	Page     int32  // 当前页
	PageSize int32  // 每页条数
	IsCount  bool   // 是否统计总数
}

// PageHandle Page 分页处理
func (p *Pages) PageHandle(db *gorm.DB) *gorm.DB {
	if p.Page == PageNoPagination {
		return db
	}
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	offset := (p.Page - 1) * p.PageSize
	return db.Offset(int(offset)).Limit(int(p.PageSize))
}

func NewPages(pages ...int32) *Pages {
	page := new(Pages)
	if len(pages) > 0 {
		page.Page = pages[0]
	}
	if len(pages) > 1 {
		page.PageSize = pages[1]
	}
	return page
}

type Fn func(db *gorm.DB) *gorm.DB

type ConnDB[T any] struct {
	source *gorm.DB
	model  *T
}

type ListRequest struct {
	Page     int32
	PageSize int32
	Where    map[string]any
	IsCount  bool // 是否统计总数
	Order    []string
	Select   any
	Fn       Fn
}
type ContentRequest struct {
	Where  map[string]any
	Select any
	Fn     Fn
	Order  []string
}

func (d *ConnDB[T]) ContentById(ctx context.Context, id uint32) (resData *T, err error) {
	var where = map[string]any{
		"id = ?": id,
	}
	return d.Content(ctx, &ContentRequest{
		Where: where,
	})
}

// Content 获取单条数据
func (d *ConnDB[T]) Content(ctx context.Context, data *ContentRequest) (resData *T, err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	db = d.ListHelper(db, &ListRequest{
		Where:  data.Where,
		Order:  data.Order,
		Select: data.Select,
		Fn:     data.Fn,
	})
	err = db.Take(&resData).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors_ez.Wrap(err, "数据不存在")
		return
	}
	if err != nil {
		err = errors_ez.Wrap(err)
		resData = nil
	}
	return
}

// ListHelper ListRequest
func (d *ConnDB[T]) ListHelper(db *gorm.DB, data *ListRequest) *gorm.DB {
	if len(data.Where) > 0 {
		db = d.Where(data.Where, db)
	}
	if len(data.Order) > 0 {
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

// ListPage 获取列表带分页
func (d *ConnDB[T]) List(ctx context.Context, req *ListRequest) (data *ListResponse[T], err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	db = d.ListHelper(db, req)
	page := NewPages(req.Page, req.PageSize)
	var total int64
	if req.Page != PageNoPagination {
		db = page.PageHandle(db)
	}
	if req.IsCount {
		if err = db.Count(&total).Error; err != nil {
			err = errors_ez.Wrap(err, "count error")
			return
		}
	}
	var list []*T
	if err = db.Find(&list).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors_ez.Wrap(err, "list error")
			return
		}
	}
	return &ListResponse[T]{
		List:     list,
		Count:    uint32(total),
		Page:     page.Page,
		PageSize: page.PageSize,
		IsCount:  req.IsCount,
	}, nil
}

// Select 设置查询字段
func (d *ConnDB[T]) Select(db *gorm.DB, data any) *gorm.DB {
	if data != nil {
		db = db.Select(data)
	}
	return db
}

// Add 新增
func (d *ConnDB[T]) Add(ctx context.Context, data *T) (err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	err = db.Create(data).Error
	return
}

// AddMore 新增
func (d *ConnDB[T]) AddMore(ctx context.Context, data []*T) (err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	err = db.Create(data).Error
	return
}

// Where map[string]any where条件
func (d *ConnDB[T]) Where(data map[string]any, db *gorm.DB) *gorm.DB {
	if len(data) > 0 {
		for k, v := range data {
			if strings.Count(k, "?") > 1 {
				itemSlice, ok := v.([]any)
				if !ok {
					continue
				}
				db = db.Where(k, itemSlice...)
			} else {
				db = db.Where(k, v)
			}
		}
	}
	return db
}

// Order []string 排序条件
func (d *ConnDB[T]) Order(data []string, db *gorm.DB) *gorm.DB {
	if len(data) > 0 {
		for _, v := range data {
			db = db.Order(v)
		}
	}
	return db
}

// Update 更新 where map  data map
func (d *ConnDB[T]) Update(ctx context.Context, where, data map[string]any) (err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	if len(where) == 0 {
		err = errors_ez.Wrap(err, "更新条件不能为空")
		return
	}
	db = d.Where(where, db)
	err = db.Updates(data).Error

	return
}

func (d *ConnDB[T]) UpdateWithResult(ctx context.Context, where, data map[string]any) (rowsAffected int64, err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	if len(where) == 0 {
		err = errors_ez.Wrap(err, "更新条件不能为空")
		return
	}
	db = d.Where(where, db)
	result := db.Updates(data)
	rowsAffected = result.RowsAffected
	err = result.Error
	if err != nil {
		err = errors_ez.Wrap(err, "更新失败")
		return
	}
	return
}
func (d *ConnDB[T]) UpdateById(ctx context.Context, id uint32, data map[string]any) (err error) {
	var where = map[string]any{
		"id = ?": id,
	}
	return d.Update(ctx, where, data)
}

// DeleteByIds 删除
func (d *ConnDB[T]) DeleteByIds(ctx context.Context, ids []uint32) (err error) {
	var where = map[string]any{
		"id in ?": ids,
	}
	return d.Delete(ctx, where)
}

// Delete 删除
func (d *ConnDB[T]) Delete(ctx context.Context, where map[string]any) (err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	if len(where) == 0 {
		err = errors_ez.Wrap(err, "删除条件不能为空")
		return
	}
	db = d.Where(where, db)
	err = db.Delete(d.model).Error
	if err != nil {
		err = errors_ez.Wrap(err, "删除失败")
	}
	return
}

// Count 统计
func (d *ConnDB[T]) Count(ctx context.Context, where map[string]any) (count int64, err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	db = d.Where(where, db)
	err = db.Count(&count).Error
	if err != nil {
		err = errors_ez.Wrap(err, "统计失败")
	}
	return
}
func (d *ConnDB[T]) Exec(ctx context.Context, where map[string]any, fn func(conn *gorm.DB) (connErr error)) (err error) {
	db := d.source.WithContext(ctx).Model(d.model)
	db = d.Where(where, db)
	err = fn(db)
	if err != nil {
		err = errors_ez.Wrap(err, "执行失败")
	}
	return
}

// NewDB 实例化 DB
func NewDB[T any](source *gorm.DB) *ConnDB[T] {
	model := new(T)
	return &ConnDB[T]{
		source: source,
		model:  model,
	}
}

func NewDBRepo[T any](source *gorm.DB) ConnRepo[T] {
	return NewDB[T](source)
}

type ConnRepo[T any] interface {
	ContentById(ctx context.Context, id uint32) (resData *T, err error)
	List(ctx context.Context, req *ListRequest) (data *ListResponse[T], err error)
	Content(ctx context.Context, req *ContentRequest) (data *T, err error)
	Add(ctx context.Context, req *T) (err error)
	AddMore(ctx context.Context, req []*T) (err error)
	Update(ctx context.Context, where, data map[string]any) (err error)
	UpdateWithResult(ctx context.Context, where, data map[string]any) (rowsAffected int64, err error)
	UpdateById(ctx context.Context, id uint32, data map[string]any) (err error)
	DeleteByIds(ctx context.Context, ids []uint32) (err error)
	Delete(ctx context.Context, where map[string]any) (err error)
	Count(ctx context.Context, where map[string]any) (count int64, err error)
	Exec(ctx context.Context, where map[string]any, fn func(conn *gorm.DB) (err error)) (err error)
}
