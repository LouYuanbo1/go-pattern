package table

import (
	"log"

	"gorm.io/gorm"
)

func NewOrderTable(db *gorm.DB) error {
	table := `CREATE TABLE IF NOT EXISTS orders(
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			product_id BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	err := db.Exec(table).Error
	if err != nil {
		log.Printf("NewOrderTable(table): 创建订单表失败: %v", err)
		return err
	}
	// 创建用户表更新时间戳触发器
	createUserUpdateTrigger := `
	DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
	CREATE TRIGGER update_orders_updated_at
	BEFORE UPDATE ON orders
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();`

	if err := db.Exec(createUserUpdateTrigger).Error; err != nil {
		log.Printf("创建用户表更新时间戳触发器失败: %v", err)
		return err
	}

	log.Println("用户表更新时间戳触发器初始化完成")
	return nil
}
