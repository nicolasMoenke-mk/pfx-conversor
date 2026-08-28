package collections

type PostSicoob struct {
	Method       string `json:"method,omitempty"`
	Endpoint     string `json:"endpoint,omitempty"`
	Token        string `json:"token,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

func AuthPixSicoob() ([]byte, error) {

	return nil, nil
}
