package br.edu.atitus.currencyservice.dtos;

public record CurrencyQuoteResponse(
        Double paridadeCompra,
        Double paridadeVenda,
        Double cotacaoCompra,
        Double cotacaoVenda,
        String dataHoraCotacao,
        String tipoBoletim
) {
}
