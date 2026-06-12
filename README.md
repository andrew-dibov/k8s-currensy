# Exchange Service

Сервис конвертации валют, состоящий из микросервисов на Go с кафкой, редисом и постгрессами

Планируется три типа развертывания :
- A : Docker Compose : готов
- B : Docker Swarm : в процессе
- C : Kubernetes : в процессе

## Docker Compose

```bash
(cd docker-compose && docker compose up -d)
./test.sh # покурлить запросы
```

## Docker Swarm

В ПРОЦЕССЕ

```bash
docker swarm init
(cd docker-swarm && ???)
```

## Kubernetes : microk8s

```bash
```

В ПРОЦЕССЕ