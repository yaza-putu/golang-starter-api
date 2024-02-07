package request

import (
	"errors"
	"fmt"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transalation_en "github.com/go-playground/validator/v10/translations/en"
	transalation_id "github.com/go-playground/validator/v10/translations/id"
	"github.com/labstack/gommon/log"
	"github.com/yaza-putu/golang-starter-api/internal/config"
	"github.com/yaza-putu/golang-starter-api/internal/database"
	"github.com/yaza-putu/golang-starter-api/internal/http/response"
	"github.com/yaza-putu/golang-starter-api/internal/pkg/logger"
	"strings"
)

type (
	optFunc func(*attr)

	attr struct {
		M map[string]string
	}
)

func defaultParam() attr {
	return attr{
		M: map[string]string{},
	}
}

func CustomMessage(m map[string]string) optFunc {
	return func(p *attr) {
		p.M = m
	}
}

// Validation form
func Validation(s any, opts ...optFunc) (response.DataApi, error) {
	o := defaultParam()

	for _, fn := range opts {
		fn(&o)
	}

	uni := ut.New(id.New(), en.New(), id.New())
	trans, found := uni.GetTranslator(config.App().Lang)

	if !found {
		logger.New(errors.New("translator not found"), logger.SetType(logger.FATAL))
	}

	v := validator.New()
	switch config.App().Lang {
	case "id":
		if err := transalation_id.RegisterDefaultTranslations(v, trans); err != nil {
			log.Fatal(err)
			return response.Api(response.SetMessage(err)), err
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
			logger.New(err, logger.SetType(logger.FATAL))
			return response.Api(response.SetMessage("bad request"), response.SetError(err)), err
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
		return response.Api(response.SetMessage(err)), err
	}
	for k, ms := range o.M {
		rError := v.RegisterTranslation(k, trans, func(ut ut.Translator) error {
			return ut.Add(k, fmt.Sprintf("{0} %s", ms), true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(k, fe.Field())
			return t
		})

		if rError != nil {
			return response.Api(response.SetError(rError), response.SetMessage("Unprocessable Content")), rError
		}

	}

	r := v.Struct(s)
	m := map[string]interface{}{}

	if r != nil {
		for _, r := range r.(validator.ValidationErrors) {
			m[strings.ToLower(r.Field())] = []string{r.Translate(trans)}
		}

		return response.Api(response.SetCode(422), response.SetMessage("Unprocessable Content"), response.SetError(m)), r
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
