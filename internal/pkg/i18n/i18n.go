package i18n

import (
	"strings"
	texttemplate "text/template"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n/template"
	"github.com/yaza-putu/golang-starter-api/internal/config"
)

var Bundle *i18n.Bundle
var lang *language

type (
	language struct {
		locale    string
		localizer *i18n.Localizer
	}
	Localize struct {
		Key            string
		Locale         string
		TemplateData   any
		PluralCount    any
		Default        *i18n.Message
		Funcs          texttemplate.FuncMap
		TemplateParser template.Parser
	}
	optFunc func(lang *language)
)

func SetLocale(locale string) optFunc {
	return func(lang *language) {
		lang.locale = locale
		lang.localizer = i18n.NewLocalizer(Bundle, locale)
	}
}

func New(opts ...optFunc) {

	o := language{
		locale:    config.App().Lang,
		localizer: i18n.NewLocalizer(Bundle, config.App().Lang),
	}

	for _, fn := range opts {
		fn(&o)
	}

	lang = &o
}

func T(localize Localize) string {

	// avoid the localizer not initialized
	if lang.locale == "" && localize.Locale == "" {
		New()
	}

	if localize.Key == "" {
		return ""
	}

	config := &i18n.LocalizeConfig{
		MessageID:      localize.Key,
		TemplateData:   localize.TemplateData,
		PluralCount:    localize.PluralCount,
		DefaultMessage: localize.Default,
		Funcs:          localize.Funcs,
		TemplateParser: localize.TemplateParser,
	}

	if localize.Locale != "" {
		return i18n.NewLocalizer(Bundle, localize.Locale).MustLocalize(config)
	}

	return lang.localizer.MustLocalize(config)
}

// Locale
// get active locale
func Locale() string {
	if lang.locale != "" {
		return simplyLang(lang.locale)
	}

	return config.App().Lang
}

func simplyLang(locale string) string {
	languages := strings.Split(locale, ",")

	lang := config.App().Lang

	for _, l := range languages {
		l = strings.Split(l, ";")[0]
		l = strings.Split(l, "-")[0]
		l = strings.TrimSpace(l)
		if l != "" {
			lang = l
		}
	}

	return lang
}
