package utils

import "time"

func Int64PrtToTimePtr(t *int64) *time.Time {
	if t == nil {
		return nil
	}
	tm := time.Unix(*t, 0)
	return &tm
}
