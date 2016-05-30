package doarama

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/twpayne/go-gpx"
)

// WriteGPX writes samples to w in GPX format.
func WriteGPX(w io.Writer, samples []Sample) error {
	var trkPt []*gpx.WptType
	for _, s := range samples {
		trkPt = append(trkPt, &gpx.WptType{
			Lat:  s.Coords.Latitude,
			Lon:  s.Coords.Longitude,
			Ele:  s.Coords.Altitude,
			Time: s.Time.Time(),
		})
	}
	t := gpx.T{
		Version: "1.1",
		Creator: "https://github.com/twpayne/go-doarama",
		Trk: []*gpx.TrkType{
			&gpx.TrkType{
				TrkSeg: []*gpx.TrkSegType{
					&gpx.TrkSegType{
						TrkPt: trkPt,
					},
				},
			},
		},
	}
	return xml.NewEncoder(w).EncodeElement(t, gpx.StartElement)
}

// dmmh splits x into degrees, milliminutes, and a hemisphere.
// hs should be "NS" for latitude and "EW" for longitude.
func dmmh(x float64, hs string) (d int, mm int, h uint8) {
	if x < 0 {
		h = hs[1]
		x = -x
	} else {
		h = hs[0]
	}
	d = int(x)
	mm = int(60000 * (x - float64(d)))
	return
}

// WriteIGC writes samples to w in IGC format.
func WriteIGC(w io.Writer, samples []Sample) error {
	var date time.Time
	for _, s := range samples {
		t := s.Time.Time()
		if t.Day() != date.Day() || t.Month() != date.Month() || t.Year() != date.Year() {
			date = t
			if _, err := fmt.Fprintf(w, "HFDTE%02d%02d%02d\r\n", date.Day(), date.Month(), date.Year()%100); err != nil {
				return err
			}
		}
		latDeg, latMMin, latHemi := dmmh(s.Coords.Latitude, "NS")
		lngDeg, lngMMin, lngHemi := dmmh(s.Coords.Longitude, "EW")
		if _, err := fmt.Fprintf(w, "B%02d%02d%02d%02d%05d%c%03d%05d%cA%05d%05d\r\n",
			t.Hour(), t.Minute(), t.Second(),
			latDeg, latMMin, latHemi,
			lngDeg, lngMMin, lngHemi,
			int(s.Coords.Altitude), int(s.Coords.Altitude)); err != nil {
			return err
		}
	}
	return nil
}
