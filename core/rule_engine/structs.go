package rule_engine

type action struct {
	Move   string `yaml:"move"`
	Delete bool   `yaml:"delete"`
	Copy   string `yaml:"copy"`
}

type dataSize struct {
	Enable bool `yaml:"enable"`
	Max    int  `yaml:"max,omitempty"`
	Min    int  `yaml:"min,omitempty"`
	Equal  int  `yaml:"equal,omitempty"`
}

type ruleStruct struct {
	Name      string   `yaml:"name"`
	Extension string   `yaml:"extension"`
	Action    action   `yaml:"action"`
	DataSize  dataSize `yaml:"dataSize"`
	Enable    bool     `yaml:"enable"`
	Priority  int      `yaml:"priority"`
}

type rule struct {
	ruleName string
	ruleData ruleStruct
}

type FilePaths map[string]Paths

type Paths struct {
	Srcpath      string
	MovePaths    []string
	CopyPaths    []string
	DeleteFlag   bool
	NoAccessFlag bool
}
