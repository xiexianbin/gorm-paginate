// Copyright 2025~time.Now xiexianbin<me@xiexianbin.cn>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package paginate

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Paginate GORM paginated Response
type Pagination struct {
	QueryParams
	Items      any   `json:"items"`
	Total      int64 `json:"total"`       // total items
	TotalPages int64 `json:"total_pages"` // total pages
}

// parseOrderBy parse order by
func (p *Pagination) parseOrderBy(orderBy string, validFields map[string]bool) (orders []OrderBy) {
	parts := strings.Split(orderBy, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fieldDir := strings.Fields(part)
		if len(fieldDir) == 0 {
			continue
		}
		field := fieldDir[0]
		direction := "asc"
		if len(fieldDir) > 1 && strings.ToLower(fieldDir[1]) == "desc" {
			direction = "desc"
		}
		if strings.HasPrefix(field, "-") {
			field = strings.Replace(field, "-", "", 1)
			direction = "desc"
		}
		if validFields[field] {
			orders = append(orders, OrderBy{Field: field, Direction: direction})
		}
	}
	return
}

// parseWhere parse where
func (p *Pagination) parseWhere(query url.Values, validFields map[string]bool) []Where {
	var filters []Where
	operatorMap := map[string]string{
		"eq":      "=",
		"ne":      "!=",
		"gt":      ">",
		"gte":     ">=",
		"lt":      "<",
		"lte":     "<=",
		"like":    "LIKE",
		"notlike": "NOT LIKE",
		"is":      "IS",
		"isnot":   "IS NOT",
		"in":      "IN",
		// "between": ?
	}

	for key, values := range query {
		if key == "page" || key == "size" || key == "order_by" {
			continue
		}
		if len(values) == 0 {
			continue
		}
		value := values[0]

		// Split fields and operators
		lastIndex := strings.LastIndex(key, "_")
		var field, operator string
		if lastIndex == -1 {
			field = key
			operator = "eq"
		} else {
			field = key[:lastIndex]
			operator = key[lastIndex+1:]
		}

		// Validate fields and operators
		if !validFields[field] || operatorMap[operator] == "" {
			continue
		}

		// Handle special operators: in
		var val interface{} = value
		if operator == "in" {
			val = strings.Split(value, ",")
		}

		filters = append(filters, Where{
			Field:    field,
			Operator: operator,
			Value:    val,
		})
	}
	return filters
}

// Parse parse all to QueryParams
func (p *Pagination) Parse(query url.Values, validFields map[string]bool) {
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	p.Page = page

	size, _ := strconv.Atoi(query.Get("size"))
	if size < 1 || size > 2000 {
		size = 10
	}
	p.Size = size

	p.QueryParams.Wheres = p.parseWhere(query, validFields)
	p.QueryParams.OrderBys = p.parseOrderBy(query.Get("order_by"), validFields)
}

// ParseModelFields parse gorm model fields name from struct
func ParseModelFields(model any) ([]string, error) {
	sch, err := schema.Parse(model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, err
	}
	fields := make([]string, 0)
	for _, field := range sch.Fields {
		fields = append(fields, field.DBName)
	}
	return fields, nil
}

// applyOffsetAndLimit apply offset and limit to *gorm.DB
func applyOffsetAndLimit(db *gorm.DB, pagination *Pagination) *gorm.DB {
	offset := (pagination.Page - 1) * pagination.Size
	return db.Offset(offset).Limit(pagination.Size)
}

// applyOrderBy apply OrderBy to *gorm.DB
func applyOrderBy(db *gorm.DB, orders []OrderBy) *gorm.DB {
	for _, order := range orders {
		db = db.Order(fmt.Sprintf("%s %s", order.Field, order.Direction))
	}
	return db
}

// applyWhere apply Where to *gorm.DB
func applyWhere(db *gorm.DB, filters []Where) *gorm.DB {
	operatorMap := map[string]string{
		"eq":       "= ?",
		"ne":       "!= ?",
		"gt":       "> ?",
		"gte":      ">= ?",
		"lt":       "< ?",
		"lte":      "<= ?",
		"like":     "LIKE ?",
		"not like": "NOT LIKE ?",
		"is":       "is ?",
		"is not":   "is not ?",
		"in":       "IN (?)",
	}

	for _, filter := range filters {
		sqlOp := operatorMap[filter.Operator]
		query := fmt.Sprintf("%s %s", filter.Field, sqlOp)
		db = db.Where(query, filter.Value)
	}
	return db
}

// Paginate use gorm [scopes](https://gorm.io/docs/scopes.html#Pagination) to apply query params
//
// Usage:
//
//	db.Scopes(Paginate(model, query, paginate)).Find(&users)
func Paginate(model any, query url.Values, pagination *Pagination, tx2 *gorm.DB) func(db *gorm.DB) *gorm.DB {
	// tx2 fix: sql: expected 8 destination arguments in Scan, not 1; sql: ...
	return func(tx *gorm.DB) *gorm.DB {
		fields, _ := ParseModelFields(model)
		validFields := make(map[string]bool, len(fields))
		for _, field := range fields {
			validFields[field] = true
		}

		pagination.Parse(query, validFields)
		var totalRows int64
		tx2.Model(model).Count(&totalRows)
		pagination.Total = totalRows
		pagination.TotalPages = int64(math.Ceil(float64(totalRows) / float64(pagination.QueryParams.Size)))

		tx = applyWhere(tx, pagination.QueryParams.Wheres)
		tx = applyOffsetAndLimit(tx, pagination)
		tx = applyOrderBy(tx, pagination.QueryParams.OrderBys)
		return tx
	}
}
