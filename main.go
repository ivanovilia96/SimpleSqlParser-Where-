package main

import (
	"errors"
	"strconv"
	"strings"
)

var (
	incorrectQueryError                                                           = errors.New("your query is incorrect")
	thereIsAnKeywordsFromAnotherQueryPart                                         = errors.New("your query has an keyword like LIMIT or ORDER BY, which can not be in this part of query, or IN statement which dont support now")
	comparisonAndLogicalOperationsWhichCanBeInWhereClauseRightValueLeftColumnName = []string{"=", ">", "<", ">=", "<=", "<>", "!=", "BETWEEN", "LIKE", "IN"}
	comparisonAndLogicalOperationsWhichCanBeInWhereClauseRightColumnNameLeftValue = []string{"AND", "OR"}
	comparisonAndLogicalOperationsWhichCanBeInWhereClause                         = []string{"=", ">", "<", ">=", "<=", "<>", "!=", "AND", "OR", "IN", "BETWEEN", "LIKE"}
	syntaxErrorQuotes                                                             = errors.New("maybe you forgot add \" ' \" to start or end of your word")
)

// IN keyword don`t support now

type Parse struct {
	sqlWhereQuery      string
	tokensListWithInfo []StatisticElement
	columnInfo         map[string]string
	// этот Where это из задания в случае успеха разбора и проверки входных условий - формировать WHERE условия для qb squirrel.SelectBuilder вида qb =
	//qb.Where(squirrel.Eq{left: val}) то есть left - columnName val - val
	Where map[string]string
}

type StatisticElement struct {
	tokenType string
	dataType  string
	value     string
}

// разбирает на токены исходную строку + первичная проверка на точно не имеющие к части запроса where ключевые слова
func (v *Parse) ParseQueryOnTokens() error {
	splittedQuery := strings.Split(strings.TrimSpace(strings.ToUpper(v.sqlWhereQuery)), " ")
	for _, value := range splittedQuery {
		if value == "LIMIT" || value == "ORDER" || value == ",LIMIT" || value == ",ORDER" || value == "LIMIT," || value == "ORDER," {
			return thereIsAnKeywordsFromAnotherQueryPart
		} else if len(value) != 0 {
			thereIsThisToken := false
			// хотим проверить, нет ли в строке знаков сравнения (кроме <>)
			if indexOfFoundElement, isThere := Find(strings.Split(value, ""), "="); isThere {

				if indexOfFoundElement >= 1 {
					if value[indexOfFoundElement-1] == '=' || value[indexOfFoundElement-1] == '>' || value[indexOfFoundElement-1] == '<' || value[indexOfFoundElement-1] == '!' {
						if len(value[:indexOfFoundElement-1]) > 0 {
							v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[:indexOfFoundElement-1]})
						}
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement-1 : indexOfFoundElement+1]})
						if len(value[indexOfFoundElement+1:]) > 0 {
							v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement+1:]})
						}
						thereIsThisToken = true
					} else {
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[:indexOfFoundElement]})
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement : indexOfFoundElement+1]})
						if len(value[indexOfFoundElement+1:]) > 0 {
							v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement+1:]})
						}
						thereIsThisToken = true
					}
				} else {
					v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[:1]})
					if len(value[1:]) != 0 {
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[1:]})
					}
					thereIsThisToken = true
				}

			}
			// хотим проверить, нет ли в строке знака сравнения  <> и ни является ли он >
			if indexOfFoundElement, isThere := Find(strings.Split(value, ""), ">"); isThere {
				if indexOfFoundElement >= 1 {
					if value[indexOfFoundElement-1] == '<' {
						if len(value[:indexOfFoundElement-1]) > 0 {
							v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[:indexOfFoundElement-1]})
						}
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement-1 : indexOfFoundElement+1]})
						if len(value[indexOfFoundElement+1:]) > 0 {
							v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement+1:]})
						}
						thereIsThisToken = true
					} else {
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[:indexOfFoundElement]})
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement : indexOfFoundElement+1]})
						if len(value[indexOfFoundElement+1:]) > 0 {
							v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[indexOfFoundElement+1:]})
						}
						thereIsThisToken = true
					}
				} else {
					v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[:1]})
					if len(value[1:]) != 0 {
						v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value[1:]})
					}
					thereIsThisToken = true
				}
			}
			if !thereIsThisToken {
				v.tokensListWithInfo = append(v.tokensListWithInfo, StatisticElement{value: value})
			}

		}
	}

	if len(v.tokensListWithInfo) < 3 {
		return incorrectQueryError
	}

	return nil
}

// функция определяет тип данных элемента и тип токена (r\l value, specWord), возвращает выявленную информацию о токенах
func splitTokensOnTypes(tokens []StatisticElement) (error, []StatisticElement) {
	var sortedTokens []StatisticElement
	for _, value := range tokens {
		tokenValue := value.value
		// определяем к какому типу относится токен (л-валуе или р-валуе или ключевое слово какое-то)
		//check for int r-value
		if _, err := strconv.Atoi(tokenValue); err == nil {
			sortedTokens = append(sortedTokens, StatisticElement{"rValue", "int", tokenValue})
			//check for string r-value\l-value and non-special word
		} else if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClause, tokenValue); !isThere {
			// проверка на кавычки если они есть с 1 стороны. то должны быть и с другой, иначе синтаксическая ошибка
			if tokenValue[0] == '\'' || tokenValue[len(tokenValue)-1] == '\'' {
				if tokenValue[0] != '\'' || tokenValue[len(tokenValue)-1] != '\'' {
					panic(syntaxErrorQuotes.Error())
				}
			}
			// если есть ковычки то это r-value
			if tokenValue[0] == '\'' && tokenValue[len(tokenValue)-1] == '\'' {
				sortedTokens = append(sortedTokens, StatisticElement{"rValue", "string", tokenValue})
			} else {
				// если нет ковычек то это l-value
				sortedTokens = append(sortedTokens, StatisticElement{"lValue", "string", tokenValue})
			}
		} else {
			// предпологаем что все оставшиеся токены- ключевые слова
			// определяем что должно быть слева\справа от ключевого слова
			if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClauseRightValueLeftColumnName, tokenValue); isThere {
				// значит токен относится к типу слева колонка справа значение
				sortedTokens = append(sortedTokens, StatisticElement{"specWord", "LeftColumnRightValue", tokenValue})
			} else if _, isThere := Find(comparisonAndLogicalOperationsWhichCanBeInWhereClauseRightColumnNameLeftValue, tokenValue); isThere {
				// значит токен относится к типу слева значение справа колонка
				sortedTokens = append(sortedTokens, StatisticElement{"specWord", "LeftValueRightColumn", tokenValue})
			} else {
				panic(thereIsAnKeywordsFromAnotherQueryPart.Error())
			}
		}
	}
	return nil, sortedTokens
}

// функция проверят синтаксис нашей where части sql запроса
func syntaxCheck(tokens []StatisticElement) {
	//первое слово не должно быть ни ключевым ни rValue
	if len(tokens) != 0 && tokens[0].tokenType != "lValue" {
		panic(incorrectQueryError.Error() + " error is near " + tokens[0].value + " word")
	}

	for index, value := range tokens {
		//fmt.Printf("%v - index, %s - value \n", index, value) -> покжет что хранится

		if index != 0 {
			// проверяем что слева от типа LeftColumnRightValue должно стоять lValue
			if value.dataType == "LeftColumnRightValue" {
				if value.value == "LIKE" {
					if len(tokens) == index-1 || tokens[index+1].dataType != "string" {
						panic(incorrectQueryError.Error() + " error with LIKE keyword")
					}
				}
				if value.value == "BETWEEN" {
					if len(tokens) == index-2 || tokens[index+2].value != "AND" || tokens[index-1].tokenType != "lValue" {
						panic(incorrectQueryError.Error() + " error with BETWEEN keyword")
					}
				} else if tokens[index-1].tokenType != "lValue" {
					panic(incorrectQueryError.Error() + " error is near " + value.value + " word")
				} else if len(tokens) == index-1 || tokens[index+1].tokenType != "rValue" {
					panic(incorrectQueryError.Error() + " error is near " + value.value + " word")
				}
				// проверяем что слева от типа LeftValueRightColumn стоит rValue
			} else if value.dataType == "LeftValueRightColumn" {
				if tokens[index-2].value == "BETWEEN" {
					if value.value != "AND" {
						panic(incorrectQueryError.Error() + "error with BETWEEN keyword")
					}
				} else if tokens[index-1].tokenType != "rValue" {
					panic(incorrectQueryError.Error() + " error is near " + value.value + " word")
				} else if len(tokens) == index+1 || tokens[index+1].tokenType != "lValue" {
					panic(incorrectQueryError.Error() + " error is near " + value.value + " word")
				}
			} else if value.dataType == "string" || value.dataType == "int" {
				// предусмотреть ситуацию когда 2 int или стринг подряд стоят
				if tokens[index-1].dataType == "string" || tokens[index-1].dataType == "int" {
					panic(incorrectQueryError.Error() + " error is near " + value.value + " word")
				}
			}
		}
	}
}

// функция проверяет типы колонок к типам из columnInfo переменной
func checkColumnTypes(tokens []StatisticElement, columnsInfo map[string]string, whereMap map[string]string) {
	//создаем локальный map в верхнем регистре что бы было регистро независимо ( элементы в tokens уже все в верхнем регистре )
	localMapUpperCase := make(map[string]string)
	for key, value := range columnsInfo {
		localMapUpperCase[strings.ToUpper(key)] = value
	}
	for index, value := range tokens {
		// проверяем на то, что такой токен присутствует в нашем массиве
		val, ok := localMapUpperCase[value.value]
		if value.tokenType == "lValue" && ok {
			//сопоставляем в Where
			if tokens[index+2].dataType == val {
				whereMap[value.value] = tokens[index+2].value
			}
			if tokens[index+2].dataType != val {
				panic("type error for column: " + value.value + ", excepted: " + val + " received " + tokens[index+2].dataType)
			}
		}
	}
}

// функция занимается сбором данных и синтаксической проверкой where части sql запроса
func (v *Parse) CheckCorrectSpelling() {
	err, sortedTokens := splitTokensOnTypes(v.tokensListWithInfo)
	if err != nil {
		panic(err.Error())
	}
	syntaxCheck(sortedTokens)
	checkColumnTypes(sortedTokens, v.columnInfo, v.Where)
	v.tokensListWithInfo = sortedTokens
}

func main() {
	// указываем то, что в определенные колонки должен быть определенный тип данных (int\string) названия колонок могут быть в любом регистре
	columnsInfo := map[string]string{
		"AlICE.name":   "int",
		"BOB.LASTNAME": "string",
		"age":          "int",
		"age2":         "string",
	}
	query := "Alice.Name=5 and Bob.LastName!='56' or  age<> 20 and  age2 like '%3' "
	// все поля из запроса нужно занести в columnsInfo, иначе они не проверятся на типы и не занесутся в firstParse.Where
	firstParse := Parse{
		query,
		[]StatisticElement{},
		columnsInfo,
		map[string]string{},
	}
	println("Часть для проверки и парсинга: " + query)
	err := firstParse.ParseQueryOnTokens()
	if err != nil {
		panic(err.Error())
	}
	firstParse.CheckCorrectSpelling()
	println("Выражение прошло валидацию")

	// этим мы проверяем что у нас есть "такого воида  токены" ключ - имя колонки , значение - значение которое мы присваиваем
	for key, val := range firstParse.Where {
		println(key, val)
	}

	// соответственно если нам требуется аналитика по каждому токену, мы можем получить её с помощью простого перебора firstParse.tokensListWithInfo
	// аналитика ( в данном случае firstParse.tokensListWithInfo ) включает в себя массив токенов у каждого из которых
	//я выделил 3 основные черты {tokenType - lValue\rValue\specWord,
	//dataType - тип данных, на данный момент есть 2: string\int - для полей где  tokenType = rValue и string для полей где tokenType = lValue,
	//value - значение токена в upperCase }

}
