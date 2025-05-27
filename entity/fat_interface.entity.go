package entity

type FatInterface struct {
	ID    uint64 `db:"id"`
	FatID uint   `db:"fat_id"`
	PonID uint64 `db:"pon_id"`
}
