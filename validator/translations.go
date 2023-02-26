package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Translations struct {
	Tag        string
	Format     string
	IsOverride bool
	Fuc        func(name string, ut ut.Translator, fe validator.FieldError) string
}

const (
	DefaultFormat = "{0} must {1}"
)

var (
	tra ut.Translator
)

// 使用默认格式
var useDefaultTranslation = []string{"INTEGER", "YYYY", "YYYYMMDD", "YYYYMM", "YYYY-MM", "YYYY-MM-DD", "hh:mm:ss", "YYYY-MM-DD hh:mm:ss"}

// 自定义格式
var translationsArray = []Translations{
	{
		Tag:        "isEmail",
		Format:     "{0} 必须是邮箱格式",
		IsOverride: true,
		Fuc:        defaultTranslation,
	},
}

func init() {
	e := en.New()
	uni := ut.New(e, e)

	// 默认元空间
	tra, _ = uni.GetTranslator("en")
}

func NewDefaultTranslations(tags ...string) (res []Translations) {
	t := Translations{
		Format:     DefaultFormat,
		IsOverride: true,
		Fuc:        defaultTranslation,
	}
	for i := range tags {
		ii := i
		res = append(res, Translations{
			Tag:        tags[ii],
			Format:     t.Format,
			IsOverride: t.IsOverride,
			Fuc:        t.Fuc,
		})
	}
	return
}

func defaultTranslation(tag string, ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(tag, fe.Field(), tag)
	return t
}

func RegisterTranslations(validate *validator.Validate) ut.Translator {
	if err := en_translations.RegisterDefaultTranslations(validate, tra); err != nil {
		panic(err)
	}
	df := NewDefaultTranslations(useDefaultTranslation...)
	translationsArray = append(translationsArray, df...)
	registerTranslation(translationsArray...)
	return tra
}

func registerTranslation(translations ...Translations) {
	for i := range translations {
		ii := i
		if err := validate.RegisterTranslation(translations[ii].Tag, tra, func(ut ut.Translator) error {
			return ut.Add(translations[ii].Tag, translations[ii].Format, translations[ii].IsOverride)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			return translations[ii].Fuc(translations[ii].Tag, ut, fe)
		}); err != nil {
			panic(err)
		}
	}
}
