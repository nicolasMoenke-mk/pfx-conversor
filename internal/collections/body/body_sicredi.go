package body

type SicrediPagador struct {
	Documento  string
	Nome       string
	TipoPessoa string
	Endereco   string
	Cidade     string
	Cep        string
	Uf         string
}

type DadosBoletoSicredi struct {
	CodBeneficiario string
	DtVencimento    string
	SeuNumero       string
	EspecieDoc      string
	Pagador         SicrediPagador
	TipoConbranca   string
	Valor           string
	TipoJuros       string
	Juros           string
	Multa           string
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
