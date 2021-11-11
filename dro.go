package main

type DaemonDro struct {
	Name    string `gorm:"primaryKey"`
	Path    string `gorm:"not null;default:null"`
	Enabled bool
}
