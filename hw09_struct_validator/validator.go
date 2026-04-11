package hw09structvalidator

import (
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
	panic("implement me")
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("invalid type: %T", v)
	}
	// get all fields of struct via reflection
	fields := val.Type().NumField()
	validators := getValidator()
	for i := 0; i < fields; i++ {
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
		_ = validator(fieldValue.Interface(), validationRules)
	}
	return nil
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
