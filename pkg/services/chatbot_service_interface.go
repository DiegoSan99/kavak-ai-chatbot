package services

type IChatbotService interface {
	QueryChatbot(query string) (string, error)
}
