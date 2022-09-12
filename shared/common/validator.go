package common

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/go-playground/validator.v9"
)

var Gender, StatusKonsumen, Channel, Lob, Incoming, Home, Education, Marital, ProfID, Photo, Relationship, AppSource, Address, Tenor, Relation string

func NewValidator() *Validator {

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &Validator{
		validator: validate,
	}
}

type Validator struct {
	validator *validator.Validate
	sync      sync.RWMutex
}

func (v *Validator) Validate(i interface{}) error {
	v.sync.Lock()
	v.validator.RegisterValidation("dateformat", dateFormatValidation)
	v.validator.RegisterValidation("url", urlFormatValidation)
	v.validator.RegisterValidation("notnull", notNullValidation)
	v.validator.RegisterValidation("mustnull", mustNullValidation)
	v.sync.Unlock()

	return v.validator.Struct(i)
}
func notNullValidation(fl validator.FieldLevel) bool {
	fmt.Println(fl.Field().Bool())
	return fl.Field().Bool()
}

func mustNullValidation(fl validator.FieldLevel) bool {
	fmt.Println(fl.Field().Bool())
	return fl.Field().Bool()
}

func ftrProspectIDValidation(fl validator.FieldLevel) bool {

	arr := strings.Split(fl.Field().String(), " - ")
	validator, _ := strconv.ParseBool(arr[1])

	return validator
}

func dateFormatValidation(fl validator.FieldLevel) bool {

	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	return re.MatchString(fl.Field().String())
}

func urlFormatValidation(fl validator.FieldLevel) bool {

	re := regexp.MustCompile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)

	return re.MatchString(fl.Field().String())
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
