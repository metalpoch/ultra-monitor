package entity

import (
	"time"
)

type Fat struct {
	ID           int32     `db:"id"`
	IP           string    `db:"ip"`
	Region       string    `db:"region"`
	State        string    `db:"state"`
	Municipality string    `db:"municipality"`
	County       string    `db:"county"`
	Odn          string    `db:"odn"`
	Fat          string    `db:"fat"`
	Bras         string    `db:"bras"`
	Shell        uint8     `db:"shell"`
	Card         uint8     `db:"card"`
	Port         uint8     `db:"port"`
	CreatedAt    time.Time `db:"created_at"`
}

type FatInfoStatus struct {
	ID                 int32     `db:"id"`
	IP                 string    `db:"ip"`
	Region             string    `db:"region"`
	State              string    `db:"state"`
	Municipality       string    `db:"municipality"`
	County             string    `db:"county"`
	Odn                string    `db:"odn"`
	Fat                string    `db:"fat"`
	Bras               string    `db:"bras"`
	Shell              uint8     `db:"shell"`
	Card               uint8     `db:"card"`
	Port               uint8     `db:"port"`
	Plans              string    `db:"plans"`
	CreatedAt          time.Time `db:"created_at"`
	Actives            uint32    `db:"actives"`
	ProvisionedOffline uint32    `db:"provisioned_offline"`
	CutOff             uint32    `db:"cut_off"`
	InProgress         uint32    `db:"in_progress"`
	Date               time.Time `db:"date"`
}

type FatStatus struct {
	FatsID             int32     `db:"fats_id"`
	Date               time.Time `db:"date"`
	Actives            uint32    `db:"actives"`
	ProvisionedOffline uint32    `db:"provisioned_offline"`
	CutOff             uint32    `db:"cut_off"`
	InProgress         uint32    `db:"in_progress"`
	Plans              string    `db:"plans"`
}

type LastFatStatus struct {
	Name               string `db:"name"`
	Bras               string `db:"bras"`
	Actives            uint32 `db:"actives"`
	ProvisionedOffline uint32 `db:"provisioned_offline"`
	CutOff             uint32 `db:"cut_off"`
	InProgress         uint32 `db:"in_progress"`
}

type FatStatusSummary struct {
	Date               time.Time `db:"date"`
	Actives            uint32    `db:"actives"`
	ProvisionedOffline uint32    `db:"provisioned_offline"`
	CutOff             uint32    `db:"cut_off"`
	InProgress         uint32    `db:"in_progress"`
}

type UpsertFat struct {
	IP                 string    `db:"ip"`
	Region             string    `db:"region"`
	State              string    `db:"state"`
	Municipality       string    `db:"municipality"`
	County             string    `db:"county"`
	Odn                string    `db:"odn"`
	Fat                string    `db:"fat"`
	Bras               string    `db:"bras"`
	Shell              uint8     `db:"shell"`
	Card               uint8     `db:"card"`
	Port               uint8     `db:"port"`
	Plans              string    `db:"plans"`
	Actives            uint32    `db:"actives"`
	ProvisionedOffline uint32    `db:"provisioned_offline"`
	CutOff             uint32    `db:"cut_off"`
	InProgress         uint32    `db:"in_progress"`
	Date               time.Time `db:"date"`
}
