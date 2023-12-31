// Code generated by entitgen. DO NOT EDIT.
package null

import "database/sql"

func NullToInt32(nul sql.NullInt32) *int32 {
	if nul.Valid {
		return &nul.Int32
	}

	return nil
}

func ToNullInt32(p *int32) sql.NullInt32 {
	if p == nil {
		return sql.NullInt32{}
	}

	return sql.NullInt32{Int32: *p, Valid: true}
}
