package request

import (
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transalation_id "github.com/go-playground/validator/v10/translations/id"
	"github.com/labstack/gommon/log"
	"strings"
	"yaza/src/http/response"
)

// Validation form
func Validation(s any) (response.DataApi, error) {
	translator := id.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("id")

	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := transalation_id.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	r := v.Struct(s)
	m := map[string]interface{}{}

	if r != nil {
		for _, r := range r.(validator.ValidationErrors) {
			m[strings.ToLower(r.Field())] = []string{r.Translate(trans)}
		}

		return response.Api(response.SetStatus(false), response.SetCode(422), response.SetMessage(m)), r
	}

	return response.Api(response.SetStatus(false), response.SetCode(422)), nil
}
