package models

type FinalDecisionPayload struct {
	Year int `json:"year"`
}

type FinalDicisionResponse struct {
	Ok             bool     `json:"ok"`
	UncertainUsers []string `json:"uncertainUsers"`
}
