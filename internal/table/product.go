package table

import (
	"log"

	"gorm.io/gorm"
)

func NewProductTable(db *gorm.DB) error {
	table := `CREATE TABLE IF NOT EXISTS products(
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			quantity BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	err := db.Exec(table).Error
	if err != nil {
		log.Printf("NewProductTable(table): 创建商品表失败: %v", err)
		return err
	}
	// 创建用户表更新时间戳触发器
	createUserUpdateTrigger := `
		DROP TRIGGER IF EXISTS update_products_updated_at ON products;
		CREATE TRIGGER update_products_updated_at
		BEFORE UPDATE ON products
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();`

	if err := db.Exec(createUserUpdateTrigger).Error; err != nil {
		log.Printf("创建用户表更新时间戳触发器失败: %v", err)
		return err
	}

	log.Println("用户表更新时间戳触发器初始化完成")
	return nil
}
