package store

import "errors"

var (
	ErrNotFound  = errors.New("没有找到记录")
	ErrOperation = errors.New("操作数据库出现问题")
)
