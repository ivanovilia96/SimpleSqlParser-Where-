SimpleSqlParser-Where-
Простой парсер where части sql запросов на golang
->(не совсем уверен что это то, что вы хотели) 
на данный момент что делает -> 

1 - Проверяет синтаксис строки , которая является "where" частью запроса ( там 1 ключевое слова не поддерживается пока что)
2 - Проверяет типы данных с указанными в Мэп columnsInfo
3 - Хранит информацию о каждом токене в переменной firstParse.columnInfo (значение\тип данных\ что за токен(rValue\lValue\specWord))
4 - Хранит информацию как в задании типа {lVal\rVal} {имя колонки : значение колонки } в переменной firstParse.Where

(вроде реализовал то, что требовалось в задании)

Вопросы: 
1 - " callback-обработчик для проверки валидности имени колонки и соответствия типов колонки и значения" - какие колонки, откуда брать? нужно подключать бд? и там создавать? -> если да , то запрос все равно не будет отправляться пока мы его в ручную не напишем, выражение where же может быть не только в select запросе, если нет, могу выдумать колонки в map структуру их собрать и выдуманно проверять что то, но тогда мы с ограничениями некими сталкнемся

2- "qb squirrel.SelectBuilder" ( еще есть метод Eq) это что ? метод структуры? он должен в себе выражения хранить, под выражениями понимается что ? "age <15"? надо распарсить на такие блоки? или как? 

3 - я спросил вопрос, мне ответили "Конечная цель - заново собранная строка SQL запроса к PostgreSQL из результата парсинга входной строки." - зачем её заново собирать если она валидная? входная строка же и есть where часть sql запроса (по вашему заданию: уметь разбирать пользовательский ввод вида Field1 = "foo" AND Field2 != 7 OR Field3 > 11.7 ) то есть она принимает where и потом что собрать? не понятно это совсем.


----------------------------------

если что, что надо я сделаю, просто не совсем понятен функционал конкретный который должен быть, в каком формате данные входят, в каком ожидается выход и тд, ввод данных происходит с cmd или просто вшить их в код?
