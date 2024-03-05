package slash

import (
	"fmt"
	"strings"
	"unicode"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/MoYoez/Lucy_reibot/utils/toolchain"
	rei "github.com/fumiama/ReiBot"
	tgba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	engine := rei.Register("slash", &ctrl.Options[*rei.Ctx]{
		DisableOnDefault:  true,
		Help:              "slash - use / pattern to make it well",
		PrivateDataFolder: "slash",
	})

	engine.OnMessage().SetBlock(false).Handle(func(ctx *rei.Ctx) {
		getReply := QuoteReply(ctx)
		if getReply == "" {
			return
		}
		ctx.Caller.Send(&tgba.MessageConfig{BaseChat: tgba.BaseChat{ChatID: ctx.Message.Chat.ID}, Text: getReply, ParseMode: "MarkdownV2", DisableWebPagePreview: true})
	})

}

func QuoteReply(message *rei.Ctx) (replyMsg string) {
	if len(message.Message.Text) < 2 {
		return ""
	}
	if !strings.HasPrefix(message.Message.Text, "/") || (isASCII(message.Message.Text[:2]) && !strings.HasPrefix(message.Message.Text, "/$")) {
		return ""
	}
	keywords := strings.SplitN(tgba.EscapeText(tgba.ModeMarkdownV2, strings.Replace(message.Message.Text, "$", "", 1)[1:]), " ", 2)
	if len(keywords) == 0 {
		return ""
	}
	senderName := tgba.EscapeText(tgba.ModeMarkdownV2, message.Message.From.FirstName+" "+message.Message.From.LastName)
	senderURI := fmt.Sprintf("tg://user?id=%d", message.Message.From.ID)
	replyToName := ""
	replyToURI := ""
	if message.Message.SenderChat != nil {
		senderName = tgba.EscapeText(tgba.ModeMarkdownV2, message.Message.SenderChat.Title)
		senderURI = fmt.Sprintf("https://t.me/%s", message.Message.SenderChat.UserName)
	}
	if message.Message.ReplyToMessage != nil {
		replyToName = tgba.EscapeText(tgba.ModeMarkdownV2, message.Message.ReplyToMessage.From.FirstName+" "+message.Message.ReplyToMessage.From.LastName)
		replyToURI = fmt.Sprintf("tg://user?id=%d", message.Message.ReplyToMessage.From.ID)

		if message.Message.ReplyToMessage.From.IsBot && len(message.Message.ReplyToMessage.Entities) != 0 {
			if message.Message.ReplyToMessage.Entities[0].Type == "text_mention" {
				replyToName = tgba.EscapeText(tgba.ModeMarkdownV2, message.Message.ReplyToMessage.Entities[0].User.FirstName+" "+message.Message.ReplyToMessage.Entities[0].User.LastName)
				replyToURI = fmt.Sprintf("tg://user?id=%d", message.Message.ReplyToMessage.Entities[0].User.ID)
			}
		}

		if message.Message.ReplyToMessage.SenderChat != nil {
			replyToName = tgba.EscapeText(tgba.ModeMarkdownV2, message.Message.ReplyToMessage.SenderChat.Title)
			replyToURI = fmt.Sprintf("https://t.me/%s", message.Message.ReplyToMessage.SenderChat.UserName)
		}

	} else {
		textNoCommand := strings.TrimPrefix(strings.TrimPrefix(keywords[0], "/"), "$")
		if text := strings.Split(textNoCommand, "@"); len(text) > 1 {
			name := toolchain.GetNickNameFromUsername(text[1])
			if name != "" {
				keywords[0] = text[0]
				replyToName = tgba.EscapeText(tgba.ModeMarkdownV2, name)
				replyToURI = fmt.Sprintf("https://t.me/%s", text[1])
			}
		}
		if replyToName == "" {
			replyToName = "自己"
			replyToURI = senderURI
		}
	}
	least := tgba.EscapeText(tgba.ModeMarkdownV2, "~")
	if len(keywords) < 2 {
		return fmt.Sprintf("[%s](%s) %s了 [%s](%s) %s ", senderName, senderURI, keywords[0], replyToName, replyToURI, least)
	}
	return fmt.Sprintf("[%s](%s) %s [%s](%s) %s %s", senderName, senderURI, keywords[0], replyToName, replyToURI, keywords[1], least)

}

func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}
