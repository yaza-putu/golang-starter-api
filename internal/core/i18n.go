package core

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	i18n2 "github.com/yaza-putu/golang-starter-api/internal/pkg/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func I18n() {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "../..")

	bundle := i18n.NewBundle(defaultLang())
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	bundle.LoadMessageFile(fmt.Sprintf("%s/internal/locales/en.yaml", root))
	bundle.LoadMessageFile(fmt.Sprintf("%s/internal/locales/id.yaml", root))

	i18n2.Bundle = bundle
	i18n2.New()
}

func defaultLang() language.Tag {
	switch config.App().Lang {
	case "id":
		return language.Indonesian
	default:
		return language.English
	}
}
