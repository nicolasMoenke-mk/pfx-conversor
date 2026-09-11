package body

type BoletoDescontoBB struct {
	Tipo        uint
	DtExpiracao string
	Valor       string
}

type BoletoJurosBB struct {
	Tipo        uint
	Porcentagem string
}

type BoletoMultaBB struct {
	Tipo         uint
	DtVencimento string
	Porcentagem  string
}

type BoletoPagadorBB struct {
	Tipo         uint
	NumInscricao uint
}

type BaixaBoletoBB struct {
	NumeroConv string
}

type DadosBoletoBB struct {
	NumConvenio            string
	NumVariacaoCart        string
	CodigoModalidade       string
	DtEmissao              string
	DtVencimento           string
	VlOriginal             string
	IndicadorTitVencido    string
	NumDiasLimiteReceb     uint
	CdTipoTitulo           uint
	IndPermissaoRecParcial string
	NumTituloBenef         string
	NumTituloClient        string
	Desconto               BoletoDescontoBB
	JurosMora              BoletoJurosBB
	Multa                  BoletoMultaBB
	Pagador                BoletoPagadorBB
	IndicadorPix           string
}

// ----------------- PIX ----------------- //

type BBCalendario struct {
	DtVencimento          string
	ValidadePosVencimento uint
}

type BBChavePix struct {
	Chave string
}

type BBDevedor struct {
	Cpf  string
	Nome string
}

type BBValor struct {
	Valor string
}

type BBCriarPix struct {
	Calendario BBCalendario
	ChavePix   BBChavePix
	Devedor    BBDevedor
	Valor      BBValor
}

type BBStatusPix struct {
	Status string
}
