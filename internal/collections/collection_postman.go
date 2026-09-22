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
	Event   []PostmanEvent  `json:"event,omitempty"`
}

type PostmanEvent struct {
	Listen string        `json:"listen"`
	Script PostmanScript `json:"script"`
}

type PostmanScript struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

// AutoSetBearerTokenEvent monta um script de "Tests" que roda logo após a
// requisição de autenticação: lê o access_token da resposta (padrão OAuth2,
// usado por todos os bancos desta ferramenta) e grava na variável de
// collection bearer_token — assim as próximas requisições já usam o token
// novo automaticamente, sem precisar copiar e colar manualmente no Postman.
func AutoSetBearerTokenEvent() []PostmanEvent {
	return []PostmanEvent{
		{
			Listen: "test",
			Script: PostmanScript{
				Type: "text/javascript",
				Exec: []string{
					"if (pm.response.code === 200 || pm.response.code === 201) {",
					"    var jsonData = pm.response.json();",
					"    if (jsonData.access_token) {",
					"        pm.collectionVariables.set(\"bearer_token\", jsonData.access_token);",
					"        console.log(\"bearer_token atualizado automaticamente.\");",
					"    }",
					"}",
				},
			},
		},
	}
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
	Mode       string              `json:"mode"`
	Urlencoded []PostmanParam      `json:"urlencoded,omitempty"`
	Raw        string              `json:"raw,omitempty"`
	Options    *PostmanBodyOptions `json:"options,omitempty"`
}

type PostmanBodyOptions struct {
	Raw *PostmanRawOptions `json:"raw,omitempty"`
}

type PostmanRawOptions struct {
	Language string `json:"language"`
}

// RawJSONBody monta um body em modo "raw" já marcado com o language "json"
// nas options — sem isso o Postman renderiza o body como texto puro (Text)
// em vez de aplicar o realce de sintaxe JSON no editor "raw".
func RawJSONBody(raw string) *PostmanBody {
	return &PostmanBody{
		Mode: "raw",
		Raw:  raw,
		Options: &PostmanBodyOptions{
			Raw: &PostmanRawOptions{Language: "json"},
		},
	}
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
