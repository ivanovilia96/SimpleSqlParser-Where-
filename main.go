package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	incorrectQueryError                                             = errors.New("your query is incorrect")
	thereIsAnKeywordsFromAnotherQueryPart                           = errors.New("your query has an keyword like LIMIT or ORDER BY, which can not be in this part of query")
	incorrectStartOfQuery                                           = errors.New("your query can not start or end with AND or OR clause ")
	comparisonAndLogicalOperationsWhichCanBeInWhereClauseEasyStruct = []string{"=", ">", "<", ">=", "<=", "<>", "!=", "BETWEEN"}
	comparisonAndLogicalOperationsWhichCanBeInWhereClause           = []string{"=", ">", "<", ">=", "<=", "<>", "!=", "AND", "OR", "IN", "BETWEEN", "LIKE", "NOT"}
	syntaxError                                                     = errors.New("syntax error")
)

type Parse struct {
	sqlWhereQuery string
	sqlParts      []string
}

func (v *Parse) ParseQueryOnTokens() error {
	splittedQuery := strings.Split(strings.TrimSpace(strings.ToUpper(v.sqlWhereQuery)), " ")
	for _, value := range splittedQuery {
		if value == "LIMIT" || value == "ORDER" || value == ",LIMIT" || value == ",ORDER" || value == "LIMIT," || value == "ORDER," {
			return thereIsAnKeywordsFromAnotherQueryPart
		} else if len(value) != 0 {
			v.sqlParts = append(v.sqlParts, value)
		}
	}

	if len(v.sqlParts) < 3 {
		return incorrectQueryError
	}

	for _, v := range v.sqlParts {
		println(v, " value")
	}
	return nil
}

func (v *Parse) CheckCorrectSpelling() {
	if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, v.sqlParts[0]); isThere {
		panic(incorrectStartOfQuery.Error())
	}

	if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, v.sqlParts[len(v.sqlParts)-1]); isThere {
		panic(incorrectStartOfQuery.Error())
	}

	if v.sqlParts[0][0] == '\'' || v.sqlParts[0][len(v.sqlParts[0])-1] == '\'' {
		panic("Первое слово  не может быть строкой")
	}

	for index, value := range v.sqlParts {
		if index != 0 {
			// проверка на то , что если элемент является одним из ключевых слов, то не перед ним ни до не должно быть ключевых слов
			if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, value); isThere {
				_, isThereBefore := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, v.sqlParts[index-1])
				_, isThereAfter := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, v.sqlParts[index+1])

				if isThereBefore || isThereAfter {
					panic(syntaxError.Error() + "near " + value + " word of query")
				}
			}

			if value[0] == '\'' || value[len(value)-1] == '\'' {
				// если это какая то строка, то до неё должно быть ключевое слово
				if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClauseEasyStruct, v.sqlParts[index-1]); !isThere {
					panic(syntaxError.Error() + "near : " + value + " word of query")
				}
				// строка должна и начинаться и заканчиваться со специального символа, иначе синтаксическая ошибка
				if value[0] != '\'' || value[len(value)-1] != '\'' {
					panic(syntaxError.Error() + "near : " + value + " word of query, may be you forgot add the quotes to the string value")
				}
			}
			//  если это какое точисло, то до него должно быть ключевое слов
			if _, err := strconv.Atoi(value); err == nil {

				if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, v.sqlParts[index-1]); !isThere {
					panic(syntaxError.Error() + "near - " + value + " word of query")
				}
			}
			// перед не специальным словом не может быть не специальное слово
			_, isCurrentSpecSymbol := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, value)
			if _, isPreviousSpecSymbol := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, v.sqlParts[index-1]); !isPreviousSpecSymbol &&
				!isCurrentSpecSymbol {
				panic(syntaxError.Error() + "near - " + value + " word of query")
			}

		}
	}
}

func main() {
	firstParse := Parse{"name <> 'ilia' and age != 15", []string{}}
	err := firstParse.ParseQueryOnTokens()
	if err != nil {
		panic(err.Error())
	}
	firstParse.CheckCorrectSpelling()

	print("Выражение прошло валидацию")
}
