package resolver

import (
	"context"
	"testing"

	"github.com/RyoheiTomiyama/phraze-api/router/graph/model"
	"github.com/RyoheiTomiyama/phraze-api/test/assertion"
	"github.com/RyoheiTomiyama/phraze-api/util/errutil"
	"github.com/stretchr/testify/assert"
)

func (s *resolverSuite) TestReviewCard() {
	ctx := context.Background()

	s.T().Run("Validationエラー", func(t *testing.T) {
		testCases := []struct {
			name   string
			input  model.ReviewCardInput
			assert func(err error)
		}{
			{
				name:  "CardIDが0値の場合",
				input: model.ReviewCardInput{Grade: 3},
				assert: func(err error) {
					assertion.AssertError(t, "CardIDは必須項目です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Gradeが0の場合",
				input: model.ReviewCardInput{CardID: 100, Grade: 0},
				assert: func(err error) {
					assertion.AssertError(t, "Gradeは1が最小です", errutil.CodeBadRequest, err)
				},
			},
			{
				name:  "Gradeが5より大きい場合",
				input: model.ReviewCardInput{CardID: 100, Grade: 6},
				assert: func(err error) {
					assertion.AssertError(t, "Gradeは5が最大です", errutil.CodeBadRequest, err)
				},
			},
		}

		for _, tc := range testCases {
			result, err := s.resolver.Mutation().ReviewCard(ctx, tc.input)
			assert.Nil(t, result)
			tc.assert(err)
		}
	})
}
