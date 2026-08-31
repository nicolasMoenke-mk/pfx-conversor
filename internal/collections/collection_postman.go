package collections

type PostmanCollection struct {
	Info     PostmanInfo       `json:"info"`
	Variable []PostmanVariable `json:"variable,omitempty"`
	Item     []PostmanItem     `json:"item"`
}

type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

type PostmanVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type PostmanItem struct {
	Name    string          `json:"name"`
	Request *PostmanRequest `json:"request,omitempty"`
	Item    []PostmanItem   `json:"item,omitempty"`
}

type PostmanRequest struct {
	Auth   *PostmanAuth    `json:"auth,omitempty"`
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header"`
	Body   *PostmanBody    `json:"body,omitempty"`
	Url    PostmanURL      `json:"url"`
}

type PostmanAuth struct {
	Type  string             `json:"type"`
	Basic []PostmanAuthParam `json:"basic,omitempty"`
}

type PostmanAuthParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type PostmanBody struct {
	Mode       string         `json:"mode"`
	Urlencoded []PostmanParam `json:"urlencoded,omitempty"`
	Raw        string         `json:"raw,omitempty"`
}

type PostmanParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type PostmanURL struct {
	Raw  string   `json:"raw"`
	Host []string `json:"host,omitempty"`
	Path []string `json:"path,omitempty"`
}
