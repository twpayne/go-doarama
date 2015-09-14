package doarama

import (
	"testing"
)

func TestDefaultActivityTypes(t *testing.T) {
	for _, at := range []ActivityType{
		{Id: 0, Name: "Undefined - Ground Based"},
		{Id: 17, Name: "Snowboard"},
		{Id: 35, Name: "Fly - Hike + Glide"},
	} {
		if got, ok := DefaultActivityTypes.Id(at.Name); !ok || got != at.Id {
			t.Errorf("DefaultActivityTypes.Id(%#v) == %#v, %#v, want %#v, true", at.Name, got, ok, at.Id)
		}
		if got, ok := DefaultActivityTypes.Name(at.Id); !ok || got != at.Name {
			t.Errorf("DefaultActivityTypes.Name(%#v) == %#v, %#v, want %#v, true", at.Id, got, ok, at.Name)
		}
	}
}
