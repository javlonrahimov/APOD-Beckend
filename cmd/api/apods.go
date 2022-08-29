package main

import (
	"fmt"
	"net/http"
	"time"

	"apod.api.javlonrahimov1212/internal/data"
	"apod.api.javlonrahimov1212/internal/validator"
)

func (a *application) createApodHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string         `json:"title"`
		Explanation string         `json:"explanation"`
		HdUrl       string         `json:"hd_url"`
		Url         string         `json:"url"`
		MediaType   data.MediaType `json:"media_type"`
	}

	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResposne(w, r, err)
		return
	}

	apod := &data.Apod{
		Date:        time.Time{},
		Explanation: input.Explanation,
		HdUrl:       input.HdUrl,
		Url:         input.Url,
		Title:       input.Title,
		MediaType:   input.MediaType,
		CreatedAt:   time.Time{},
	}

	v := validator.New()

	if data.ValidateApod(v, apod); !v.Valid() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (a *application) showApodHandler(w http.ResponseWriter, r *http.Request) {

	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	apod := data.Apod{
		ID:          id,
		Date:        time.Now(),
		Explanation: "At the core of the Crab Nebula lies a city-sized, magnetized neutron star spinning 30 times a second. Known as the Crab Pulsar, it is the bright spot in the center of the gaseous swirl at the nebula's core. About twelve light-years across, the spectacular picture frames the glowing gas, cavities and swirling filaments near the Crab Nebula's center.  The featured picture combines visible light from the Hubble Space Telescope in purple, X-ray light from the Chandra X-ray Observatory in blue, and infrared light from the Spitzer Space Telescope in red.  Like a cosmic dynamo the Crab pulsar powers the emission from the nebula, driving a shock wave through surrounding material and accelerating the spiraling electrons. With more mass than the Sun and the density of an atomic nucleus,the spinning pulsar is the collapsed core of a massive star that exploded. The outer parts of the Crab Nebula are the expanding remnants of the star's component gasses. The supernova explosion was witnessed on planet Earth in the year 1054.   Explore Your Universe: Random APOD Generator",
		HdUrl:       "https://apod.nasa.gov/apod/image/2208/Crab_HubbleChandraSpitzer_3600.jpg",
		Url:         "https://apod.nasa.gov/apod/image/2208/Crab_HubbleChandraSpitzer_1080.jpg",
		Title:       "The Spinning Pulsar of the Crab Nebula",
		MediaType:   data.Image,
		CreatedAt:   time.Now(),
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"apod": apod}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
