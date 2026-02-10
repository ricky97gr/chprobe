package pagination

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ricky97gr/chprobe/chprobe_common/utils"
)

// PageQuery 分页查询结构
type PageQuery struct {
	Page       int         `json:"page" form:"page"`             // 页码
	PageSize   int         `json:"pageSize" form:"pageSize"`     // 每页大小
	StartTime  int32       `json:"startTime" form:"startTime"`   // 开始时间
	EndTime    int32       `json:"endTime" form:"endTime"`       // 结束时间
	Conditions []Condition `json:"conditions" form:"conditions"` // 查询条件
	Sorts      []Sort      `json:"sorts" form:"sorts"`           // 排序规则
}

// Condition 查询条件结构
type Condition struct {
	Field     string      `json:"field" form:"field"`         // 字段名
	Value     interface{} `json:"value" form:"value"`         // 字段值
	Operation int         `json:"operation" form:"operation"` // 操作符
}

// Sort 排序规则结构
type Sort struct {
	Field   string `json:"field" form:"field"`     // 排序字段
	OrderBy int    `json:"orderBy" form:"orderBy"` // 排序方向
}

// 操作符常量
const (
	Equal        = iota + 1 // = 等于
	NotEqual                // != 不等于
	GreaterThan             // > 大于
	GreaterEqual            // >= 大于等于
	LessThan                // < 小于
	LessEqual               // <= 小于等于
	Like                    // like 模糊查询
	In                      // in 包含
	NotIn                   // not in 不包含
)

// 排序方向常量
const (
	Asc  = 1 // 升序
	Desc = 2 // 降序
)

// GetPageQuery 从请求中获取分页查询参数
func GetPageQuery(ctx *gin.Context) (PageQuery, error) {
	var page PageQuery
	var s []Sort
	var conditons []Condition
	var pageNumber int = 1
	var pageSize int = 10

	// 解析sorts参数（JSON格式）
	if orderStr, ok := ctx.GetQuery("sorts"); ok {
		err := json.Unmarshal([]byte(orderStr), &s)
		if err != nil {
			utils.Logger.Errorf("failed to unmarshal sorts, err: %+v\n", err)
		}
	}

	// 解析conditions参数（JSON格式）
	if conStr, ok := ctx.GetQuery("conditions"); ok {
		err := json.Unmarshal([]byte(conStr), &conditons)
		if err != nil {
			utils.Logger.Errorf("failed to unmarshal conditions, err: %+v\n", err)
		}
	}

	// 解析page参数
	if pageNumberStr, ok := ctx.GetQuery("page"); ok {
		n, err := strconv.Atoi(pageNumberStr)
		if err != nil {
			utils.Logger.Errorf("failed to get page number, err: %+v, then set pageNumber to 1\n", err)
			n = 1
		}
		pageNumber = n
	}

	// 解析pageSize参数
	if pageSizeStr, ok := ctx.GetQuery("pageSize"); ok {
		n, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			utils.Logger.Errorf("failed to get page size, err: %+v, then set pageSize to 20\n", err)
			n = 20
		}
		pageSize = n
	}

	// 设置参数
	page.Page = pageNumber
	page.PageSize = pageSize
	page.Sorts = s
	page.Conditions = conditons

	return page, nil
}

// GetOffset 获取偏移量
func (p *PageQuery) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取限制数
func (p *PageQuery) GetLimit() int {
	return p.PageSize
}

// Order 生成排序条件
func Order(sorts []Sort) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		resultDB := db
		for _, sort := range sorts {
			if sort.OrderBy == 1 {
				resultDB = resultDB.Order(sort.Field + " asc")
			}
			if sort.OrderBy == -1 {
				resultDB = resultDB.Order(sort.Field + " desc")
			}
		}
		return resultDB
	}
}

// GetCondition 获取指定字段的条件
func (q PageQuery) GetCondition(field string) (Condition, bool) {
	for _, cond := range q.Conditions {
		if cond.Field == field {
			return cond, true
		}
	}
	return Condition{}, false
}

// ParseQuery 生成完整的查询条件
func ParseQuery(q PageQuery) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		resultDB := db
		for _, v := range q.Conditions {
			resultDB = resultDB.Scopes(QueryFilter(v.Field, v.Value, v.Operation))
		}
		resultDB = resultDB.Scopes(Order(q.Sorts))
		resultDB = resultDB.Scopes(QueryLimitShip(q.Page, q.PageSize))
		return resultDB
	}
}

// QueryFilter 生成过滤条件
func QueryFilter(field string, value interface{}, operation int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch operation {
		case Equal:
			return db.Where(field+" = ?", value)
		case NotEqual:
			return db.Where(field+" != ?", value)
		case GreaterThan:
			return db.Where(field+" > ?", value)
		case GreaterEqual:
			return db.Where(field+" >= ?", value)
		case LessThan:
			return db.Where(field+" < ?", value)
		case LessEqual:
			return db.Where(field+" <= ?", value)
		case Like:
			return db.Where(field+" like ?", "%"+value.(string)+"%")
		case In:
			return db.Where(field+" in ?", value)
		case NotIn:
			return db.Where(field+" not in ?", value)
		default:
			return db
		}
	}
}

// QueryLimitShip 生成分页限制
func QueryLimitShip(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
