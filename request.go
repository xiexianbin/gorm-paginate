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

// Where transform to SQL `Where` Condition
type Where struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // eq, ne, gt, etc.
	Value    any    `json:"value"`
}

// WhereCondition and? or?

// OrderBy transform to SQL `order by`
type OrderBy struct {
	Field     string `json:"field"`
	Direction string `json:"direction"` // asc, desc
}

// QueryParams Where + OrderBy + Pagination
type QueryParams struct {
	Page     int       `json:"page"` // page and size params, transform to SQL `limit` and `offset`
	Size     int       `json:"size"`
	OrderBys []OrderBy `json:"order_bys"`
	Wheres   []Where   `json:"wheres"`
	Comments []string  `json:"comments"` // when some condition not workk, comment it
}
