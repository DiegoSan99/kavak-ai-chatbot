package web

import (
	"fmt"
	"net/http"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/services"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/utils"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/web/dto/request"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ChatbotController struct {
	logger         *zap.SugaredLogger
	chatbotService services.IChatbotService
}

func NewChatbotController(e *echo.Echo, logger *zap.SugaredLogger, chatbotService services.IChatbotService) *ChatbotController {
	c := &ChatbotController{
		logger:         logger,
		chatbotService: chatbotService,
	}

	api := e.Group("/api/v1/chatbot")

	api.POST("/", c.QueryChatbot)

	return c
}

func (c *ChatbotController) QueryChatbot(ctx echo.Context) error {
	var twilioRequest request.TwilioRequest
	if err := ctx.Bind(&twilioRequest); err != nil {
		c.logger.Error("Error binding query: ", zap.Error(err))
		return ctx.XML(http.StatusBadRequest, "<Response><Message>Error processing request</Message></Response>")
	}

	responseText, err := c.chatbotService.QueryChatbot(twilioRequest.Body)
	if err != nil {
		c.logger.Error("Error getting chatbot response: ", err)
		return ctx.XML(http.StatusInternalServerError, "<Response><Message>Internal server error</Message></Response>")
	}

	twiml := fmt.Sprintf(`<Response><Message>%s</Message></Response>`, utils.EscapeXML(responseText))

	// Twilio expects XML, not JSON
	return ctx.XMLBlob(http.StatusOK, []byte(twiml))
}
