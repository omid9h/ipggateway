package repo

type Terminal struct {
	Terminal string `gorm:"not null" json:"terminal"`
	Addr     string `gorm:"not null" json:"addr"`
}
