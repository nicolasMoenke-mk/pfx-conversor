package body

type SicrediPagador struct {
	Documento  string `json:"documento"`
	Nome       string `json:"nome"`
	TipoPessoa string `json:"tipoPessoa"`
	Endereco   string `json:"endereco"`
	Cidade     string `json:"cidade"`
	Cep        string `json:"cep"`
	Uf         string `json:"uf"`
}

type DadosBoletoSicredi struct {
	CodBeneficiario string         `json:"codigoBeneficiario"`
	DtVencimento    string         `json:"dataVencimento"`
	SeuNumero       string         `json:"seuNumero"`
	EspecieDoc      string         `json:"especieDocumento"`
	Pagador         SicrediPagador `json:"pagador"`
	TipoConbranca   string         `json:"tipoCobranca"`
	Valor           string         `Json:"valor"`
	TipoJuros       string         `Json:"tipoJuros"`
	Juros           string         `Json:"juros"`
	Multa           string         `Json:"multa"`
}

// ----------------- PIX ----------------- //

type SicrediCalendario struct {
	DtVencimento          string `json:"dataDevencimento"`
	ValidadePosVencimento uint   `json:"validadeAposVencimento"`
}

type SicrediChavePix struct {
	Chave string `json:"chave_pix"`
}

type SicrediDevedor struct {
	Cpf  string `json:"cpf"`
	Nome string `json:"nome"`
}

type SicrediValor struct {
	Valor string `json:"original"`
}

type SicrediCriarPix struct {
	Calendario SicrediCalendario `json:"calendario"`
	ChavePix   SicrediChavePix   `json:"chave"`
	Devedor    SicrediDevedor    `json:"devedor"`
	Valor      SicrediValor      `json:"valor"`
}

type SicrediStatusPix struct {
	Status string `json:"status"`
}
