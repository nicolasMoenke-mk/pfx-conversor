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
