package types

type DBTy string

const (
	MemDBTy   DBTy = "memory"
	GoCacheTy      = "go-cache"
	RedisDBTy      = "redis"
	MysqlDBTy      = "mysql"
)

type FindSqlField struct {
	TBName string
	Select SelectField
	Where  WhereField
	Limit  interface{}
	Offset interface{}
	Group  string
	Having HavingField
}

type HavingField struct {
	Query  interface{}
	Values interface{}
}

type SelectField struct {
	Query interface{}
	Args  interface{}
}

type WhereField struct {
	Query interface{}
	Args  interface{}
}

type UpdateField struct {
	TBName string
	Where  WhereField
	Values interface{}
}

type DeleteField struct {
	TBName string
	Where  WhereField
	Values interface{}
}

type RawField struct {
	Sql    string
	Values interface{}
}
