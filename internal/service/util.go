package service

import (
	"database/sql"
	"time"
)

func toNullableString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{Valid: true, String: *s}
}

func toNullableTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Valid: true, Time: *t}
}

func fromNullString(sn sql.NullString) *string {
	if sn.Valid {
		return &sn.String
	}
	return nil
}

func FromNullableTime(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
