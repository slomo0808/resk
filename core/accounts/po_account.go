package accounts

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"time"
)

// 持久化对象是ORM映射的基础
// 1. dbx支持自动映射名称，默认是把驼峰命名转换为下划线命名
// 2. 表名默认是结构体名称转换成下划线命名来映射
// 3. 字段名称默认是field name 转换成下划线命名l来映射，字段映射可以用tag描述
// 4. 使用 uni|unique 的tag值来标识字段唯一索引
// 5. 使用 id|pk 的tag值来表示主键
// 6. 使用 omitempty 的tag值来标识字段更新和写入时会被忽略
// 7. 用 — 的tag值来标识字段在更新，写入，查询时会被忽略

// 账户持久化对象
type Account struct {
	Id           int64           `db:"id,omitempty"`
	AccountNo    string          `db:"account_no,unique"`
	AccountName  string          `db:"account_name"`
	AccountType  int             `db:"account_type"`
	CurrencyCode string          `db:"currency_code"`
	UserId       string          `db:"user_id"`
	Username     sql.NullString  `db:"username"`
	Balance      decimal.Decimal `db:"balance"`
	Status       int             `db:"status"`
	CreatedAt    time.Time       `db:"created_at,omitempty"`
	UpdatedAt    time.Time       `db:"updated_at,omitempty"`
}
