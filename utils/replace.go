package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/xuri/excelize/v2"
)

// invoiceColumnAliases são os nomes de coluna aceitos, já normalizados
// (minúsculo, sem espaço/underscore/hífen/acento) via normalizeHeader.
var invoiceColumnAliases = []string{
	"fatura",
	"faturas",
	"codfatura",
	"codigofatura",
	"numfatura",
	"numerofatura",
	"nrfatura",
}

var accentReplacer = strings.NewReplacer(
	"á", "a", "à", "a", "ã", "a", "â", "a", "ä", "a",
	"é", "e", "ê", "e", "è", "e", "ë", "e",
	"í", "i", "ì", "i", "î", "i", "ï", "i",
	"ó", "o", "ô", "o", "õ", "o", "ò", "o", "ö", "o",
	"ú", "u", "û", "u", "ù", "u", "ü", "u",
	"ç", "c",
)

func normalizeHeader(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = accentReplacer.Replace(s)
	s = strings.NewReplacer(" ", "", "_", "", "-", "", ".", "").Replace(s)
	return s
}

// maxHeaderCellLength evita que uma frase solta (título do relatório,
// legenda etc.) seja confundida com um cabeçalho de coluna só porque
// contém a palavra "fatura" em algum ponto do texto.
const maxHeaderCellLength = 30

func findInvoiceColumn(headers []string) (int, bool) {
	for i, h := range headers {
		norm := normalizeHeader(h)
		if len(norm) > maxHeaderCellLength {
			continue
		}
		for _, alias := range invoiceColumnAliases {
			if strings.Contains(norm, alias) {
				return i, true
			}
		}
	}
	return -1, false
}

// findHeaderRow procura, dentre as primeiras linhas da planilha, a primeira
// que contenha uma coluna de fatura — cobre planilhas onde a linha 0 não é
// o cabeçalho de verdade (título do relatório, linhas em branco etc.). Uma
// planilha com uma única coluna (só "Fatura", sem mais nada) é um caso
// válido e não é descartada aqui — quem filtra título de relatório solto
// (frase longa que contenha a palavra "fatura") é o maxHeaderCellLength em
// findInvoiceColumn, não a quantidade de colunas da linha.
func findHeaderRow(rows [][]string) (headerRow int, colIndex int, err error) {
	maxScan := 10
	if maxScan > len(rows) {
		maxScan = len(rows)
	}

	for i := 0; i < maxScan; i++ {
		if idx, ok := findInvoiceColumn(rows[i]); ok {
			return i, idx, nil
		}
	}

	return -1, -1, fmt.Errorf("nenhuma coluna de fatura encontrada nas primeiras %d linhas da planilha — cabeçalho da linha 1: %v", maxScan, rows[0])
}

// numeroComSeparadorDeMilhar reconhece um número inteiro formatado com
// separador de milhar — vírgula ou espaço (padrão en-US/agrupamento SI) ou
// ponto (padrão pt-BR, o mais comum vindo do Excel em português), sempre
// agrupado de 3 em 3 dígitos. O Excel aplica esse formato automaticamente
// em muitas colunas numéricas, e a planilha do usuário pode ter isso mesmo
// numa coluna que é só um código/identificador, não um valor monetário.
// Não casa com texto solto (ex.: "Total: 15") nem com números decimais.
var numeroComSeparadorDeMilhar = regexp.MustCompile(`^\d{1,3}([.,]\d{3}|\x{00a0}\d{3}| \d{3})+$`)

// limparSeparadorDeMilhar remove o separador de milhar de um número que
// bate com numeroComSeparadorDeMilhar, devolvendo só os dígitos. Qualquer
// outro valor (incluindo números decimais de verdade ou texto) passa
// intocado.
func limparSeparadorDeMilhar(s string) string {
	if !numeroComSeparadorDeMilhar.MatchString(s) {
		return s
	}
	return strings.NewReplacer(",", "", ".", "", " ", "", " ", "").Replace(s)
}

// somenteNumeros verifica se a string contém apenas dígitos — usado para
// descartar células que não são um código de fatura de verdade, como uma
// linha de rodapé/total ("Qtd: 15") que às vezes aparece na mesma coluna.
func somenteNumeros(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func readCSV(file []byte) ([][]string, error) {
	reader := csv.NewReader(bytes.NewReader(file))
	reader.FieldsPerRecord = -1 // permite linhas com número de colunas diferente

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o csv: %w", err)
	}
	return rows, nil
}

func readXLSX(file []byte) ([][]string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir o arquivo xlsx: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return nil, errors.New("planilha não possui nenhuma aba")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler as linhas da planilha: %w", err)
	}
	return rows, nil
}

// ExtrairValoresFatura lê um arquivo .csv ou .xlsx, localiza a coluna cujo
// cabeçalho corresponde a "fatura" (ou variações como "codFatura",
// "Número Fatura" etc.) e devolve os valores dessa coluna, um por posição
// do slice, na ordem em que aparecem na planilha. Formatação (com ou sem
// aspas, separador etc.) fica por conta de quem consome o resultado.
func ExtractValues(file []byte, nomeArquivo string) ([]string, error) {
	var rows [][]string
	var err error

	nome := strings.ToLower(nomeArquivo)

	switch {
	case strings.HasSuffix(nome, ".csv"):
		rows, err = readCSV(file)
	case strings.HasSuffix(nome, ".xlsx"), strings.HasSuffix(nome, ".xls"):
		rows, err = readXLSX(file)
	default:
		return nil, fmt.Errorf("formato de arquivo não suportado: %s (use .csv ou .xlsx)", nomeArquivo)
	}

	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, errors.New("planilha vazia ou sem linhas de dados após o cabeçalho")
	}

	headerRow, colIndex, err := findHeaderRow(rows)
	if err != nil {
		return nil, err
	}

	values := make([]string, 0, len(rows)-headerRow-1)
	for _, row := range rows[headerRow+1:] {
		if colIndex >= len(row) {
			continue
		}

		v := strings.TrimSpace(row[colIndex])
		v = limparSeparadorDeMilhar(v)
		if !somenteNumeros(v) {
			continue
		}

		values = append(values, v)
	}

	if len(values) == 0 {
		return nil, errors.New("nenhum valor encontrado na coluna de fatura")
	}

	return values, nil
}
