package body

type ChavePix struct {
	Type    string `json:"type"`
	DictKey string `json:"dictKey"`
}

type PagadorSantander struct {
	Name         string `json:"name"`
	DocumentType string `json:"documentType"`
	Document     string `json:"documentNumber"`
	Address      string `json:"address"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
}

type DadosBoletoSantander struct {
	ClientNumber         string           `json:"clientNumber"`
	NsuCode              string           `json:"nsuCode"`
	CovenantCode         string           `json:"covenantCode"`
	BankNumber           string           `json:"bankNumber"`
	DueDate              string           `json:"dueDate"`
	NsuDate              string           `json:"nsuDate"`
	IssueDate            string           `json:"issueDate"`
	Key                  ChavePix         `json:"key"`
	Environment          string           `json:"environment"`
	ParticipantCode      *string          `json:"participantCode"`
	NominalValue         string           `json:"nominalValue"`
	Payer                PagadorSantander `json:"payer"`
	Beneficiary          *string          `json:"beneficiary"`
	DocumentKind         string           `json:"documentKind"`
	FinePercentage       string           `json:"finePercentage"`
	FineQuantityDays     string           `json:"fineQuantityDays"`
	InterestPercentage   string           `json:"interestPercentage"`
	DeductionValue       *string          `json:"deductionValue"`
	ProtestType          string           `json:"protestType"`
	ProtestQuantityDays  *string          `json:"protestQuantityDays"`
	WriteOffQuantityDays *string          `json:"writeOffQuantityDays"`
	PaymentType          string           `json:"paymentType"`
	ParcelsQuantity      *string          `json:"parcelsQuantity"`
	ValueType            *string          `json:"valueType"`
	MinValueOrPercentage *string          `json:"minValueOrPercentage"`
	MaxValueOrPercentage *string          `json:"maxValueOrPercentage"`
	IofPercentage        *string          `json:"iofPercentage"`
	Sharing              *string          `json:"sharing"`
	DigitableLine        string           `json:"digitableLine"`
	Barcode              string           `json:"barcode"`
	QrCodePix            string           `json:"qrCodePix"`
	QrCodeUrl            string           `json:"qrCodeUrl"`
	Txid                 *string          `json:"txid"`
	Messages             *string          `json:"messages"`
}

type BaixaBoletoSantander struct {
	ConvenantCode string `json:"convenantCode"`
	BankNumber    string `json:"bankNumber"`
	Operation     string `json:"operation"`
}

// ------------------ PIX ------------------ //

type SantanderCalendario struct {
	DtVencimento          string `json:"dataDeVencimento"`
	ValidadePosVencimento uint   `json:"validadeAposVencimento"`
}

type SantanderChavePix struct {
	Chave string `json:"chave_pix"`
}

type SantanderPixDevedor struct {
	Cpf  string `json:"cpf"`
	Nome string `json:"nome"`
}

type SantanderPixValor struct {
	Valor string `json:"original"`
}

type SantanderStatusPix struct {
	Status string `json:"status"`
}

type SantanderCriarPix struct {
	Calendario SantanderCalendario `json:"calendario"`
	Chave      SantanderChavePix   `json:"chave"`
	Devedor    SantanderPixDevedor `json:"devedor"`
	Valor      SantanderPixValor   `json:"valor"`
}
