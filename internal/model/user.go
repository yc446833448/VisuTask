package model

type User struct {
	ID            string `gorm:"primaryKey" json:"id"`
	VIPLevel      int    `gorm:"default:0" json:"vipLevel"`
	MaxConcurrent int    `gorm:"default:3" json:"maxConcurrent"`
	Balance       float64 `gorm:"default:0" json:"balance"`
}

// GetMaxConcurrent returns the concurrency limit based on VIP level
func (u *User) GetMaxConcurrent() int {
	if u.MaxConcurrent > 0 {
		return u.MaxConcurrent
	}
	switch u.VIPLevel {
	case 0:
		return 3
	case 1:
		return 5
	case 2:
		return 6
	case 3:
		return 7
	case 4:
		return 8
	case 5:
		return 10
	default:
		return 3
	}
}
