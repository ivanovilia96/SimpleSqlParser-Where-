SimpleSqlParser-Where-
Простой парсер where части sql запросов на golang
->(не совсем уверен что это то, что вы хотели, не очень понятно что там куда) на данный момент это фактически " проверка синтаксиса части where sql выражения " то есть если вы передали синтаксически верную конструкцию, то в конце выведется "выражение прошло валидацию" если нет, то будет или обьяснение почему или где ошибка.


-> сейчас он не доделал, это как бы идея, он сейча разбирает общий синтаксис, подряд два ключевых слова нельзя поставить, limit,order by нельзя и тд, не работает с IN, LIKE(еще) , ну это сыро еще, просто идея ( (сейчас это как в окно workbench sql заходишь и пишешь что то не то, и тебе высвечивается около чего ошибка, вот я что-то такое сейчас создаю, я вообще то делаю?)


Вопросы: 
1 - " callback-обработчик для проверки валидности имени колонки и соответствия типов колонки и значения" - какие колонки, откуда брать? нужно подключать бд? и там создавать? -> если да , то запрос все равно не будет отправляться пока мы его в ручную не напишем, выражение where же может быть не только в select запросе, если нет, могу выдумать колонки в map структуру их собрать и выдуманно проверять что то, но тогда мы с ограничениями некими сталкнемся

2- "qb squirrel.SelectBuilder" ( еще есть метод Eq) это что ? метод структуры? он должен в себе выражения хранить, под выражениями понимается что ? "age <15"? надо распарсить на такие блоки? или как? 

3 - я спросил вопрос, мне ответили "Конечная цель - заново собранная строка SQL запроса к PostgreSQL из результата парсинга входной строки." - зачем её заново собирать если она валидная? входная строка же и есть where часть sql запроса (по вашему заданию: уметь разбирать пользовательский ввод вида Field1 = "foo" AND Field2 != 7 OR Field3 > 11.7 ) то есть она принимает where и потом что собрать? не понятно это совсем.


----------------------------------

если что, что надо я сделаю, просто не совсем понятен функционал конкретный который должен быть, в каком формате данные входят, в каком ожидается выход и тд, ввод данных происходит с cmd или просто вшить их в код?
