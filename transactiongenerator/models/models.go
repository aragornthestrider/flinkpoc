package models

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID     string  `json:"id"`
	Source Account `json:"source"`
	Dest   Account `json:"dest"`
	Value  int64   `json:"value"`
	Time   string  `json:"time"`
}
