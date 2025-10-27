# modbus-cli
## Использование
```cmd
./mb.exe -read 51 -ip 192.168.2.99 -port 502 -id 1
```
или
```cmd
./mb.exe -write 407:2 -ip 192.168.2.99 -port 502 -id 1
```
## Сборка
```cmd
go build -o mb.exe
```
Сборка под windows
```bash
GOOS=windows GOARCH=amd64 go build -o mb.exe
```
Сборка под linux
```bash
GOOS=linux GOARCH=amd64 go build -o mb_linux
```
# modbus-web
## Установка зависимостей
```cmd
pip install -r requirements.txt
```
## Использоваение
```cmd
python manage.py runserver
```
Доступ по [http://localhost:8000]