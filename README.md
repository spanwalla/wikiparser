# Сборка
```bash
cd cmd/main
go build
```
Если не получается собрать или запустить готовый файл, попробуйте использовать опцию `-builmode=exe`.

# Ключи
Программа имеет три сценария работы:
1. parse: прочитать данные из дампа формата .sql.gz и перенести их в формат .csv, разделитель `\t`.
2. replace: заменить названия страниц в файлах redirect, pagelinks (всё, где второй колонкой идёт название страницы) на ID соответствующих страниц. Обратите внимание, что файл будет перезаписан.
3. combine: сгруппировать информацию о входящих и исходящих ссылках на страницу в формате .csv (ID страницы\tLIST[INCOMING LINKS]\tLIST[OUTGOING LINKS]). Список ссылок представлен целыми числами, идентификаторами страниц, разделёнными запятой. Если ссылки отсутствуют, указано '-'.

`-input [string]` - ОБЯЗАТЕЛЬНО - путь к файлу ввода (.sql.gz в первом сценарии и файл pagelinks во втором и третьем (расширение .sql.csv, разделитель `\t`).

`-mode [parse/replace/combine]` - ОБЯЗАТЕЛЬНО - режим работы программы.

`-path [string]` - ОБЯЗАТЕЛЬНО В РЕЖИМЕ replace - путь к файлу с информацией о страницах (..-pages.sql.csv).

`-silent` - ОПЦИОНАЛЬНО - тихий режим при работе. Если ключ указан, прогресс не будет выводиться в терминал (однако сообщения об ошибках всё ещё будут). Также принимает значения [true, false, T, F, t, f, 1, 0, ...].

# Описание
Программа для чтения данных дампа Википедии, а именно таблиц page, pagelinks, redirect (остальные легко добавляются модификацией файла [data_structures.go](cmd/parser/data_structures.go)).
Результат работы сохраняется в каталоге output. Правила хорошего тона предписывают вам использовать каталог resources для размещения входных данных для программы.

**Я не ручаюсь за поведение программы в исключительных ситуациях: повреждённый файл, неподдерживаемый тип, неправильное выражение, ручные изменения файлов. Скорее всего программа завершится с ошибкой.**
