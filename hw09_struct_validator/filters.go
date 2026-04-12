package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type FilterFunc func(string, interface{}) error

type FilterFuncValuePair struct {
	value  string
	filter FilterFunc
}

func getFilter(filterName string) (FilterFunc, error) {
	switch filterName {
	case "len":
		return lenFilter, nil
	case "regexp":
		return regExpFilter, nil
	case "in":
		return inFilter, nil
	case "min":
		return minFilter, nil
	case "max":
		return maxFilter, nil
	}
	return nil, fmt.Errorf("unknown filter: %s", filterName)
}

func rulesToFilters(rules []rule) ([]FilterFuncValuePair, error) {
	filters := make([]FilterFuncValuePair, len(rules))
	for i, rule := range rules {
		filter, err := getFilter(rule.name)
		if err != nil {
			return nil, err
		}
		filters[i] = FilterFuncValuePair{
			value:  rule.value,
			filter: filter,
		}
	}
	return filters, nil
}

func lenFilter(ruleValue string, value interface{}) error {
	length, err := strconv.Atoi(ruleValue)
	if err != nil || length < 0 {
		return fmt.Errorf("%w: invalid length %s", ErrInvalidTag, ruleValue)
	}

	var actualLength int
	switch v := value.(type) {
	case string:
		actualLength = utf8.RuneCountInString(v)
	case []string:
		actualLength = len(v)
	case []int:
		actualLength = len(v)
	default:
		return fmt.Errorf("invalid type for lenFilter: %T", value)
	}

	if actualLength != length {
		return fmt.Errorf("%w: expected %d, got %d", ErrLengthMismatch, length, actualLength)
	}
	return nil
}

func regExpFilter(ruleValue string, value interface{}) error {
	switch v := value.(type) {
	case string:
		pattern := ruleValue
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrInvalidRegexp, pattern, err)
		}
		if !re.MatchString(v) {
			return fmt.Errorf("%w: %s", ErrInvalidRegexp, pattern)
		}
	default:
		return fmt.Errorf("invalid type for regExpFilter: %T", value)
	}
	return nil
}

func inFilter(ruleValue string, value interface{}) error {
	allowed := strings.Split(ruleValue, ",")

	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case int:
		strValue = strconv.Itoa(v)
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.String {
			strValue = rv.String()
		} else if rv.Kind() == reflect.Int {
			strValue = strconv.FormatInt(rv.Int(), 10)
		} else {
			return fmt.Errorf("invalid type for inFilter: %T", value)
		}
	}

	if !slices.Contains(allowed, strValue) {
		return fmt.Errorf("%w: %s not in %v", ErrNotInList, strValue, allowed)
	}
	return nil
}

func minFilter(ruleValue string, value interface{}) error {
	minV, err := strconv.Atoi(ruleValue)
	if err != nil {
		return fmt.Errorf("%w: invalid min %s", ErrInvalidTag, ruleValue)
	}

	var intValue int
	switch v := value.(type) {
	case int:
		intValue = v
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Int {
			intValue = int(rv.Int())
		} else {
			return fmt.Errorf("invalid type for minFilter: %T", value)
		}
	}

	if intValue < minV {
		return fmt.Errorf("%w: expected min %d, got %d", ErrMinValue, minV, intValue)
	}
	return nil
}

func maxFilter(ruleValue string, value interface{}) error {
	maxV, err := strconv.Atoi(ruleValue)
	if err != nil {
		return fmt.Errorf("%w: invalid max %s", ErrInvalidTag, ruleValue)
	}

	var intValue int
	switch v := value.(type) {
	case int:
		intValue = v
	default:
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Int {
			intValue = int(rv.Int())
		} else {
			return fmt.Errorf("invalid type for maxFilter: %T", value)
		}
	}

	if intValue > maxV {
		return fmt.Errorf("%w: expected max %d, got %d", ErrMaxValue, maxV, intValue)
	}
	return nil
}
