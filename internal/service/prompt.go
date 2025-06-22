package service

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/fukaraca/skypiea/internal/storage"
	"github.com/google/uuid"
)

func TitleFromString(input string, maxWords, maxLen int) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return fmt.Sprintf("untitled-%s", uuid.New().String()[:8])
	}

	words := strings.Fields(input)
	var selected []string
	currLen := 0

	for i, w := range words {
		if i >= maxWords {
			break
		}

		addLen := len(w)
		if len(selected) > 0 {
			addLen++
		}
		if currLen+addLen > maxLen {
			break
		}
		selected = append(selected, w)
		currLen += addLen
	}

	if len(selected) > 0 {
		return strings.Join(selected, " ")
	}

	// fallback for giant first word: take a prefix of the raw input
	truncated := input
	if len(truncated) > maxLen {
		truncated = truncated[:maxLen]
		truncated = strings.TrimSpace(truncated)
	}
	return fmt.Sprintf("%s-%s", truncated, uuid.New().String()[:8])
}

func (s *Service) Sanitize(txt string, safe bool) *template.HTML {
	var sb strings.Builder
	var err error
	if safe {
		err = s.MDSafe.Convert([]byte(s.policy.Sanitize(txt)), &sb)
	} else {
		err = s.MDUnsafe.Convert([]byte(s.policy.Sanitize(txt)), &sb) // trust the answer from gemini
	}
	if err != nil {
		return nil
	}
	temp := template.HTML(sb.String()) //nolint:gosec
	return &temp
}

func BuildPrompt(aboutMe *string, summary *string, conversations []*storage.Message, newQuestion *storage.Message) string {
	var b strings.Builder

	b.WriteString("You are an assistant helping the user with general questions.\n\n")

	if aboutMe != nil && strings.TrimSpace(*aboutMe) != "" {
		b.WriteString("## About the user:\n")
		b.WriteString(*aboutMe)
		b.WriteString("\n\n")
	}

	if summary != nil && strings.TrimSpace(*summary) != "" {
		b.WriteString("## Background summary of the conversation so far from descending by time :\n")
		b.WriteString(*summary)
		b.WriteString("\n\n")
	}

	b.WriteString("## Recent conversation:\n")
	for _, msg := range conversations {
		if msg.ID == newQuestion.ID {
			continue
		}
		role := func() string {
			if msg.ByUser {
				return "user"
			}
			return "assistant"
		}()
		msgStr := func() string {
			if msg.MessageText != nil {
				return string(*msg.MessageText)
			}
			return ""
		}()
		b.WriteString(fmt.Sprintf("%s: %s\n", role, msgStr))
	}
	b.WriteString("\n")
	q := func() string {
		if newQuestion.MessageText != nil {
			return string(*newQuestion.MessageText)
		}
		return ""
	}()
	b.WriteString("## Now the user asked:\n")
	b.WriteString(q)
	b.WriteString("\n\n")

	b.WriteString("Please respond based on the user’s background and the conversation history above.\n")

	return b.String()
}
