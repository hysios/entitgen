package model

import (
	"database/sql"
	"testing"

	"gorm.io/datatypes"
)

func TestUser(t *testing.T) {
	var u = &User{
		ID:          0,
		Name:        "test",
		Username:    "",
		Namespace:   "",
		Nickname:    "",
		Email:       "",
		Password:    sql.NullString{},
		Avatar:      "",
		Phone:       "",
		Address:     "",
		Description: "",
		Score:       0,
		Role:        0,
		IsActive:    false,
		InScopes:    datatypes.NewJSONSlice([]string{"hello", "world"}),
	}

	t.Log(u)
}
