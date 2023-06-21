# ipc-messager
Библиотека для обмена сообщениями между сервисами через UDS(unix domain sockets).
## Описание
Библиотека предоставляет пакеты клиента и сервера.

Формат данных для обмена - raw bytes ([]byte)

### Возможности

Сервер поддерживает обработку нескольких клиентов одновременно (пример в examples).

Клиентов отличает друг от друга сочетание PID, UID, GID - сертификат.

Доступна возможность блокировки/разблокировки отдельных клиентов.

Обмен осуществляется через UDS. Используется протокол TCP.

Для обработки сообщений сервер должен предоставить функцию с определенной сигнатурой (ниже)

## Примеры использования:  

### Сервер:

      s := server.New("/tmp/server.sock", server.DEFAULT_TIMEOUT, handler)

handler - пользовательская функция, в которой обрабатываются сообщения.

Ее сигнатура - func handler(request []byte) (response []byte) 

      func handler(request []byte) (response []byte) {
      return []byte("Получил запрос")
      }

Для запуска сервера нужно обратится к s.Listen().

      <-s.Listen()

Так сервер будет работать до первой ошибки.

Есть еще один пример с таймером, в examples/server/main.go

### Клиент:

#### 1-й пример - получаем raw bytes:
      c, err := client.Connect("/tmp/server.sock", time.Millisecond*200)
      if err == nil {
            c.Send([]byte("raw bytes"))
            if data, ok := c.Receive(); ok{
                  //проводим любые операции с data - это []byte
            }
            
      }

#### 2-й пример - получаем string:

      c, err := client.Connect("/tmp/server.sock", time.Millisecond*200)
         if err == nil {
               c.Send([]byte("Hi!")) //послали строку
               if data, ok := c.Receive(); ok{
                  var buffer bytes.Buffer = *bytes.NewBuffer(data)
                  fmt.Println(buffer.String()) //ответ от сервера - в данном случае строка
               }
            } 
