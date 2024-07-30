# GophKeeper

### Клиентское приложение
В дир-ию с исполняемым файлом приложения скопировать файл *client-config-example.yaml* и переименовать в *client-config.yaml*
Скомпилировать приложение, указав ldflags флаги:
- main.buildDate - дата сборки
- main.buildVersion - версия приложения
- main.secretKey - секрет для шифрования данных в репозитории

### Серверное приложение
В дир-ию с исполняемым файлом приложения скопировать файл *.server.env.example* и переименовать в *.server.env*
Сгенерировать ключ и сертификат для TLS:
```
openssl req -x509 -newkey rsa:4096 -keyout ./cert/private.pem -out ./cert/cert.pem -passout pass:<ваш_пароль> -sha256 -days 365  
openssl rsa -in ./cert/private.pem -out ./cert/private.pem -passin pass:<ваш_пароль>
```
Скомпилировать приложение, указав ldflags флаги:
- main.buildDate - дата сборки
- main.buildVersion - версия приложения

#### Генерация из .proto файлов
```
protoc --go_out=internal/adapter/handler/proto \
--go_opt=paths=source_relative \
--go-grpc_out=internal/adapter/handler/proto \
--go-grpc_opt=paths=source_relative \
--proto_path=internal/adapter/handler/proto \
service.proto user.proto vault.proto
```

#### Ограничения версии PoC:
- Максимально допустимый размер текстовой заметки составляет 1 MiB
- Максимально допустимый размер загружаемого файла составляет 1 MiB
- Фиксированное именование и расположение файлов конфигурации
- Хранение содержимого файла в реляционной СУБД, а не в S3
- Регистрация и аутентификация - только online
- Функция изменения или восстановления пароля не реализована