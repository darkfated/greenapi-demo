# GreenAPI Demo
HTML-страница для взаимодействия с GreenAPI

## Скриншот сайта
<img width="1903" height="943" alt="image" src="https://github.com/user-attachments/assets/16c477aa-cd9c-468a-b242-adcf9062737b" />

## Как пользоваться
1. Создать инстанс в [личном кабинете GREEN-API](https://console.green-api.com/).
2. Сканировать QR-код в приложении WhatsApp для подключения своего номера.
3. На открытой странице ввести `idInstance` и `ApiTokenInstance` (их можно посмотреть в кабинете).
4. Нажимать кнопки:
  - **getSettings** — посмотреть настройки инстанса
  - **getStateInstance** — узнать текущее состояние
  - **sendMessage** — отправить сообщение (указать номер и текст)
  - **sendFileByUrl** — отправить файл по ссылке

Ответ методов появится в правом большом поле.

## Инициализация
1. Перейти в репозиторий
```bash
    git clone https://github.com/darkfated/greenapi-demo
    cd greenapi-demo
```
2. Запустить программу
```bash
    go run main.go
```
3. Перейти на `http://127.0.0.1:8080/`
