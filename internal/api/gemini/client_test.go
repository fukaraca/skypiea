package gemini

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/genai"
)

func TestAskToGemini(t *testing.T) {

	cfg := &genai.ClientConfig{
		APIKey:  "test",
		Backend: genai.BackendGeminiAPI,
	}
	client, err := NewClient(&Config{cfg, nil})
	require.NoError(t, err)
	require.NotNil(t, client)
	/*	var wg sync.WaitGroup
		wg.Add(5)
		ctx := context.Background()
		for i := 0; i < 5; i++ {
			go func() {
				str := fmt.Sprintf("%d * 5 = ? do the math", i)
				resp, err := client.AskToGemini(ctx, str, "gemini-2.0-flash-lite")
				if err != nil {
					t.Errorf("goroutine %d: error asking Gemini: %v", i, err)
					return
				}
				t.Logf("goroutine %d: response: %v", i, resp)
				wg.Done()
			}()
		}
		wg.Wait()*/
}
