package validator

import (
	"fmt"
	"testing"
	"time"
)

type student struct {
	Age string `json:"age" validate:"INTEGER"`
}

func TestValidateInt(t *testing.T) {
	value := "22001212"
	_, err := time.Parse("20060102", value)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(true)
}

type MyStruct struct {
	String []string `validate:"is-awesome=123456"`
}

func TestNewValidator(t *testing.T) {
	val := validator.New()
	val.RegisterValidation("is-awesome", ValidateMyVal)
	s := MyStruct{String: []string{"awesome"}}
	err := val.Struct(s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("校验通过")
}

// map 收集校验
func ValidateMyVal(fl validator.FieldLevel) bool {
	fmt.Println(fl)
	fmt.Println("afdafdasf", fl.Param())
	fmt.Println(fl.Field().String())
	return true
}
