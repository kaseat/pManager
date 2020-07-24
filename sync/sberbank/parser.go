package sberbank

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

const reportHeaderMarker = "<br>Отчет брокера</br>"
const operationsTableMarker = "<br>Сделки купли/продажи ценных бумаг</br>"
const securitiesInfoTableMarker = "<br>Справочник Ценных Бумаг**</br>"
const cashFlowTableMarker = "<br>Движение денежных средств за период</br>"
const buybacksTableMarker = "<br>Движение ЦБ, не связанное с исполнением сделок</br>"
const endTableMarker = "</table>"

func parseReport(r io.Reader) report {
	result := report{}
	var operationsTable strings.Builder
	var securitiesInfoTable strings.Builder
	var cashFlowTable strings.Builder
	var buybackTable strings.Builder
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		// discard monthly reports
		if scanner.Text() == reportHeaderMarker {
			scanner.Scan()
			re := regexp.MustCompile(`\d{2}.\d{2}.\d{4}`)
			match := re.FindAllString(scanner.Text(), -1)
			if match[0] != match[1] {
				result.IsEmpty = true
				break
			}
			result.Date = match[0]
		}
		if scanner.Text() == operationsTableMarker {
			scanner.Scan()
			for scanner.Scan() {
				operationsTable.Write(scanner.Bytes())
				if scanner.Text() == endTableMarker {
					break
				}
			}
		}
		if scanner.Text() == securitiesInfoTableMarker {
			scanner.Scan()
			for scanner.Scan() {
				securitiesInfoTable.Write(scanner.Bytes())
				if scanner.Text() == endTableMarker {
					break
				}
			}
		}
		if scanner.Text() == cashFlowTableMarker {
			scanner.Scan()
			for scanner.Scan() {
				cashFlowTable.Write(scanner.Bytes())
				if scanner.Text() == endTableMarker {
					break
				}
			}
		}
		if scanner.Text() == buybacksTableMarker {
			scanner.Scan()
			for scanner.Scan() {
				buybackTable.Write(scanner.Bytes())
				if scanner.Text() == endTableMarker {
					break
				}
			}
		}
	}

	if !result.IsEmpty {
		result.SecuritiesInfo = processSecuritiesInfoTable(parseTable(securitiesInfoTable.String()))
		result.Operations = processOperationsTable(parseTable(operationsTable.String()), result.SecuritiesInfo)
		result.CashFlow = processCashFlowTable(parseTable(cashFlowTable.String()))
		result.Buybacks = processBuybacksTable(parseTable(buybackTable.String()), result.SecuritiesInfo)
	}

	return result
}

func parseTable(rawTable string) [][]string {
	slice := [][]string{}
	dict := map[string]pair{}
	begin, end, n := 0, 0, 0

	for i, char := range rawTable {
		if char == 60 { // rune "<"
			begin = i + 1
			n = 0
		}

		if char == 32 { // rune " "
			n++
			if n == 1 {
				end = i
			}
		}

		if char == 62 { // rune ">"
			if n == 0 {
				end = i
			}
			rawTag := rawTable[begin:end]
			if []rune(rawTag)[0] == 47 { // closing tag
				tag := rawTag[1:] // remove "/" from tag name
				bounds := dict[tag]
				bounds.end = begin - 1
				dict[tag] = bounds

				if tag == "td" { // process delimiter logic
					lastElement := len(slice) - 1
					arr := slice[lastElement]
					slice[lastElement] = append(arr, rawTable[bounds.begin:bounds.end])
				}
			} else { // opening tag
				bounds := dict[rawTag]
				bounds.begin = i + 1
				dict[rawTag] = bounds

				if rawTag == "tr" { // process delimiter logic
					slice = append(slice, []string{})
				}
			}
		}
	}
	return slice
}
