package doarama

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	for _, tc := range []struct {
		ts Timestamp
		t  time.Time
	}{
		{
			t:  time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			ts: Timestamp(0),
		},
		{
			t:  time.Date(2006, 1, 2, 15, 4, 5, 999000000, time.UTC),
			ts: Timestamp(1136214245999),
		},
		{
			t:  time.Date(2015, 7, 5, 9, 30, 0, 0, time.UTC),
			ts: Timestamp(1436088600000),
		},
	} {
		if got := tc.ts.Time(); !got.Equal(tc.t) {
			t.Errorf("Timestamp(%#v).Time() == %#v, want %#v", tc.ts, got, tc.t)
		}
		if got := NewTimestamp(tc.t); got != tc.ts {
			t.Errorf("NewTimestamp(%#v).Time == %#v, want %#v", tc.t, got, tc.ts)
		}
	}
}
