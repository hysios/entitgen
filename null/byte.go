// Code generated by entitgen. DO NOT EDIT.
package null

import "database/sql"

func NullToByte(nul sql.NullByte) *byte {
	if nul.Valid {
		return &nul.Byte
	}

	return nil
}

func ToNullByte(p *byte) sql.NullByte {
	if p == nil {
		return sql.NullByte{}
	}

	return sql.NullByte{Byte: *p, Valid: true}
}
