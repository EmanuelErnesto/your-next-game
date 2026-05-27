package domain

type Game struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	CoverURL    string   `json:"coverUrl"`
	Status      string   `json:"status"`
	Genre       string   `json:"genre"`
	Platform    string   `json:"platform"`
	Description string   `json:"description"`
	Developer   string   `json:"developer"`
	ReleaseDate string   `json:"releaseDate"`
	HoursPlayed float64  `json:"hoursPlayed"`
	Rating      *float64 `json:"rating"`
	SteamAppID  *string  `json:"steamAppId"`
}
