package notificationserviceimpl

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/im/telegram"
)

func (s *service) telegramSendMsg(
	ctx context.Context,
	setting *entity.IMTelegram,
	msg string,
) error {
	botToken, err := setting.BotToken.GetPlain()
	if err != nil {
		return apperrors.Wrap(err)
	}
	err = telegram.NewClient().SendMessage(ctx, botToken, setting.ChatID, msg, "HTML")
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}
