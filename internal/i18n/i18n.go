package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed messages/*.json
var messagesFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
)

// Init initializes the i18n system
func Init() error {
	bundle = i18n.NewBundle(language.Korean)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load Korean messages
	koData, err := messagesFS.ReadFile("messages/ko.json")
	if err != nil {
		return fmt.Errorf("failed to read Korean messages: %w", err)
	}
	bundle.MustParseMessageFileBytes(koData, "ko.json")

	// Load English messages
	enData, err := messagesFS.ReadFile("messages/en.json")
	if err != nil {
		return fmt.Errorf("failed to read English messages: %w", err)
	}
	bundle.MustParseMessageFileBytes(enData, "en.json")

	// Initialize localizer with detected language
	lang := detectLanguage()
	localizer = i18n.NewLocalizer(bundle, lang)

	return nil
}

// detectLanguage always returns Korean as this is a Korean law information tool
func detectLanguage() string {
	// Always return Korean - this is a Korean law information tool
	return "ko"
}

// T translates a message
func T(messageID string, data ...map[string]interface{}) string {
	if localizer == nil {
		// Fallback if i18n is not initialized
		return messageID
	}

	var templateData map[string]interface{}
	if len(data) > 0 {
		templateData = data[0]
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})
	if err != nil {
		// Fallback to message ID if translation not found
		return messageID
	}
	return msg
}

// Tf translates a message with formatting (convenience function)
func Tf(messageID string, args ...interface{}) string {
	translated := T(messageID)
	if len(args) > 0 {
		return fmt.Sprintf(translated, args...)
	}
	return translated
}

// GetCurrentLanguage returns the current language code
func GetCurrentLanguage() string {
	return detectLanguage()
}
