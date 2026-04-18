package hw09structvalidator

import (
	"fmt"
	"strings"
)

func prepareRules(rules string) ([]rule, error) {
	rawRules := strings.Split(rules, "|")
	res := make([]rule, 0, len(rawRules))
	for _, rawRule := range rawRules {
		r, err := rawRuleToRule(rawRule)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, nil
}

func rawRuleToRule(rawRule string) (rule, error) {
	name, value, found := strings.Cut(rawRule, ":")
	if !found {
		return rule{}, fmt.Errorf("invalid rule: %s", rawRule)
	}
	return rule{
		name:  name,
		value: value,
	}, nil
}
