# Servidor de Broadcast Concurrente

Proyecto base para la Clase sobre Sockets de Sistemas Distribuidos.

## Integrantes

- Gallo, Guillermo Ariel
- Pedernera Theisen, Nahuel Thomas

## Ejecución

### Local

```bash
# Terminal 1: servidor
go run ./cmd/servidor

# Terminal 2: cliente
go run ./cmd/cliente
```

### Docker Compose

```bash
docker-compose up --build
```

## Requisitos completados

- [ ] Servidor TCP concurrente
- [ ] Protocolo JSON
- [x] Registro de clientes con sync.RWMutex
- [x] Broadcast a todos los clientes
- [ ] Cliente interactivo (stdin + recepción paralela)
- [ ] Docker + docker-compose
- [x] Bonus: descubrimiento UDP

## Captura de ejecución

(Adjuntar log o captura de pantalla con múltiples clientes conectados)
