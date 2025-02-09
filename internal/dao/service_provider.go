package dao

type ServiceProvider struct {
    dao    *Dao
    Id     int    `gorm:"primaryKey;autoIncrement"`
    Name   string `gorm:"uniqueIndex:idx_name"`
    Url    string
    ApiKey string
}

func (ServiceProvider) TableName() string {
    return "service_provider"
}
