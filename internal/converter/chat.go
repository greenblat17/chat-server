package converter

import (
	"github.com/greenblat17/chat-server/internal/model"
	desc "github.com/greenblat17/chat-server/pkg/chat_v1"
)

// ToChatFromCreateAPI converts a proto  to a chat model
func ToChatFromCreateAPI(req *desc.CreateRequest) *model.Chat {
	return &model.Chat{
		ChatName:  req.GetChatName(),
		Usernames: req.GetUsernames(),
	}
}
