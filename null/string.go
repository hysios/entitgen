// Code generated by entitgen. DO NOT EDIT.
package null

import "database/sql"

func NullToString(nul sql.NullString) *string {
	if nul.Valid {
		return &nul.String
	}

	return nil
}

func ToNullString(p *string) sql.NullString {
	if p == nil {
		return sql.NullString{}
	}

	return sql.NullString{String: *p, Valid: true}
}
