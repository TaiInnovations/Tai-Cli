package dao

import (
    "Tai/internal/database"
    "gorm.io/gorm"
    "log"
    "time"
)

type Setting struct {
    dao       *Dao
    Name      string `gorm:"primaryKey"`
    Value     string
    UpdatedAt time.Time
    CreatedAt time.Time
}

func (Setting) TableName() string {
    return "setting"
}

func GetSetting(name string) (Setting, error) {
    setting := Setting{}
    result := database.GetDB().Where("name=?", name).Take(&setting)
    if result.Error == gorm.ErrRecordNotFound {
        return setting, result.Error
    }
    return setting, nil
}

func GetSettingValue(name string, defaultValue string) string {
    setting, err := GetSetting(name)
    if err == nil {
        return setting.Value
    }
    if err == gorm.ErrRecordNotFound {
        return defaultValue
    }
    log.Fatal("Failed to get setting value: ", err)
    return ""
}

func UpdateSetting(name string, value string) {
    setting, err := GetSetting(name)
    if err == gorm.ErrRecordNotFound {
        setting = Setting{Name: name, Value: value}
        result := database.GetDB().Create(&setting)
        if result.Error != nil {
            log.Fatal("Failed to create a new setting: ", result.Error)
        }
    }
    setting.Value = value
    setting.UpdatedAt = time.Now()
    result := database.GetDB().Save(&setting)
    if result.Error != nil {
        log.Fatal("Failed to update setting: ", result.Error)
    }
}
