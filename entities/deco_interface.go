package entities

type IDecorator interface {
	Name() string
}

type ITableDecorator interface {
	RunTableDecorator(tbl *Table) error
}
type IFieldDecorator interface {
	RunFieldDecorator(tbl *Table, field *Field) error
}
