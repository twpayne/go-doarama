package doarama_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/twpayne/go-doarama"
)

func TestWrite(t *testing.T) {
	for _, tc := range []struct {
		samples []doarama.Sample
		wantGPX string
		wantIGC string
	}{
		{
			samples: []doarama.Sample{
				{
					Time: doarama.NewTimestamp(time.Date(2015, 7, 5, 9, 30, 0, 0, time.UTC)),
					Coords: doarama.Coords{
						Latitude:  47.79885,
						Longitude: 13.04840,
						Altitude:  430,
					},
				},
				{
					Time: doarama.NewTimestamp(time.Date(2015, 7, 5, 11, 15, 0, 0, time.UTC)),
					Coords: doarama.Coords{
						Latitude:  47.80413,
						Longitude: 13.11091,
						Altitude:  1272,
					},
				},
			},
			wantGPX: "" +
				"<gpx version=\"1.1\" creator=\"https://github.com/twpayne/go-doarama\">" +
				"<trk>" +
				"<trkseg>" +
				"<trkpt lat=\"47.798850\" lon=\"13.048400\">" +
				"<ele>430.000000</ele>" +
				"<time>2015-07-05T09:30:00Z</time>" +
				"</trkpt>" +
				"<trkpt lat=\"47.804130\" lon=\"13.110910\">" +
				"<ele>1272.000000</ele>" +
				"<time>2015-07-05T11:15:00Z</time>" +
				"</trkpt>" +
				"</trkseg>" +
				"</trk>" +
				"</gpx>",
			wantIGC: "" +
				"HFDTE050715\r\n" +
				"B0930004747931N01302904EA0043000430\r\n" +
				"B1115004748247N01306654EA0127201272\r\n",
		},
	} {
		bGPX := &bytes.Buffer{}
		if err := doarama.WriteGPX(bGPX, tc.samples); err != nil {
			t.Errorf("doarama.WriteGPX(b, %#v) == %v, want nil", tc.samples, err)
		}
		if bGPX.String() != tc.wantGPX {
			t.Errorf("doarama.WriteGPX(b, %#v) wrote %#v, want %#v", tc.samples, bGPX.String(), tc.wantGPX)
		}
		bIGC := &bytes.Buffer{}
		if err := doarama.WriteIGC(bIGC, tc.samples); err != nil {
			t.Errorf("doarama.WriteIGC(b, %#v) == %v, want nil", tc.samples, err)
		}
		if bIGC.String() != tc.wantIGC {
			t.Errorf("doarama.WriteIGC(b, %#v) wrote %#v, want %#v", tc.samples, bIGC.String(), tc.wantIGC)
		}
	}
}
