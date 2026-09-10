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
	DtVencimento          string
	ValidadePosVencimento uint
}

type SicrediChavePix struct {
	Chave string
}

type SicrediDevedor struct {
	Cpf  string
	Nome string
}

type SicrediValor struct {
	Valor string
}

type SicrediCriarPix struct {
	Calendario SicrediCalendario
	ChavePix   SicrediChavePix
	Devedor    SicrediDevedor
	Valor      SicrediValor
}

type SicrediStatusPix struct {
	Status string
}
