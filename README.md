# P2P words [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/izveigor/p2p-words/blob/main/LICENSE) ![Http tests](https://github.com/izveigor/p2p-words/actions/workflows/http-tests.yml/badge.svg)

Одноранговая сеть для поиска слов.
![1](https://user-images.githubusercontent.com/68601180/186383274-18883d3e-f83f-4244-a182-278553537795.png)
![2](https://user-images.githubusercontent.com/68601180/186383342-b40ef881-0bcd-4dda-941d-7bdc7d5ffa3f.png)

## Возможности
- Загрузка локального сайта для удобного пользования
- Загрузить книгу на сайт (сохраняется в локальной папке в проекте)
- Найти предложения, содержащие слово, которое было введено в поисковой
строке, с помощью одноранговой сети и локальному поиску загруженных книг
## Установка
Для демонстрации локального поиска:
1) Скачать исходник кода (с помощью приложения или сайта Github или с помощью команды git clone)
2) В главной папке запустить docker-compose
3) Зайти на сайт http://localhost:8000.
4) Загрузить небольшую книгу
5) Зайти во вкладку "Поиск" и ввести слово (если оно было найдено, то высветятся предложения, содержащее это слово).

Для демонстрации одноранговой сети (localhost):
1) Сделать все вышеперечисленные пункты
2) Скопировать проект приложения в другую папку
3) Изменить http/pkg/config/envs/dev.env:
```
P2P_SVC_URL=localhost:50054
```
4) Изменить http/cmd/main.go:
```
	if err := http.ListenAndServe(":8001", router); err != nil {
		panic(err)
	}
```
5) Изменть lemmatizer/main.py:
```
model = spacy.load("ru_core_news_sm")
MAX = 2**16
ALPHABET = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
P2P_SVC_URL = "localhost:50053"
```
6) Изменить network/pkg/config/envs/dev.env:
```
LEMMATIZER_SVC_URL=localhost:50053
P2P_SVC_URL=localhost:50054
PATH_TO_DATA=./pkg/p2p/data/
P2P_ADDRESS=localhost
INITIAL_PORT=2000
```
7) Запустить три команды на разных cmd:
- Первая команда:
```
$ cd http
$ make start
```
- Вторая команда:
```
$ cd lemmatizer
$ make start
```
- Третья команда:
```
$ cd network
$ make start
```
8) После всех правильных действий, ваше приложения на трех командных строках соединится с приложением на docker-compose. Попробуйте зайти на сайт http://localhost:8000 и http://localhost:8001 и загрузите две разные книги на одном сайте и на другом. После зайдите во вкладку "Поиск" и вбейте любое слово. В случае совпадения, у вас должны появиться предложения, содержащее это слово, как в одной книге, так и в другой книге.

Для демонстрации одноранговой сети (Интернет):
1) Зайти в network/pkg/config/envs/dev.env и измените адрес на интернет-адрес:
```
LEMMATIZER_SVC_URL=localhost:50051
P2P_SVC_URL=localhost:50052
PATH_TO_DATA=./pkg/p2p/data/
P2P_ADDRESS=АДРЕС
INITIAL_PORT=2000
```
2) Зайти на http://localhost:8000. После этого вы можете загрузить книгу, которую увидят все в этой сети, или посетить вкладку "Поиск" и найти все предложения с конкретным слово из книг, загруженных пользователями из этой сети.