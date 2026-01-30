package table

import (
	"log"

	"gorm.io/gorm"
)

func NewUpdateAtTrigger(db *gorm.DB) error {
	// 创建通用的更新时间戳函数（只执行一次）
	createUpdateTimeFunc := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;`

	if err := db.Exec(createUpdateTimeFunc).Error; err != nil {
		log.Printf("创建通用更新时间戳函数失败: %v", err)
		return err
	}

	log.Println("数据库通用函数初始化完成")
	return nil
}
