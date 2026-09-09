package body

//Body de requisição para a Van Bancária

type SicoobDadosBoleto struct {
	NumeroCliente          int     `json:"numeroCliente"`
	CodigoModalidade       int     `json:"codigoModalidade"`
	ContaCorrente          int     `json:"numeroContaCorrente"`
	DataVencimento         string  `json:"dataVencimento"`
	DataEmissao            string  `json:"dataEmissao"`
	SeuNumero              int     `json:"seuNumero"`
	NossoNumero            int     `json:"nossoNumero"`
	IdBoletoEmpresa        int     `json:"identificacaoBoletoEmpresa"`
	CodigoEspecieDocumento string  `json:"codigoEspecieDocumento"`
	IdEmissaoBoleto        int     `json:"identificacaoEmissaoBoleto"`
	IdDistribuicaoBoleto   int     `json:"identificacaoDistribuicaoBoleto"`
	Valor                  float64 `json:"valor"`
	TipoDesconto           int     `json:"tipoDesconto"`
	TipoMulta              int     `json:"tipoMulta"`
	TipoJurosMora          int     `json:"tipoJurosMora"`
	NumeroParcela          int     `json:"numeroParcela"`
	Aceite                 bool    `json:"aceite"`
	CdNegativacao          int     `json:"codigoNegativacao"`
	CdProtesto             int     `json:"codigoProtesto"`
}

type SicoobPagadorBoleto struct {
	Documento string `json:"numeroCpfCnpj"`
	Nome      string `json:"nome"`
	Endereco  string `json:"endereco"`
	Bairro    string `json:"bairro"`
	Cidade    string `json:"cidade"`
	Cep       string `json:"cep"`
	Uf        string `json:"uf"`
	Email     string `json:"email,omitempty"`
}

type SicoobPixBoleto struct {
	CdCadastroPix int `json:"codigoCadastrarPIX"`
}

type SicoobCriarBoleto struct {
	Dados   SicoobDadosBoleto   `json:"dados"`
	Pagador SicoobPagadorBoleto `json:"pagador"`
	Pix     SicoobPixBoleto     `json:"pix"`
}

type SicoobRemoverBoleto struct {
	Boleto SicoobDadosBoleto `json:"boleto"`
}

// ------------------------------------------- //

//Body de requisição para o PIX

type SicoobCalendario struct {
	DtVencimento string `json:"dataDeVencimento,omitempty"`
	Validade     uint8  `json:"validadeAposVencimento,omitempty"`
}

type SicoobChavePix struct {
	Chave string `json:"chave,omitempty"`
}

type SicoobDevedor struct {
	Cpf  string `json:"cpf,omitempty"`
	Nome string `json:"nome,omitempty"`
}

type SicoobValor struct {
	Valor string `json:"valor"`
}

type SicoobCriarPix struct {
	Calendario SicoobCalendario `json:"calendario"`
	Chave      SicoobChavePix   `json:"chave"`
	Devedor    SicoobDevedor    `json:"devedor"`
	Valor      SicoobValor      `json:"valor"`
}

type SicoobStatusPix struct {
	Status string `json:"status"`
}
