# Misrious pong
Misra PING-PONG algorithm GO implementation, using TCP server.

### Flags
Configuration options available in config.go

| flag          | description                                                                                           | default value | example       |
|---------------|-------------------------------------------------------------------------------------------------------|---------------|:--------------|
| mode          | Application mode. Allowed values SERVER, INITIATOR (case insensitive)                                 | server        | initiator     |
| sleep         | Number of seconds to sleep between sending messages - imitates processing in the critical section.    | 5             | 1             |
| serverAddress | Server listen address                                                                                 | localhost     | 127.0.0.1     |
| serverPort    | Server listen port                                                                                    | 5999          | 8080          |
| clientAddress | Destination server address - the address of the destination server where all the message will be sent | localhost     | 192.168.100.1 |
| clientPort    | Destination server port                                                                               | 5998          | 8081          |
