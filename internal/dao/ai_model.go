package dao

import (
    "Tai/internal/database"
    "time"
)

type AiModel struct {
    dao        *Dao
    Id         int    `gorm:"primaryKey;autoIncrement"`
    Name       string `gorm:"index:idx_name"`
    ProviderId int    `gorm:"index:idx_name"`
    CreatedAt  time.Time
}

func (AiModel) TableName() string {
    return "ai_model"
}

func GetAiModelById(id int) *AiModel {
    aiModel := &AiModel{}
    database.GetDB().Where("id = ?", id).First(aiModel)
    return aiModel
}
