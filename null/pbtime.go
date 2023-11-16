package null

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func PbtimeToSQLTime(nul *timestamppb.Timestamp) sql.NullTime {
	if nul != nil {
		return sql.NullTime{Time: nul.AsTime(), Valid: true}
	}

	return sql.NullTime{}
}

func SQLTimeToPbtime(p sql.NullTime) *timestamppb.Timestamp {
	if p.Valid {
		return timestamppb.New(p.Time)
	}

	return nil
}
