package concurrent

import (
	"context"
	"testing"

	"github.com/ovechkin-dm/mockio/v2/mock"
)

type TestService interface {
	DoSomething(ctx context.Context, input string) (string, error)
}

func Test_MockioRaceCondition(t *testing.T) {
	for i := 0; i < 3; i++ {
		i := i
		t.Run("subtest", func(t *testing.T) {
			t.Parallel()

			ctrl := mock.NewMockController(t)
			testService := mock.Mock[TestService](ctrl)

			mock.WhenDouble(testService.DoSomething(
				mock.AnyContext(),
				mock.Any[string](),
			)).ThenReturn("mocked response", nil)

			response, err := testService.DoSomething(context.Background(), "test input")
			if err != nil {
				t.Errorf("Test %d: unexpected error: %v", i, err)
			}
			if response != "mocked response" {
				t.Errorf("Test %d: expected 'mocked response', got %q", i, response)
			}

			_, _ = mock.Verify(testService, mock.Once()).DoSomething(
				mock.AnyContext(),
				mock.Any[string](),
			)
		})
	}
}
