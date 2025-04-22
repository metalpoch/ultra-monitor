package entity

type Location struct {
	ID           uint   `gorm:"primaryKey"`
	State        string `gorm:"uniqueIndex:idx_unique_location"`
	County       string `gorm:"uniqueIndex:idx_unique_location"`
	Municipality string `gorm:"uniqueIndex:idx_unique_location"`
}
