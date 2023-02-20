package validator

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	validate *validator.Validate
)

func InitValidator() {
	validate = getCustomValidate()
	tra = RegisterTranslations(validate)
}

func AddRegisterVal(tag string, v func(fl validator.FieldLevel) bool, useDefaultTra bool, t ...Translations) error {
	if validate == nil {
		panic(fmt.Errorf("校验器未初始化"))
	}
	if err := validate.RegisterValidation(tag, v); err != nil {
		return err
	}
	// 添加默认的翻译器
	if useDefaultTra {
		arr := NewDefaultTranslations(tag)
		registerTranslation(arr...)
		return nil
	}
	if len(t) == 0 {
		return fmt.Errorf("需要自定义翻译器")
	}
	for i := range t {
		registerTranslation(t[i])
		// 执行第一个就结束
		return nil
	}
	return nil
}

// getCustomValidate 自定义validate
func getCustomValidate() *validator.Validate {
	var err error
	validate = validator.New()
	// 添加校验tag：RFC3339
	if err = validate.RegisterValidation("RFC3339", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		_, err := String2Time(s)
		if err != nil {
			return false
		}
		return true
	}); err != nil {
		panic(err)
	}
	// 添加校验字符串为整数
	if err = validate.RegisterValidation("INTEGER", checkInt); err != nil {
		panic(err)
	}
	// 添加校验tag：YYYY
	if err = validate.RegisterValidation("YYYY", checkYear); err != nil {
		panic(err)
	}
	// 添加校验tag：YYYYMMDD
	if err = validate.RegisterValidation("YYYYMMDD", checkDate); err != nil {
		panic(err)
	}
	// 添加校验tag：YYYYMM
	if err = validate.RegisterValidation("YYYYMM", checkMonth1); err != nil {
		panic(err)
	}
	// 添加校验tag：YYYY-MM
	if err = validate.RegisterValidation("YYYY-MM", checkMonth); err != nil {
		panic(err)
	}
	// 添加校验tag：YYYY-MM-DD
	if err = validate.RegisterValidation("YYYY-MM-DD", checkStrDate); err != nil {
		panic(err)
	}
	// 添加校验tag：hh:mm:ss
	if err = validate.RegisterValidation("hh:mm:ss", validateHMS); err != nil {
		panic(err)
	}
	// 添加校验tag：hhmmss
	if err = validate.RegisterValidation("hhmmss", validateShortHMS); err != nil {
		panic(err)
	}
	// 添加校验tag：YYYY-MM-DD hh:mm:ss
	if err = validate.RegisterValidation("YYYY-MM-DD hh:mm:ss", validateYMDHMS); err != nil {
		panic(err)
	}
	//添加tag: isEmail
	if err = validate.RegisterValidation("isEmail", checkEmail); err != nil {
		panic(err)
	}
	return validate
}

// VelidatorParms 校验参数
type VelidatorParms struct {
	Value    interface{}
	Rule     string
	ErrorMsg string
}

// Validator 单独值校验
func Validator(vParms *[]VelidatorParms) error {
	for _, parms := range *vParms {
		err := validate.Var(parms.Value, parms.Rule)
		if err != nil {
			return errors.New(parms.ErrorMsg)
		}
	}
	return nil
}

func checkInt(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^[0-9]*$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

func checkEmail(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

func checkYear(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^(19\d{2}|20\d{2}|21\d{2})$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}
func checkMonth(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^(19\d{2}|20\d{2}|21\d{2})-(0[1-9]|1[0-2])$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

func checkMonth1(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^(19\d{2}|20\d{2}|21\d{2})(0[1-9]|1[0-2])$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

func checkStrDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}
	return true
}

func checkDate(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, err := time.Parse("20060102", value)
	if err != nil {
		return false
	}
	return true
}

// validateHMS 校验时分秒
func validateHMS(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

// validateHMS 校验时分秒
func validateShortHMS(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^([01]\d|2[0-3])[0-5]\d[0-5]\d$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

// validateYMDHMS 校验时分秒
func validateYMDHMS(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	reg, err := regexp.Compile(`^((19|20)\d{2}|21\d{2})-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]) ([01]\d|2[0-3]):[0-5]\d:[0-5]\d$`)
	if err != nil {
		return false
	}
	return reg.MatchString(value)
}

func ValidatorStructForJson(s interface{}) (map[string]string, error) {
	res := make(map[string]string)
	sType := reflect.TypeOf(s)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}
	if err := validate.Struct(s); err != nil {
		if rErr, ok := err.(validator.ValidationErrors); ok {
			for _, errItem := range rErr {
				key := errItem.Field()
				sField, _ := sType.FieldByName(key)
				res[sField.Tag.Get("json")] = sField.Tag.Get("validate")
			}
		} else {
			return nil, err
		}
	}
	return res, nil
}

// ValidatorStruct 新的校验器  errorMsg已无用处,可以自定义翻译器translations.go
func ValidatorStruct(s interface{}, errorMsg ...map[string]string) error {
	err := validate.Struct(s)
	sType := reflect.TypeOf(s)
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}
	if err != nil {
		msg := ""
		if rErr, ok := err.(validator.ValidationErrors); ok {
			for _, e := range rErr {
				ss := strings.Split(e.StructNamespace(), ".")
				ss = ss[1:]
				jsonKey := strings.ReplaceAll(Snake(strings.Join(ss, ".")), "._", ".")
				result := e.Translate(tra)
				results := strings.Split(result, " ")
				if len(results) > 0 {
					msg += strings.Replace(result+";", results[0], jsonKey, 1)
				} else {
					msg += result + ";"
				}
			}
		}
		return errors.New(msg)
	}
	return nil
}

// String2Time 将时间字符串转换time.Time
//
// "" -> ZeroTime, "RFC3339" -> Time, "other" -> error
func String2Time(str string) (time.Time, error) {
	t := time.Time{}
	var e error
	if str != "" {
		t, e = time.Parse(time.RFC3339, str)
	}
	return t, e
}
