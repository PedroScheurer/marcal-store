package dtos

type BCBResponse struct {
	Value []CurrencyQuoteResponse `json:"value"`
}

type CurrencyQuoteResponse struct {
	ParidadeCompra  float64 `json:"paridadeCompra"`
	ParidadeVenda   float64 `json:"paridadeVenda"`
	CotacaoCompra   float64 `json:"cotacaoCompra"`
	CotacaoVenda    float64 `json:"cotacaoVenda"`
	DataHoraCotacao string  `json:"dataHoraCotacao"`
	TipoBoletim     string  `json:"tipoBoletim"`
}
