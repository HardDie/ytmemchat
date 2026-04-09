package chat

// WebsocketPayload defines the JSON structure sent to the browser.
type WebsocketPayload struct {
	AuthorName    string `json:"authorName"`
	AuthorPicture string `json:"authorPicture"`
	MessageText   string `json:"messageText"`
	PublishedAt   string `json:"publishedAt"`
	IsModerator   bool   `json:"isModerator"`
	IsOwner       bool   `json:"isOwner"`
}
