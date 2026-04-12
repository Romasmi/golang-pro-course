package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type validator[T any] = func(v T, rules string) ValidationErrors
type validatorsMap = map[string]validator[any]

func (v ValidationErrors) Error() string {
	return errors.Join(v).Error()
}

func ErrToValidationErrors(err error) ValidationErrors {
	return ValidationErrors{ValidationError{Err: err}}
}

func Validate(v interface{}) ValidationErrors {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrToValidationErrors(fmt.Errorf("invalid type: %T", v))
	}

	validators := getValidator()
	allErrors := make(ValidationErrors, 0)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		typeName := field.Type.Name()
		validator, ok := validators[typeName]
		if !ok {
			continue
		}
		validationRules := field.Tag.Get("validate")
		if validationRules == "" {
			continue
		}
		fieldValue := val.Field(i)
		allErrors = append(allErrors, validator(fieldValue.Interface(), validationRules)...)
	}
	return allErrors
}

func getValidator() validatorsMap {
	return validatorsMap{
		"int":      intValidator,
		"[]int":    intSliceValidator,
		"string":   stringValidator,
		"[]string": stringSliceValidator,
	}
}

func stringValidator(v any, rules string) ValidationErrors {
	return nil
}

func stringSliceValidator(v any, rules string) ValidationErrors {
	return nil
}

func intValidator(v any, rules string) ValidationErrors {
	return nil
}

func intSliceValidator(v any, rules string) ValidationErrors {
	return nil
}
