package response

type TwilioResponse struct {
	From string `json:"from"`
	To   string `json:"to"`
	Body string `json:"body"`
}
