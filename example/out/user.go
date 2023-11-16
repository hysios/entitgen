package out

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
)

type Member struct {
	ID           uint
	UserID       uint
	EnterpriseID uint
	User         *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

//go:generate entitgen -type User -O no_models=User
type User struct {
	ID          uint
	Name        string
	Username    string
	Namespace   string
	Nickname    string
	Email       string
	Password    sql.NullString
	Avatar      string
	Phone       string
	Address     string
	Description string
	Score       float64
	Role        int32
	IsActive    bool
	InScopes    datatypes.JSONSlice[string]
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MemberID    uint
	Member      *Member
	Friends     []*Friend
}
