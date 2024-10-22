package converter

import (
	"github.com/greenblat17/chat-server/internal/model"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
)

// ToMessageFromAPI converts proto message to message model
func ToMessageFromAPI(req *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		From: req.GetFrom(),
		Text: req.GetText(),
	}
}
