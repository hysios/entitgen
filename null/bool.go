// Code generated by entitgen. DO NOT EDIT.
package null

import "database/sql"

func NullToBool(nul sql.NullBool) *bool {
	if nul.Valid {
		return &nul.Bool
	}

	return nil
}

func ToNullBool(p *bool) sql.NullBool {
	if p == nil {
		return sql.NullBool{}
	}

	return sql.NullBool{Bool: *p, Valid: true}
}
