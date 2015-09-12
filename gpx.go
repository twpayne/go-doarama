package doarama

import (
	"fmt"
	"io"
	"time"
)

// WriteGPX writes samples to w in GPX format.
func WriteGPX(w io.Writer, samples []Sample) error {
	if _, err := fmt.Fprintf(w, ""+
		"<gpx version=\"1.1\" creator=\"https://github.com/twpayne/go-doarama\">"+
		"<trk>"+
		"<trkseg>"); err != nil {
		return err
	}
	for _, s := range samples {
		if _, err := fmt.Fprintf(w, ""+
			"<trkpt lat=\"%f\" lon=\"%f\">"+
			"<ele>%f</ele>"+
			"<time>%s</time>"+
			"</trkpt>",
			s.Coords.Latitude, s.Coords.Longitude,
			s.Coords.Altitude,
			time.Unix(s.Time/1000, s.Time%1000).UTC().Format("2006-01-02T15:04:05Z")); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(w, ""+
		"</trkseg>"+
		"</trk>"+
		"</gpx>"); err != nil {
		return err
	}
	return nil
}
