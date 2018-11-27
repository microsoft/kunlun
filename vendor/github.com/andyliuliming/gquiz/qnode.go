package gquiz

// QNode sample:
// name: root
// questions:
// 	- var_name: kind
// 	  description: "Which kind of artifacts do you want?"
// 	  candidates:
// 		- value: infrastructure
// 		description: "To deployment the infrastructures, like: kubernetes..."
// 		- value: application
// 		description: "To deploy your application with code, like: php, java..."
// transitions:
// 	 - name: infrastructure
// 		 condition: "kind == 'infrastructure'"
// 	 - name: application
// 		 condition: "kind == 'application'"
type QNode struct {
	Name        string       `yaml:"name"`
	Questions   []Question   `yaml:"questions,omitempty"`
	Transitions []Transition `yaml:"transitions"`
}

type Question struct {
	VarName     string            `yaml:"var_name"`
	Type        string            `yaml:"type"`
	Description string            `yaml:"description"`
	Candidates  []CandidateAnswer `yaml:"candidates,omitempty"`
	Default     string            `yaml:"default"`
	DefaultEnv  string            `yaml:"default_env"`
	Persistent  bool              `yaml:"persistent"`
}

type CandidateAnswer struct {
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}

type Transition struct {
	Name      string `yaml:"name"`
	Condition string `yaml:"condition"`
}
