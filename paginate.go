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
	"net/url"
	"strconv"
	"strings"
)

// parsePagination parse page and size from url.Values
func parsePagination(query url.Values) Pagination {
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(query.Get("size"))
	if size < 1 || size > 2000 {
		size = 10
	}
	return Pagination{Page: page, Size: size}
}

// parseOrderBy parse order by
func parseOrderBy(orderBy string, validFields map[string]bool) (orders []OrderBy) {
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
func parseWhere(query url.Values, validFields map[string]bool) []Where {
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
