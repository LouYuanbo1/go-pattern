package table

import (
	"fmt"

	"gorm.io/gorm"
)

func NewAllTables(db *gorm.DB) error {
	err := NewUpdateAtTrigger(db)
	if err != nil {
		return fmt.Errorf("创建触发器函数错误")
	}
	err = NewUserTable(db)
	if err != nil {
		return fmt.Errorf("创建用户表错误")
	}
	err = NewProductTable(db)
	if err != nil {
		return fmt.Errorf("创建商品表错误")
	}
	err = NewOrderTable(db)
	if err != nil {
		return fmt.Errorf("创建订单表错误")
	}
	return err
}
