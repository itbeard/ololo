package main

import (
	"strings"
)

type ContentModerator interface {
	IsForbiddenContent(text string) bool
}

type SimpleDrugModerator struct {
	keywords []string
}

func NewSimpleDrugModerator() *SimpleDrugModerator {
	return &SimpleDrugModerator{
		keywords: []string{"drugs", "cocaine", "heroin", "weed", "meth"},
	}
}

func (m *SimpleDrugModerator) IsForbiddenContent(text string) bool {
	lowercaseText := strings.ToLower(text)
	for _, keyword := range m.keywords {
		if strings.Contains(lowercaseText, keyword) {
			return true
		}
	}
	return false
}
