package presenter

type Event struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	StartedAt   string `json:"startedAt"`
	FinishedAt  string `json:"finishedAt"`
	Description string `json:"description"`
	NotifyAt    string `json:"notifyAt"`
}
