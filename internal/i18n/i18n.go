package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed messages/*.json
var messagesFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	langFlag  string // Language set by --lang flag
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

// SetLanguage sets the language for the application
func SetLanguage(lang string) {
	langFlag = lang
	localizer = i18n.NewLocalizer(bundle, lang)
}

// detectLanguage detects the user's preferred language
func detectLanguage() string {
	// Priority 1: --lang flag (set by SetLanguage)
	if langFlag != "" {
		return langFlag
	}

	// Priority 2: LANG environment variable
	if envLang := os.Getenv("LANG"); envLang != "" {
		lang := parseLocale(envLang)
		if lang != "" {
			return lang
		}
	}

	// Priority 3: LC_ALL environment variable
	if lcAll := os.Getenv("LC_ALL"); lcAll != "" {
		lang := parseLocale(lcAll)
		if lang != "" {
			return lang
		}
	}

	// Default: Korean
	return "ko"
}

// parseLocale extracts the language code from a locale string
// Handles formats like: en_US.UTF-8, en-US.UTF-8, ko_KR.UTF-8, C.UTF-8
func parseLocale(locale string) string {
	// Handle special case for C locale
	if locale == "C" || strings.HasPrefix(locale, "C.") {
		return "" // Use default
	}

	// Remove encoding suffix (e.g., ".UTF-8")
	if idx := strings.Index(locale, "."); idx != -1 {
		locale = locale[:idx]
	}

	// Replace '-' with '_' for consistency
	locale = strings.ReplaceAll(locale, "-", "_")

	// Split on '_' and take the first part
	parts := strings.Split(locale, "_")
	if len(parts) > 0 {
		lang := strings.ToLower(parts[0])
		// Only return if it's a supported language
		if lang == "ko" || lang == "en" {
			return lang
		}
	}

	return "" // Use default
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
	if langFlag != "" {
		return langFlag
	}
	return detectLanguage()
}
