package bot

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log/level"

	"cthulhu/telegram"
)

const unbanCommand = "unban"

func (s *service) handleUnban(ctx context.Context, updateReq *telegram.Update) error {
	if updateReq.Message.ReplyToMessage == nil {
		level.Info(s.Logger).Log("msg", "no message quoted")
		return nil
	}

	var (
		chatID   = updateReq.Message.Chat.ID
		authorID = updateReq.Message.From.ID
		userID   = updateReq.Message.ReplyToMessage.From.ID
	)

	level.Info(s.Logger).Log(
		"msg", "received new unban request",
		"chat_id", chatID,
		"author_id", authorID,
		"user_id", userID,
	)

	if !s.Config.hasPermissions(chatID, authorID, unbanCommand) {
		level.Info(s.Logger).Log("msg", "not enough privileges")
		return nil
	}
	if err := telegram.UnbanChatMember(ctx, string(s.GetToken()), chatID, userID); err != nil {
		return err
	}
	telegram.SendMessage(ctx, string(s.GetToken()), chatID, fmt.Sprintf("user %s unbanned", updateReq.Message.ReplyToMessage.From.UserName))
	return nil
}
