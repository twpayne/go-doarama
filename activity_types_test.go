package doarama

import (
	"testing"
)

func TestDefaultActivityTypes(t *testing.T) {
	for _, at := range []ActivityType{
		{ID: 0, Name: "Undefined - Ground Based"},
		{ID: 17, Name: "Snowboard"},
		{ID: 35, Name: "Fly - Hike + Glide"},
	} {
		if got, ok := DefaultActivityTypes.ID(at.Name); !ok || got != at.ID {
			t.Errorf("DefaultActivityTypes.ID(%#v) == %#v, %#v, want %#v, true", at.Name, got, ok, at.ID)
		}
		if got, ok := DefaultActivityTypes.Name(at.ID); !ok || got != at.Name {
			t.Errorf("DefaultActivityTypes.Name(%#v) == %#v, %#v, want %#v, true", at.ID, got, ok, at.Name)
		}
	}
}
