package request

import (
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transalation_en "github.com/go-playground/validator/v10/translations/en"
	transalation_id "github.com/go-playground/validator/v10/translations/id"
	"github.com/labstack/gommon/log"
	"strings"
	"yaza/src/config"
	"yaza/src/database"
	"yaza/src/http/response"
)

// Validation form
func Validation(s any, msg map[string]string) (response.DataApi, error) {
	uni := ut.New(id.New(), en.New(), id.New())
	trans, found := uni.GetTranslator(config.App().Lang)

	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()
	switch config.App().Lang {
	case "id":
		if err := transalation_id.RegisterDefaultTranslations(v, trans); err != nil {
			log.Fatal(err)
			return response.Api(response.SetStatus(false), response.SetMessage(err)), err
		}
		_ = v.RegisterTranslation("unique", trans, func(ut ut.Translator) error {
			return ut.Add("unique", "{0} ini sudah terdaftar di database", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("unique", fe.Field())
			return t
		})
		break
	default:
		if err := transalation_en.RegisterDefaultTranslations(v, trans); err != nil {
			log.Fatal(err)
			return response.Api(response.SetStatus(false), response.SetMessage(err)), err
		}
		_ = v.RegisterTranslation("unique", trans, func(ut ut.Translator) error {
			return ut.Add("unique", "The {0} already exists in database", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("unique", fe.Field())
			return t
		})
		break
	}

	err := v.RegisterValidation("unique", unique)

	if err != nil {
		return response.Api(response.SetStatus(false), response.SetMessage(err)), err
	}

	for k, ms := range msg {
		rError := v.RegisterTranslation(k, trans, func(ut ut.Translator) error {
			return ut.Add(k, fmt.Sprintf("{0} %s", ms), true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(k, fe.Field())
			return t
		})

		if rError != nil {
			return response.Api(response.SetStatus(false), response.SetMessage(rError)), rError
		}

	}

	r := v.Struct(s)
	m := map[string]interface{}{}

	if r != nil {
		for _, r := range r.(validator.ValidationErrors) {
			m[strings.ToLower(r.Field())] = []string{r.Translate(trans)}
		}

		return response.Api(response.SetStatus(false), response.SetCode(422), response.SetMessage(m)), r
	}

	return response.Api(response.SetCode(200)), nil
}

func unique(fl validator.FieldLevel) bool {
	var count int64
	param := fl.Param()
	params := strings.Split(param, ":")
	switch len(params) {
	case 1:
		return true
	case 2:
		dField := fl.Field().String()

		database.Instance.Table(params[0]).Where(fmt.Sprintf("%s = ?", strings.ToLower(params[1])), dField).Count(&count)

		if count > 0 {
			return false
		} else {
			return true
		}
	case 3:
		ignore := fl.Parent().FieldByName(params[2]).String()
		dField := fl.Field().String()
		database.Instance.Table(params[0]).Where(fmt.Sprintf("%s = ?", strings.ToLower(params[1])), dField).Not(map[string]any{fmt.Sprintf("%s", strings.ToLower(params[2])): []string{ignore}}).Count(&count)
		if count > 0 {
			return false
		} else {
			return true
		}
	default:
		return true
	}
}
