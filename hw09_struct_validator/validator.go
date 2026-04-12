package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotStruct      = errors.New("not a struct")
	ErrInvalidTag     = errors.New("invalid tag")
	ErrUnknownFilter  = errors.New("unknown filter")
	ErrInvalidRegexp  = errors.New("invalid regexp")
	ErrLengthMismatch = errors.New("length mismatch")
	ErrMinValue       = errors.New("less than min")
	ErrMaxValue       = errors.New("greater than max")
	ErrNotInList      = errors.New("not in list")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type ValidatorFunc func(name string, v any, filters []FilterFuncValuePair) ValidationErrors
type typedValidator[T any] func(name string, v T, filters []FilterFuncValuePair) ValidationErrors
type validatorsMap = map[string]ValidatorFunc

type rule struct {
	name  string
	value string
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %s: %v", v.Field, v.Err)
}

func (v ValidationError) Unwrap() error {
	return v.Err
}

func (v ValidationErrors) Error() string {
	errs := make([]error, len(v))
	for i, err := range v {
		errs[i] = err
	}
	return errors.Join(errs...).Error()
}

func (v ValidationErrors) Unwrap() []error {
	errs := make([]error, len(v))
	for i, err := range v {
		errs[i] = err
	}
	return errs
}

func ErrToValidationErrors(err error) ValidationErrors {
	return ValidationErrors{{Err: err}}
}

func Validate(v any) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %T", ErrNotStruct, v)
	}

	validators := getValidator()
	var allErrors ValidationErrors
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		if !field.IsExported() {
			continue
		}

		validationRules := field.Tag.Get("validate")
		if validationRules == "" {
			continue
		}

		fieldValue := val.Field(i)
		typeName := field.Type.String()
		validator, ok := validators[typeName]
		if !ok {
			switch field.Type.Kind() {
			case reflect.Int:
				validator = validators["int"]
			case reflect.String:
				validator = validators["string"]
			case reflect.Slice:
				elemKind := field.Type.Elem().Kind()
				if elemKind == reflect.Int {
					validator = validators["[]int"]
				} else if elemKind == reflect.String {
					validator = validators["[]string"]
				}
			default:
				panic("unhandled default case")
			}
		}

		if validator == nil {
			continue
		}

		rules, err := prepareRules(validationRules)
		if err != nil {
			return err
		}
		filters, err := rulesToFilters(rules)
		if err != nil {
			return err
		}
		fieldErrors := validator(field.Name, fieldValue.Interface(), filters)
		if len(fieldErrors) > 0 {
			allErrors = append(allErrors, fieldErrors...)
		}
	}

	if len(allErrors) == 0 {
		return nil
	}
	return allErrors
}

func getValidator() validatorsMap {
	return validatorsMap{
		"int":      wrapValidator(intValidator),
		"[]int":    wrapValidator(intSliceValidator),
		"string":   wrapValidator(stringValidator),
		"[]string": wrapValidator(stringSliceValidator),
	}
}

func wrapValidator[T any](fn typedValidator[T]) ValidatorFunc {
	return func(name string, v any, filters []FilterFuncValuePair) ValidationErrors {
		typed, ok := v.(T)
		if !ok {
			rv := reflect.ValueOf(v)
			targetType := reflect.TypeFor[T]()
			if rv.Type().ConvertibleTo(targetType) {
				typed = rv.Convert(targetType).Interface().(T)
			} else {
				return ValidationErrors{{Field: name, Err: fmt.Errorf("expected %T, got %T", *new(T), v)}}
			}
		}
		return fn(name, typed, filters)
	}
}

func stringValidator(name, v string, filters []FilterFuncValuePair) ValidationErrors {
	var validationErrors ValidationErrors
	for _, filter := range filters {
		if err := filter.filter(filter.value, v); err != nil {
			validationErrors = append(validationErrors, ValidationError{Field: name, Err: err})
		}
	}
	return validationErrors
}

func stringSliceValidator(name string, v []string, filters []FilterFuncValuePair) ValidationErrors {
	var validationErrors ValidationErrors
	for _, value := range v {
		validationErrors = append(validationErrors, stringValidator(name, value, filters)...)
	}
	return validationErrors
}

func intValidator(name string, v int, filters []FilterFuncValuePair) ValidationErrors {
	var validationErrors ValidationErrors
	for _, filter := range filters {
		if err := filter.filter(filter.value, v); err != nil {
			validationErrors = append(validationErrors, ValidationError{Field: name, Err: err})
		}
	}
	return validationErrors
}

func intSliceValidator(name string, v []int, filters []FilterFuncValuePair) ValidationErrors {
	var validationErrors ValidationErrors
	for _, value := range v {
		validationErrors = append(validationErrors, intValidator(name, value, filters)...)
	}
	return validationErrors
}
