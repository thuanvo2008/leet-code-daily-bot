package model

type LeetCodeResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	ActiveDailyCodingChallengeQuestion ActiveDailyCodingChallengeQuestion `json:"activeDailyCodingChallengeQuestion"`
}

type ActiveDailyCodingChallengeQuestion struct {
	Date     string   `json:"date"`
	Question Question `json:"question"`
}

type Question struct {
	Title      string `json:"title"`
	TitleSlug  string `json:"titleSlug"`
	Difficulty string `json:"difficulty"`
	Url        string `json:"url,omitempty"` // Optional field, not in original response
	Content    string `json:"content"`
}
