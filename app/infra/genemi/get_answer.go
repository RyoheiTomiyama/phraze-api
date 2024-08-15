package genemi

import (
	"context"
	"fmt"
	"strings"

	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/RyoheiTomiyama/phraze-api/util/logger"
	"github.com/google/generative-ai-go/genai"
)

func (c *client) GenAnswer(ctx context.Context, q string) (string, error) {
	log := logger.FromCtx(ctx)

	session := c.model.StartChat()
	setAnswerSessionHistory(session)

	resp, err := session.SendMessage(ctx, genai.Text(q))
	if err != nil {
		return "", errutil.Wrap(err)
	}

	log.Debug(ctx, "GenAnswer", "q", q, "resp", resp, "Parts", fmt.Sprintf("%+v", resp.Candidates[0].Content.Parts))

	answer := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		answer += fmt.Sprintf("%v\n", part)
	}

	if !strings.HasPrefix(answer, "**【日本語訳】**") {
		err := fmt.Errorf("Geminiから異常な解答を検出しました。")
		log.ErrorWithNotify(ctx, err, "q", q, "answer", answer)
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
				//nolint:lll
				genai.Text("markdown形式で入力された英語の**単語または文章**に対して、以下の処理を行ってください。\n入力されたテキストの太字のところを重点的に説明してください。\nすべての項目で改行を挟むようにしてください。\n英英辞典の説明の部分は\n英単語： 英英辞典の説明\nの形にしてください。\n返答は以下のmarkdownフォーマットで回答してください。\n\n**【日本語訳】**  \n[ここに日本語訳を入力]\n\n**【英英辞典の説明】**  \n[ここに英英辞典の説明]\n\n**英英辞典の日本語訳】**  \n[ここに英英辞書の日本語訳]"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				//nolint:lll
				genai.Text("かしこまりました。markdown形式で単語や文章を記載してください。 \n太字になっている箇所を重点的に説明し、日本語訳、英英辞典の説明、英英辞典の日本語訳をmarkdown形式で記述します。 \n\n例：\n\n```markdown\n**This is a sentence.** \n**This** is a **word**.\n```\n\n上記のような形式で単語や文章を記載していただければ、的確な説明を提供できます。 \n"),
			},
		},
		{
			Role: "user",
			Parts: []genai.Part{
				//nolint:lll
				genai.Text("この後入力するものを指定した形式で出力してください。もし英語ではない場合は指定した形式ではなく、自動翻訳に失敗しました。とだけ出力してください。自動翻訳に失敗した場合は理由も教えてください。"),
			},
		},
	}
}
