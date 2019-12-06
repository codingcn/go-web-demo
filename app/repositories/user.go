package repositories

import (
	"go-web-demo/app/models"
	"go-web-demo/kernel/database/gorm"
	"go-web-demo/kernel/zlog"
	"go.uber.org/zap"
)

// 更新用户手机号
func UpdateUserPhoneByUserId(userId uint64, phone string) error {
	userProfile := models.User{}
	err := gorm.Connect.Default.Model(&userProfile).Where("user_id = ?", userId).
		Updates(map[string]interface{}{
			"phone": phone,
		}).Error
	if err != nil {
		zlog.Logger.Error("更新用户手机号失败", zap.Any("error", err))
		return err
	}
	return  nil
}
