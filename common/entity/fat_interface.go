package entity

type FatInterface struct {
	ID          uint      `gorm:"primaryKey"`
	FatID       uint      `gorm:"uniqueIndex:idx_unique_fat_interface"`
	InterfaceID uint      `gorm:"uniqueIndex:idx_unique_fat_interface"`
	Fat         Fat       `gorm:"constraint:OnDelete:SET NULL"`
	Interface   Interface `gorm:"constraint:OnDelete:CASCADE"`
}
