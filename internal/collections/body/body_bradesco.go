package body

type DadosBoletoBradesco struct {
	NuCPFCNPJ                            int     `json:"nuCPFCNPJ"`
	FilialCPFCNPJ                        int     `json:"filialCPFCNPJ"`
	CtrlCPFCNPJ                          int     `json:"ctrlCPFCNPJ"`
	IdProduto                            int     `json:"idProduto"`
	NuNegociacao                         int64   `json:"nuNegociacao"`
	NuTitulo                             int     `json:"nuTitulo"`
	NuCliente                            string  `json:"nuCliente"`
	DtEmissaoTitulo                      string  `json:"dtEmissaoTitulo"`
	DtVencimentoTitulo                   string  `json:"dtVencimentoTitulo"`
	TpVencimento                         int     `json:"tpVencimento"`
	VlNominalTitulo                      float64 `json:"vlNominalTitulo"`
	CdEspecieTitulo                      int     `json:"cdEspecieTitulo"`
	TpProtestoAutomaticoNegativacao      int     `json:"tpProtestoAutomaticoNegativacao"`
	PrazoProtestoAutomaticoNegativacao   int     `json:"prazoProtestoAutomaticoNegativacao"`
	TipoPrazoDecursoTres                 int     `json:"tipoPrazoDecursoTres"`
	ControleParticipante                 string  `json:"controleParticipante,omitempty"`
	CdPagamentoParcial                   string  `json:"cdPagamentoParcial"`
	QtdePagamentoParcial                 int     `json:"qtdePagamentoParcial"`
	PercentualJuros                      int     `json:"percentualJuros"`
	VlJuros                              float64 `json:"vlJuros"`
	QtdeDiasJuros                        int     `json:"qtdeDiasJuros"`
	PercentualMulta                      int     `json:"percentualMulta"`
	VlMulta                              float64 `json:"vlMulta"`
	QtdeDiasMulta                        int     `json:"qtdeDiasMulta"`
	PercentualDesconto1                  int     `json:"percentualDesconto1"`
	VlDesconto1                          float64 `json:"vlDesconto1"`
	DataLimiteDesconto1                  string  `json:"dataLimiteDesconto1,omitempty"`
	PercentualDesconto2                  int     `json:"percentualDesconto2"`
	VlDesconto2                          float64 `json:"vlDesconto2"`
	DataLimiteDesconto2                  string  `json:"dataLimiteDesconto2,omitempty"`
	PercentualDesconto3                  int     `json:"percentualDesconto3"`
	VlDesconto3                          float64 `json:"vlDesconto3"`
	DataLimiteDesconto3                  string  `json:"dataLimiteDesconto3,omitempty"`
	PrazoBonificacao                     int     `json:"prazoBonificacao"`
	PercentualBonificacao                int     `json:"percentualBonificacao"`
	VlBonificacao                        float64 `json:"vlBonificacao"`
	DtLimiteBonificacao                  string  `json:"dtLimiteBonificacao,omitempty"`
	VlAbatimento                         float64 `json:"vlAbatimento"`
	VlIOF                                float64 `json:"vlIOF"`
	NomePagador                          string  `json:"nomePagador"`
	LogradouroPagador                    string  `json:"logradouroPagador"`
	NuLogradouroPagador                  string  `json:"nuLogradouroPagador"`
	ComplementoLogradouroPagador         string  `json:"complementoLogradouroPagador,omitempty"`
	CepPagador                           int     `json:"cepPagador"`
	ComplementoCepPagador                int     `json:"complementoCepPagador"`
	BairroPagador                        string  `json:"bairroPagador"`
	MunicipioPagador                     string  `json:"municipioPagador"`
	UfPagador                            string  `json:"ufPagador"`
	CdIndCpfcnpjPagador                  int     `json:"cdIndCpfcnpjPagador"`
	NuCpfcnpjPagador                     int64   `json:"nuCpfcnpjPagador"`
	EndEletronicoPagador                 string  `json:"endEletronicoPagador,omitempty"`
	NomeSacadorAvalista                  string  `json:"nomeSacadorAvalista,omitempty"`
	LogradouroSacadorAvalista            string  `json:"logradouroSacadorAvalista,omitempty"`
	NuLogradouroSacadorAvalista          string  `json:"nuLogradouroSacadorAvalista,omitempty"`
	ComplementoLogradouroSacadorAvalista string  `json:"complementoLogradouroSacadorAvalista,omitempty"`
	CepSacadorAvalista                   int     `json:"cepSacadorAvalista,omitempty"`
	ComplementoCepSacadorAvalista        int     `json:"complementoCepSacadorAvalista,omitempty"`
	BairroSacadorAvalista                string  `json:"bairroSacadorAvalista,omitempty"`
	MunicipioSacadorAvalista             string  `json:"municipioSacadorAvalista,omitempty"`
	UfSacadorAvalista                    string  `json:"ufSacadorAvalista,omitempty"`
	CdIndCpfcnpjSacadorAvalista          int     `json:"cdIndCpfcnpjSacadorAvalista,omitempty"`
	NuCpfcnpjSacadorAvalista             int64   `json:"nuCpfcnpjSacadorAvalista,omitempty"`
	EnderecoSacadorAvalista              string  `json:"enderecoSacadorAvalista,omitempty"`
}

type BoletoBradesco struct {
	BoletoBradesco DadosBoletoBradesco `json:"boleto_bradesco"`
}

type DocPagadorBradesco struct {
	CpfCnpj  int `json:"cpfCnpj"`
	Filial   int `json:"filial"`
	Controle int `json:"controle"`
}

type BaixaBoletoBradesco struct {
	CpfCnpj     DocPagadorBradesco `json:"cpfCnpj"`
	Produto     int                `json:"produto"`
	Negociacao  int64              `json:"negociacao"`
	NossoNumero int                `json:"nossoNumero"`
	Sequencia   int                `json:"sequencia"`
	CodigoBaixa int                `json:"codigoBaixa"`
}

// ------------------- PIX Bradesco ------------------- //

type BradescoCalendario struct {
	DtVencimento          string `json:"dataDeVencimento"`
	ValidadePosVencimento uint   `json:"ValidadeAposVencimento"`
}

type BradescoChavePix struct {
	ChavePix string `json:"chave"`
}

type BradescoPixDevedor struct {
	Cpf  string `json:"cpf"`
	Nome string `json:"nome"`
}

type BradescoPixValor struct {
	Valor string `json:"original"`
}

type BradescoStatusPix struct {
	Status string `json:"status"`
}

type BradescoCriarPix struct {
	Calendario BradescoCalendario `json:"calendario"`
	Chave      BradescoChavePix   `json:"chave"`
	Devedor    BradescoPixDevedor `json:"devedor"`
	Valor      BradescoPixValor   `json:"valor"`
}
