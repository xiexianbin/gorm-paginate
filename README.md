# gorm-paginate

[GORM]（https://gorm.io/） 分页插件（[program implementation for]（https://github.com/xiexianbin/gin-template/commit/926b25832fefe611dcd119b6cd46cbd80806c386））

## 特性

- **分页**
  - `page`：当前页（默认 1）
  - `size`：每页记录数（默认 10）

- **排序**
  - `order_by`：排序字段和方向（如 `created_at desc, name`）

- **过滤条件**
  - 格式：`<field>_<operator>=<value>`
    - 如：`age_gt=20`（where age > 20）、`name_like=John%`（where LIKE 'John%'）
  - 支持的比较
    - `eq` 等于（默认）
    - `ne` 不等于
    - `gt` 大于
    - `gte` 大于等于
    - `lt` 小于
    - `lte` 小于等于
    - `between` 范围搜索
    - `like` 模糊匹配
    - `notlike` 非模糊匹配
    - `is`
    - `isnot`
    - `in`

## 示例

1. 简单分页和排序

```
GET http://localhost:8080/users?size=10&page=0&order_by=-name,id
```

- 等价的 SQL

```
SELECT * FROM users ORDER BY name DESC, id ASC LIMIT 10 OFFSET 0
```

- JSON 响应:

```
{
    "items": [
        {
            "id": 1,
            "name": "xiexianbin",
            "age": 18
        }
    ],
    "page": 0,
    "page_size": 10,
    "total_pages": 1
}
```

2. 条件搜索

```
GET http://localhost:8080/users?size=10&page=0&age_gt=16&name_like=xie%&balance_between=20,25&account_manager_in=zhangsan,lisi
```

- 等价的 SQL

```
SELECT * FROM users WHERE age > 16 AND name LIKE "xie%" and balance BETWEEN 20 AND 25 and account_manager in ("zhangsan", "lisi") LIMIT 10 OFFSET 0
```

## License

© xiexianbin, 2025~time.Now

Released under the [Apache License](https://github.com/xiexianbin/gorm-paginate/blob/main/LICENSE)
