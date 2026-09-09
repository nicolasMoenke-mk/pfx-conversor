package body

type DadosBoletoSantander struct {
	ClientNumber   string
	NsuCode        string
	ConventCode    string
	BankNumber     string
	DueDate        string
	NsuDate        string
	IssueDate      string
	Key            string
	Enviroment     string
	PartCode       string
	NominalVl      string
	Payer          PagadorSantander
	Beneficiary    string
	DocumentKind   string
	FinePercentage string
	FineQtDays     string
	InterestPerc   string
	DeductionVl    string
	ProtestType    string
	ProtestQtDays  string
	WriteOffQtDays string
	PaymentType    string
	ParcelsQt      string
	VlType         string
	MinVlPerc      string
	MaxVlPerc      string
	IofPerc        string
	Sharing        string
	DigLine        string
	Barcode        string
	QrCodePix      string
	QrCodeUrl      string
	Txid           string
	Message        string
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

type BaixaBoletoSantander struct {
	ConventCode string
	BankNumber  string
	Operation   string
}
