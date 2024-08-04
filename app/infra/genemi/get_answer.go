package genemi

import (
	"context"
	"fmt"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
	"github.com/google/generative-ai-go/genai"
)

func (c *client) GenAnswer(ctx context.Context, q string) (string, error) {
	session := c.model.StartChat()
	setAnswerSessionHistory(session)

	resp, err := session.SendMessage(ctx, genai.Text(q))
	if err != nil {
		return "", errutil.Wrap(err)
	}

	logger.FromCtx(ctx).Debug("GenAnswer", "resp", resp)

	answer := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		answer += fmt.Sprintf("%v\n", part)
	}

	return answer, nil
}

func setAnswerSessionHistory(session *genai.ChatSession) {
	// model.SafetySettings = Adjust safety settings
	// See https://ai.google.dev/gemini-api/docs/safety-settings
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("わからない英語の単語や文章の翻訳と説明をお願いしたいです。\n\nmardown形式で渡すので、太字になっているところを重点的に説明してください。\nすべての項目で改行を挟むようにしてください。\n\n返答は以下のmarkdownフォーマットで回答してください。\n---\n# 日本語訳\n\n[ここに日本語訳を入力]\n\n# 英英辞典の説明\n\n[ここに英英辞典の説明]\n\n# 英英辞典の日本語訳\n\n[ここに英英辞書の日本語訳]\n\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("わかりました。**太字**の単語や文章を重点的に説明し、日本語訳、英英辞典の説明、英英辞典の日本語訳をmarkdown形式で記述します。\n\n**例として、以下の文章を翻訳してください。**\n\nThis is a **pen**. \n---\n# 日本語訳\n\nこれはペンです。\n\n# 英英辞典の説明\n\nA **pen** is a writing instrument used to apply ink to paper, typically in the form of a nib or ballpoint.\n\n# 英英辞典の日本語訳\n\n**ペン**は、通常はニブまたはボールポイントの形をしたインクを紙に塗布するために使用される筆記用具です。 \n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("here you can only enter **as part of** an organized tour group"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("---\n# 日本語訳\n\nここは、団体ツアーの一部としてのみ入場できます。\n\n# 英英辞典の説明\n\n**as part of** means \"included in\" or \"being a component of\".\n\n# 英英辞典の日本語訳\n\n**as part of** は、「～に含まれる」「～の一部である」という意味です。 \n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("英英辞典の説明の部分を\n\n英単語： 英英辞典の説明\n\nの形にしてください"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("了解しました。英英辞典の説明を「英単語： 英英辞典の説明」の形式で記述します。\n\n---\n# 日本語訳\n\nここは、団体ツアーの一部としてのみ入場できます。\n\n# 英英辞典の説明\n\n**as part of**:  included in or being a component of\n\n# 英英辞典の日本語訳\n\n**as part of**:  ～に含まれる、～の一部である \n"),
			},
		},
	}
}
