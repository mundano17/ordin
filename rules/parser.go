// Package rules aims to parse, validate, plan and execute rules
package rules

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

type Action struct {
	Move   string `yaml:"move"`
	Delete bool   `yaml:"delete"`
	Copy   string `yaml:"copy"`
}

type DataSize struct {
	Enable bool `yaml:"enable"`
	Max    int  `yaml:"max,omitempty"`
	Min    int  `yaml:"min,omitempty"`
	Equal  int  `yaml:"equal,omitempty"`
}

type RuleOptions struct {
	Name      string   `yaml:"name"`
	Extension string   `yaml:"extension"`
	Action    Action   `yaml:"action"`
	DataSize  DataSize `yaml:"dataSize"`
	Enable    bool     `yaml:"enable"`
	Priority  int      `yaml:"priority"`
}

type Rule struct {
	Name    string
	Options RuleOptions
}

type Rules []Rule

func LoadRules(data []byte) (Rules, error) {
	t := map[string]RuleOptions{}
	x := Rules{}
	err := yaml.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	for Name, Options := range t {
		Options.normalize()
		newRule := Rule{Name: Name, Options: Options}
		x = append(x, newRule)
	}
	return x, nil
}

func validateRegex(data string) error {
	_, err := regexp.Compile(data)
	return err
}

func validatePath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// to change
			return fmt.Errorf("%s does not exist", path)
		}
		return err
	}
	return nil
}

func (r *RuleOptions) normalize() {
	if r.Name == "" {
		r.Name = "."
	}
	if r.Extension == "" {
		r.Extension = "."
	}
	r.Action.Copy = strings.TrimSpace(r.Action.Copy)
	r.Action.Move = strings.TrimSpace(r.Action.Move)
}

func (data Action) validate() error {
	if data.Copy == "" && data.Move == "" && !data.Delete {
		return fmt.Errorf("useless action field")
	}

	err := validatePath(data.Move)
	if data.Move != "" && err != nil {
		return fmt.Errorf("path error in move, %w", err)
	}
	err = validatePath(data.Copy)
	if data.Copy != "" && err != nil {
		return fmt.Errorf("path error in copy, %w", err)
	}
	return nil
}

func (data DataSize) validate() error {
	if data.Min < 0 {
		return fmt.Errorf("min is negative")
	}
	if data.Equal < 0 {
		return fmt.Errorf("equal is negative")
	}
	if data.Max < 0 {
		return fmt.Errorf("max is negative")
	}
	if data.Min >= data.Max && data.Max != 0 {
		return fmt.Errorf("max is lesser than min")
	}
	return nil
}

func (r Rule) validate() error {
	if err := validateRegex(r.Options.Name); err != nil {
		return fmt.Errorf("%s: %w", r.Name, err)
	}
	if err := validateRegex(r.Options.Extension); err != nil {
		return fmt.Errorf("%s: %w", r.Name, err)
	}
	if err := r.Options.Action.validate(); err != nil {
		return fmt.Errorf("%s: %w", r.Name, err)
	}
	if err := r.Options.DataSize.validate(); err != nil {
		return fmt.Errorf("%s: %w", r.Name, err)
	}
	return nil
}

func (data Rules) validate() error {
	for _, rule := range data {
		if err := rule.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateFile Made for check, validates the rule file and return errors if any.
func ValidateFile(rulesDest string) error {
	data, err := os.ReadFile(rulesDest)
	if err != nil {
		return err
	}
	r, err := LoadRules(data)
	if err != nil {
		return err
	}
	return r.validate()
}

// ValidateAndSort Made for dryrun and run, validates the rule file and returns a sorted rules
func ValidateAndSort(rulesDest string) (Rules, error) {
	data, err := os.ReadFile(rulesDest)
	if err != nil {
		return nil, err
	}
	r, err := LoadRules(data)
	if err != nil {
		return nil, err
	}
	if err = r.validate(); err != nil {
		return nil, err
	}
	r.sort()
	return r, nil
}

func (data Rules) sort() {
	slices.SortFunc(data, func(a, b Rule) int {
		return cmp.Compare(a.Options.Priority, b.Options.Priority)
	})
}
