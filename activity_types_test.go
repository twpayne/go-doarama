package doarama

import (
	"reflect"
	"testing"
)

func mustFind(at ActivityType, ok bool) ActivityType {
	if !ok {
		panic("activity type not found")
	}
	return at
}

func TestActivityTypesFind(t *testing.T) {
	activityTypes := DefaultActivityTypes
	for _, tc := range []struct {
		s       string
		wantAT  ActivityType
		wantErr error
	}{
		{
			s:      "Swim",
			wantAT: mustFind(activityTypes.FindByID(23)),
		},
		{
			s:      "29",
			wantAT: mustFind(activityTypes.FindByName("Fly - Paraglide")),
		},
		{
			s:       "glide",
			wantErr: &ErrAmbiguousActivityType{},
		},
		{
			s:       "spelunk",
			wantErr: &ErrUnknownActivityType{},
		},
	} {
		gotAT, gotErr := activityTypes.Find(tc.s)
		if gotErr == nil && tc.wantErr == nil {
			if gotAT != tc.wantAT {
				t.Errorf("activityTypes.Find(%q) == %v, nil, want %v, nil", tc.s, gotAT, tc.wantAT)
			}
		} else {
			if reflect.TypeOf(gotErr) != reflect.TypeOf(tc.wantErr) {
				t.Errorf("activityTypes.Find(%q) == ..., %#v, a want ..., %T", tc.s, gotErr, tc.wantErr)
			}
		}
	}
}

func TestDefaultActivityTypes(t *testing.T) {
	for _, at := range []ActivityType{
		{ID: 0, Name: "Undefined - Ground Based"},
		{ID: 17, Name: "Snowboard"},
		{ID: 35, Name: "Fly - Hike + Glide"},
	} {
		if got, ok := DefaultActivityTypes.FindByName(at.Name); !ok || got.ID != at.ID {
			t.Errorf("DefaultActivityTypes.FindByName(%#v) == %#v, %#v, want %#v, true", at.Name, got, ok, at.ID)
		}
		if got, ok := DefaultActivityTypes.FindByID(at.ID); !ok || got.Name != at.Name {
			t.Errorf("DefaultActivityTypes.FindByID(%#v) == %#v, %#v, want %#v, true", at.ID, got, ok, at.Name)
		}
	}
}
