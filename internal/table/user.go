package table

import (
	"log"

	"gorm.io/gorm"
)

func NewUserTable(db *gorm.DB) error {
	table := `CREATE TABLE IF NOT EXISTS users(
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) ,
			phone VARCHAR(20) ,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	err := db.Exec(table).Error
	if err != nil {
		log.Printf("NewUserTable(table): 创建用户表失败: %v", err)
		return err
	}
	// 创建用户表更新时间戳触发器
	createUserUpdateTrigger := `
		DROP TRIGGER IF EXISTS update_users_updated_at ON users;
		CREATE TRIGGER update_users_updated_at
		BEFORE UPDATE ON users
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();`

	if err := db.Exec(createUserUpdateTrigger).Error; err != nil {
		log.Printf("创建用户表更新时间戳触发器失败: %v", err)
		return err
	}

	log.Println("用户表更新时间戳触发器初始化完成")
	return nil
}
