package sdk

type Properties struct {
	Category  string `json:"category"`
	Fodmap    string `json:"fodmap"`
	Condition string `json:"condition,omitempty"`
	Note      string `json:"note,omitempty"`
}
