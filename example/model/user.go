package model

import (
	"database/sql"
	"time"

	"gorm.io/datatypes"
)

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
	// Permissions map[string]*github.com/hysios/entitgen/example/gen/proto.User_PermissionList
	CreatedAt time.Time
	UpdatedAt time.Time
}
