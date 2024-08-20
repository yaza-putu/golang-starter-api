package i18n

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	l "golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

type (
	unitTestSuite struct {
		suite.Suite
	}
	testTable struct {
		name   string
		data   any
		expect any
	}
)

func (s *unitTestSuite) SetupSuite() {
	bundle := i18n.NewBundle(l.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	_, b, _, _ := runtime.Caller(0)
	Root := filepath.Join(filepath.Dir(b), "../..")

	bundle.LoadMessageFile(fmt.Sprintf("%s/locales/en.yaml", Root))
	bundle.LoadMessageFile(fmt.Sprintf("%s/locales/id.yaml", Root))

	Bundle = bundle
}

func TestSuite(t *testing.T) {
	suite.Run(t, &unitTestSuite{})
}

func (s *unitTestSuite) TestRunMethodNew() {

	testTable := []testTable{
		{
			name:   "without parameter",
			data:   nil,
			expect: config.App().Lang,
		},
		{
			name:   "change locale",
			data:   "ja",
			expect: "ja",
		},
	}

	for _, t := range testTable {
		s.Run(t.name, func() {
			if t.data != nil {
				New(SetLocale(t.data.(string)))
			} else {
				New()
			}

			assert.Equal(s.T(), t.expect, lang.locale)
		})
	}
}

func (s *unitTestSuite) TestTranslate() {
	testTable := []testTable{
		{
			name: "to english",
			data: map[string]any{
				"locale": "en",
				"key":    "greeting",
			}, // key
			expect: "Welcome",
		},
		{
			name: "nested key",
			data: map[string]any{
				"locale": "en",
				"key":    "validations.badrequest",
			}, // key
			expect: "Unprocessable Content",
		},
		{
			name: "to indonesian",
			data: map[string]any{
				"locale": "id",
				"key":    "greeting",
			}, // key
			expect: "Selamat Datang",
		},
		{
			name: "not found the key",
			data: map[string]any{
				"locale": "id",
				"key":    "",
			}, // key
			expect: "",
		},
		{
			name: "unknow locale",
			data: map[string]any{
				"locale": "",
				"key":    "greeting",
			}, // key
			expect: "Welcome",
		},
		{
			name: "uninitialize method new",
			data: map[string]any{
				"locale": "",
				"key":    "greeting",
			}, // key
			expect: "Welcome",
		},
	}

	for _, t := range testTable {
		s.Run(t.name, func() {

			data := t.data.(map[string]any)
			if t.name != "uninitialize method new" {
				New(SetLocale(data["locale"].(string)))
			}

			result := T(
				Localize{
					Key: data["key"].(string),
				},
			)

			assert.Equal(s.T(), t.expect, result)
		})
	}
}

func (s *unitTestSuite) TestTranslateSetLocaleWithoutSetGlobalLocal() {
	New(SetLocale("en"))

	result := T(Localize{
		Key:    "greeting",
		Locale: "id",
	})

	assert.Equal(s.T(), "Selamat Datang", result)
	assert.Equal(s.T(), lang.locale, "en")
}
