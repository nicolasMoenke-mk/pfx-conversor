package body

type BoletoDescontoBB struct {
	Tipo        uint   `json:"tipo"`
	DtExpiracao string `json:"dataExpiracao"`
	Valor       string `json:"valor"`
}

type BoletoJurosBB struct {
	Tipo        uint   `json:"tipo"`
	Porcentagem string `json:"porcentagem"`
}

type BoletoMultaBB struct {
	Tipo         uint   `json:"tipo"`
	DtVencimento string `json:"data"`
	Porcentagem  string `json:"porcentagem"`
}

type BoletoPagadorBB struct {
	Tipo         uint `json:"tipoInscricao"`
	NumInscricao uint `json:"numeroInscricao"`
}

type BaixaBoletoBB struct {
	NumeroConv string `json:"numeroConvenio"`
}

type DadosBoletoBB struct {
	NumConvenio            string           `json:"numeroConvenio"`
	NumCarteira            string           `json:"numeroCarteira"`
	NumVariacaoCart        string           `json:"numeroVariacaoCarteira"`
	CodigoModalidade       string           `json:"codigoModalidade"`
	DtEmissao              string           `json:"dataEmissao"`
	DtVencimento           string           `json:"dataVencimento"`
	VlOriginal             string           `json:"valorOriginal"`
	IndicadorTitVencido    string           `json:"indicadorAceiteTituloVencido"`
	NumDiasLimiteReceb     uint             `json:"numeroDiasLimiteRecebimento"`
	CdTipoTitulo           uint             `json:"codigoTipoTitulo"`
	IndPermissaoRecParcial string           `json:"indicadorPermissaoRecebimentoParcial"`
	NumTituloBenef         string           `json:"numeroTituloBeneficiario"`
	NumTituloClient        string           `json:"numeroTituloCliente"`
	Desconto               BoletoDescontoBB `json:"desconto"`
	JurosMora              BoletoJurosBB    `json:"jurosMora"`
	Multa                  BoletoMultaBB    `json:"multa"`
	Pagador                BoletoPagadorBB  `json:"pagador"`
	IndicadorPix           string           `json:"indicadorPix"`
}

// ----------------- PIX ----------------- //

type BBCalendario struct {
	DtVencimento          string `json:"dataDeVencimento"`
	ValidadePosVencimento uint   `json:"validadeAposVencimento"`
}

type BBChavePix struct {
	Chave string `json:"chave_pix"`
}

type BBDevedor struct {
	Cpf  string `json:"cpf"`
	Nome string `json:"nome"`
}

type BBValor struct {
	Valor string `json:"original"`
}

type BBCriarPix struct {
	Calendario BBCalendario `json:"calendario"`
	ChavePix   BBChavePix   `json:"chave"`
	Devedor    BBDevedor    `json:"devedor"`
	Valor      BBValor      `json:"valor"`
}

type BBStatusPix struct {
	Status string `json:"status"`
}
