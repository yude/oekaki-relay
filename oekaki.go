package main

type OekakiRequest struct {
	Image      string `json:"image"`
	Answer     string `json:"answer"`
	NextAnswer string `json:"next_answer"`
}

type Oekaki struct {
	Answer string `json:"answer"`
	Image  string `json:"image"`
}

type OekakiStore struct {
	Id         string `gorm:"uniqueIndex" json:"id"`
	Answer     string `json:"answer"`
	UserAnswer string `json:"user_answer"`
	Image      string `json:"image"`
}
