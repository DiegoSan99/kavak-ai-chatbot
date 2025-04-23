package request

type TwilioRequest struct {
	From string `form:"From"`
	To   string `form:"To"`
	Body string `form:"Body"`
}
