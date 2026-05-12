// Package ruleengine aims to parse, validate, plan and execute rules
package ruleengine

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

func getData(data []byte) (map[string]ruleStruct, error) {
	var t map[string]ruleStruct
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func regexValidator(data string) error {
	_, err := regexp.Compile(data)
	return err
}

func PathValidator(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist", path)
		} else {
			return err
		}
	}
	return nil
}

// allow !delete and copy/move/copy+move or only delete
func actionValidator(data *action) error {
	if *data == (action{}) {
		return fmt.Errorf("empty action field")
	}

	data.Copy = strings.TrimSpace(data.Copy)
	data.Move = strings.TrimSpace(data.Move)

	if data.Copy == "" && data.Move == "" && !data.Delete {
		return fmt.Errorf("invalid action field, useless action field")
	}

	if data.Delete && data.Move != "" || data.Delete && data.Copy != "" {
		return fmt.Errorf("delete and any other action is not allowed")
	} else {
		err := PathValidator(data.Move)
		if data.Move != "" && err != nil {
			return fmt.Errorf("path error in move\n %s", err)
		}
		err = PathValidator(data.Copy)
		if data.Copy != "" && err != nil {
			return fmt.Errorf("path error in copy\n %s", err)
		}
		return nil
	}
}

func dataSizeValidator(data dataSize) error {
	if data == (dataSize{}) {
		return fmt.Errorf("empty dataSize field")
	}

	if data.Min >= data.Max && data.Max != 0 {
		return fmt.Errorf("max is lesser than min")
	}
	if data.Min < 0 {
		return fmt.Errorf("min is negative")
	}
	if data.Max < 0 {
		return fmt.Errorf("max is negative")
	}
	return nil
}

func validator(data map[string]ruleStruct) error {
	for key, value := range data {
		err := regexValidator(value.Name)
		if err != nil {
			return fmt.Errorf("error in %s's name \n %s", key, err)
		}
		err = regexValidator(value.Extension)
		if err != nil {
			return fmt.Errorf("error in %s's extension \n %s", key, err)
		}
		err = actionValidator(&value.Action)
		if err != nil {
			return fmt.Errorf("error in %s's action \n %s", key, err)
		}
		err = dataSizeValidator(value.DataSize)
		if err != nil {
			return fmt.Errorf("error in %s's data size \n %s", key, err)
		}
	}
	return nil
}

func CheckPipeline(path string) (map[string]ruleStruct, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	yamlData, err := getData(data)
	if err != nil {
		return nil, err
	}
	err = validator(yamlData)
	if err != nil {
		return nil, err
	}
	return yamlData, nil
}

func CheckSort(path string) ([]rule, error) {
	yamlData, err := CheckPipeline(path)
	if err != nil {
		return nil, err
	}

	sortedRules := []rule{}

	for key, value := range yamlData {
		if !value.Enable {
			continue
		}
		if value.Name == "" {
			value.Name = "."
		}
		if value.Extension == "" {
			value.Extension = "."
		}

		sortedRules = append(sortedRules, rule{key, value})
	}

	slices.SortFunc(sortedRules, func(a rule, b rule) int {
		return cmp.Compare(a.ruleData.Priority, b.ruleData.Priority)
	})

	return sortedRules, nil
}
