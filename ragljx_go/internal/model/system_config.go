package model

import "time"

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	ConfigKey   string    `gorm:"size:128;uniqueIndex;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	ConfigType  string    `gorm:"size:32" json:"config_type"`
	Description string    `gorm:"type:text" json:"description"`
	IsEncrypted bool      `gorm:"default:false" json:"is_encrypted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

