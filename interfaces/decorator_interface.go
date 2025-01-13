package interfaces

import "cfg_exporter/entities"

type IDecorator interface {
	Name() string
}

type ITableDecorator interface {
	RunTableDecorator(tbl *entities.Table) error
}
type IFieldDecorator interface {
	RunFieldDecorator(tbl *entities.Table, field *entities.Field) error
}
