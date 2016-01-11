package doarama_test

import (
	"log"
	"os"
	"path/filepath"

	"github.com/twpayne/go-doarama"
)

func Example() (*doarama.Visualisation, error) {
	// Create the client using anonymous authentication
	apiKey := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	apiName := "Your API Name"
	userID := "userid"
	client := doarama.NewClient(doarama.APIURL, apiKey, apiName).Anonymous(userID)

	// Open the GPS track
	filename := "activity GPS filename (GPX or IGC)"
	gpsTrack, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer gpsTrack.Close()

	// Create the activity
	activity, err := client.CreateActivity(filepath.Base(filename), gpsTrack)
	if err != nil {
		return nil, err
	}
	log.Printf("ActivityId: %d", activity.ID)

	// Set the activity info
	if err := activity.SetInfo(&doarama.ActivityInfo{
		TypeID: doarama.FlyParaglide,
	}); err != nil {
		return nil, err
	}

	// Create the visualisation
	activities := []*doarama.Activity{activity}
	visualisation, err := client.CreateVisualisation(activities)
	if err != nil {
		return nil, err
	}
	log.Printf("VisualisationKey: %s", visualisation.Key)
	log.Printf("VisualisationURL: %s", visualisation.URL(nil))

	return visualisation, nil
}
