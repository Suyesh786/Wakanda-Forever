package main

// Era represents a single historical point
type Era struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	SearchTerm string `json:"searchTerm"` // Used for Wikimedia searches
	Desc       string `json:"desc"`
	Year       string `json:"year"`
	Coord      string `json:"coord"`
	ArtName    string `json:"artName"`
	ArtData    string `json:"artData"`
	Color      string `json:"color"`
	Filter     string `json:"filter"`
}

// GetAllEras acts as our database query
func GetAllEras() []Era {
	return []Era{
		{
			ID:         "egypt",
			Title:      "Giza Plateau",
			SearchTerm: "Great_Pyramid_of_Giza",
			Desc:       "Walk beneath the Great Pyramid during construction.",
			Year:       "2560 BCE",
			Coord:      "LOC: 29.9792° N, 31.1342° E",
			ArtName:    "The Khufu Bark",
			ArtData:    "A full-size solar vessel intended for use in the afterlife.",
			Color:      "#00d4ff",
			Filter:     "empire",
		},
		{
			ID:         "rome",
			Title:      "Roman Forum",
			SearchTerm: "Roman_Forum",
			Desc:       "The nerve center of the Mediterranean.",
			Year:       "115 CE",
			Coord:      "LOC: 41.8925° N, 12.4853° E",
			ArtName:    "Augustus Prima Porta",
			ArtData:    "A high-marble statue of Augustus Caesar.",
			Color:      "#ff4757",
			Filter:     "empire",
		},
	}
}