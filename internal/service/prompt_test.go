package service

import (
	"html/template"
	"testing"

	"github.com/fukaraca/skypiea/internal/storage"
)

func TestBuildPrompt(t *testing.T) {
	aboutMe := "My name is Karl, working on a project that involves different tools and tech like Go, K8s, Aws, Gemini etc.."

	summary := "User previously discussed prompt handling for Gemini API."
	history := make([]*storage.Message, 0)
	q1 := template.HTML("I want to implement my own memory functionality for the chat.")
	history = append(history, &storage.Message{ByUser: true, MessageText: &q1})
	q2 := template.HTML("You can summarize conversations after N messages...")
	history = append(history, &storage.Message{ByUser: false, MessageText: &q2})
	q3 := template.HTML("What is the best low-cost method?")
	history = append(history, &storage.Message{ByUser: true, MessageText: &q3})
	qNow := template.HTML("How do I manage memory usage in a 1GB k3s node?")
	prompt := BuildPrompt(&aboutMe, &summary, history, &storage.Message{ByUser: true, MessageText: &qNow})
	t.Log(prompt)
}
