package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrValidationLength   = errors.New("invalid: length")
	ErrValidationMinimum  = errors.New("invalid: minimum")
	ErrValidationMaximum  = errors.New("invalid: maximum")
	ErrValidationContains = errors.New("invalid: not contains")
	ErrNotStruct          = errors.New("not struct is passed to validation")
	ErrNoRule             = errors.New("no such validation rule")
	ErrValidationRegexp   = errors.New("invalid: regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	buff := strings.Builder{}

	for _, e := range v {
		buff.WriteString(fmt.Sprintf("Field: %s, Error: %v\n", e.Field, e.Err))
	}

	return buff.String()
}

func Validate(v interface{}) error {
	vErrors := ValidationErrors{}
	e := reflect.ValueOf(v)

	if e.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < e.NumField(); i++ {
		field := e.Type().Field(i)
		varName := field.Name
		varTag := field.Tag
		if varTag == "" {
			continue
		}

		varValue := e.Field(i).Interface()

		if e.Field(i).Kind() == reflect.Slice {
			for _, sliceVal := range varValue.([]string) {
				if err := validateValue(varTag, sliceVal); err != nil {
					vErrors = append(vErrors, ValidationError{Field: varName, Err: err})
				}
			}
		} else if err := validateValue(varTag, varValue); err != nil {
			vErrors = append(vErrors, ValidationError{Field: varName, Err: err})
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}

	return nil
}

func validateValue(varTag reflect.StructTag, varValue interface{}) error {
	if varTagVal := varTag.Get("validate"); varTagVal != "" {
		valRules := strings.Split(varTagVal, "|")

		for _, rawRule := range valRules {
			if err := validateRule(rawRule, varValue); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRule(rawRule string, varValue interface{}) error {
	valRule := strings.Split(rawRule, ":")

	if len(valRule) != 2 {
		return ErrNoRule
	}

	rule := valRule[0]
	val := valRule[1]

	switch rule {
	case "len":
		if intVal, err := strconv.Atoi(val); err != nil || intVal != len(varValue.(string)) {
			return ErrValidationLength
		}
	case "min":
		if intVal, err := strconv.Atoi(val); err != nil || intVal > varValue.(int) {
			return ErrValidationMinimum
		}
	case "max":
		if intVal, err := strconv.Atoi(val); err != nil || intVal < varValue.(int) {
			return ErrValidationMaximum
		}
	case "in":
		inString := strings.Split(val, ",")

		if !contains(inString, fmt.Sprintf("%v", varValue)) {
			return ErrValidationContains
		}
	case "regexp":
		matched, err := regexp.MatchString(val, fmt.Sprintf("%v", varValue))
		if err != nil {
			return err
		}

		if !matched {
			return ErrValidationRegexp
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
