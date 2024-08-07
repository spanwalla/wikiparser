# Сборка
```bash
go build ./cmd/combiner
go build ./cmd/parser
go build ./cmd/replacer
go build ./cmd/validator
```
Если не получается собрать или запустить готовый файл, попробуйте использовать опцию `-buildmode=exe`.

# Ключи
В комплекте имеется четыре программы:
1. parser: прочитать данные из дампа формата .sql.gz и перенести их в формат .csv, разделитель `\t`.
2. replacer: заменить названия страниц в файлах redirect, pagelinks (всё, где второй колонкой идёт название страницы) на ID соответствующих страниц. Обратите внимание, что файл будет перезаписан. Также убирает ссылки на несуществующие статьи.
3. validator: заменить ссылки-перенаправления на истинные страницы.  Файл redirect должен быть в формате `sourcePageId\ttargetPageId`.
4. combiner: сгруппировать информацию о входящих и исходящих ссылках на страницу в формате .csv (ID страницы\tLIST[INCOMING LINKS]\tLIST[OUTGOING LINKS]). Список ссылок представлен целыми числами, идентификаторами страниц, разделёнными запятой. Если ссылки отсутствуют, указано '-'.

`-input [string]` - ОБЯЗАТЕЛЬНО ДЛЯ parser И replacer - путь к файлу ввода (.sql.gz в первом сценарии и файл .sql.csv во втором (разделитель `\t`)).

`-pagelinks [string]` - ОБЯЗАТЕЛЬНО ДЛЯ combiner И validator - путь к файлу с информацией о ссылках (..-pagelinks.sql.csv).

`-redirect [string]` - ОБЯЗАТЕЛЬНО ДЛЯ validator - путь к файлу с информацией о перенаправлениях (..-redirect.sql.csv)

`-page [string]` - ОБЯЗАТЕЛЬНО ДЛЯ replacer - путь к файлу с информацией о страницах (..-pages.sql.csv).

`-silent` - ОПЦИОНАЛЬНО - тихий режим при работе. Если ключ указан, прогресс не будет выводиться в терминал (однако сообщения об ошибках всё ещё будут). Также принимает значения [true, false, T, F, t, f, 1, 0, ...]

# Описание
Программа для чтения данных дампа Википедии, а именно таблиц page, pagelinks, redirect (остальные легко добавляются модификацией файла [parse.go](internal/parse/parse.go)).
Результат работы сохраняется в каталоге output. Правила хорошего тона предписывают вам использовать каталог resources для размещения входных данных для программы.
