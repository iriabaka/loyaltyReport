## Loyalty Report
### Описание:
Утилита предназначена для выгрузки отчётов из БД с последующей отправкой по Email
### Использование:
1. Подготовить `env.yml` файл
2. Запустить утилиту с указанием имени отчёта:
```shell script
sendReport -n all_users_win
```
В данном случае "all_users_win" - имя отчёта
### Описание env файла:
```yaml
# Параметры подключения к БД
datasource:
  url: 'postgresql://localhost:5432/postgres'
  user: 'postgres'
  password: 'postgres'

# Параметры подключения к почтовому серверу
mail:
  host: 'smtp.example.com'
  port: '587'
  # (Опционально) Логин и пароль для авторизации
  username: 'myUserName'
  password: 'mySecretPassword'
  # Адрес отправителя
  send_from: 'sender@example.com'

# Словарь отчётов
reports:
  # Имя отчёта
  all_users_win:
    # Запрос к БД
    query: 'select * from rg_user;'
    # Тема письма
    subject: 'all_users_report (win)'
    # Массив email адресов для отправки
    send_to: ['user1@example.com', 'user2@example.com']
    # Текст сообщения
    text: 'Test report with win-1251 encoding from golang'
    # (Опционально) Использовать кодировку win-1251 при выгрузке отчёта
    win_encoding: true
  all_users:
    query: 'select * from rg_user;'
    subject: 'all_users_report'
    send_to: ['user3@example.com', 'user4@example.com']
    text: 'Test report from golang'
```
