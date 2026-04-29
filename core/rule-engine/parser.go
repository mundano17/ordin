package parser

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type action struct {
	Move   string `yaml:"move"`
	Delete bool   `yaml:"delete"`
	Copy   string `yaml:"copy"`
}

type data_size struct {
	Enable bool `yaml:"enable"`
	Max    int  `yaml:"max,omitempty"`
	Min    int  `yaml:"min,omitempty"`
	Equal  int  `yaml:"equal,omitempty"`
}

type ruleStruct struct {
	Name      string    `yaml:"name"`
	Extension string    `yaml:"extension"`
	Action    action    `yaml:"action"`
	Data_size data_size `yaml:"data_size"`
	Enable    bool      `yaml:"enable"`
	Priority  int       `yaml:"priority"`
}

func parser() ([]byte, error) {

	data, err := os.ReadFile("rules.yaml")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	fmt.Println(string(data))
	return data, nil

}

func getData(data []byte) (map[string]ruleStruct, error) {
	var t map[string]ruleStruct
	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		return nil, err
	}
	return t, nil

}

func regex_validator(data string) error {
	_, err := regexp.Compile(data)
	return err
}

func action_validator(data *action) error {
	if *data == (action{}) {
		return fmt.Errorf("empty action field")
	}

	data.Copy = strings.TrimSpace(data.Copy)
	data.Move = strings.TrimSpace(data.Move)

	if data.Copy == "" && data.Move == "" && data.Delete == false {
		return fmt.Errorf("invalid action field, useless action field")
	}
	if data.Delete && data.Move != "" {
		return fmt.Errorf("delete and move in same action")
	} else {
		if data.Move != "" && regex_validator(data.Move) != nil {
			return fmt.Errorf("regex error in move")
		}
		if data.Copy != "" && regex_validator(data.Copy) != nil {
			return fmt.Errorf("regex error in copy")
		}
		return nil
	}

}

func data_size_validator(data data_size) error {

	if data == (data_size{}) {
		return fmt.Errorf("empty data_size field")
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
		err := regex_validator(value.Name)
		if err != nil {
			return fmt.Errorf("Error in %s's name \n %s \n", key, err)
		}
		err = regex_validator(value.Extension)
		if err != nil {
			return fmt.Errorf("Error in %s's extension \n %s \n", key, err)
		}
		err = action_validator(&value.Action)
		if err != nil {
			return fmt.Errorf("Error in %s's action \n %s \n", key, err)
		}
		err = data_size_validator(value.Data_size)
		if err != nil {
			return fmt.Errorf("Error in %s's data size \n %s \n", key, err)
		}
	}
	return nil
}
