# Microblogging Platform

Esta es una plataforma de microblogging desarrollada en Go con PostgreSQL como base de datos. La aplicación está completamente dockerizada y se puede ejecutar con un simple comando.

## Requisitos Previos

- Docker
- Docker Compose
- Go (para desarrollo local)

## Instrucciones para levantar el proyecto

### Opción 1: Con Docker

Para levantar toda la aplicación (API + Base de datos PostgreSQL), ejecuta:

```bash
docker compose up -d
```

Este comando:
1. Construye la imagen de la aplicación Go
2. Levanta un contenedor PostgreSQL con la base de datos
3. Ejecuta las migraciones automáticamente
4. Inicia la aplicación conectada a PostgreSQL

### Opción 2: Desarrollo Local

1. Clona el repositorio en tu máquina local.
2. Dentro del directorio raíz del proyecto `microblogging-pltf` navega hasta la carpeta `cmd` y ejecuta el archivo main.go utilizando el comando `go run main.go`.

## Servicios

### API
- **Puerto (Docker)**: `8082`
- **Puerto (Local)**: `8080`
- **Base URL (Docker)**: `http://localhost:8082`
- **Base URL (Local)**: `http://localhost:8080`

### PostgreSQL (Solo Docker)
- **Contenedor**: `microblog_postgres`
- **Puerto**: `5432`
- **Base de datos**: `microblog`
- **Usuario**: `microblog_user`
- **Contraseña**: `microblog_password`

## Documentación

En la carpeta docs se encuentra un png que ilustra de arquitectura general de la plataforma de microblogging como solución escalable y optimizada para lecturas a implementar.

## Endpoints y Ejemplos de Uso

A continuación se detallan los endpoints disponibles con ejemplos usando curl:

### 1. Verificar conexión (Ping)
```bash
# Docker
curl -X GET http://localhost:8082/ping

# Local
curl -X GET http://localhost:8080/ping
```

### 2. Crear Usuario
```bash
# Docker
curl -X POST http://localhost:8082/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan Pérez","email":"juan@example.com"}'

# Local
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan Pérez","email":"juan@example.com"}'
```

### 3. Obtener Usuario por ID
```bash
# Docker
curl -X GET http://localhost:8082/api/v1/users/{userID}

# Local
curl -X GET http://localhost:8080/api/v1/users/{userID}
```

### 4. Publicar Tweet
```bash
# Docker
curl -X POST http://localhost:8082/api/v1/users/{userID}/tweet \
  -H "Content-Type: application/json" \
  -d '{"message":"Mi primer tweet desde Docker!"}'

# Local
curl -X POST http://localhost:8080/api/v1/users/{userID}/tweet \
  -H "Content-Type: application/json" \
  -d '{"message":"Mi primer tweet!"}'
```

### 5. Seguir a un Usuario
```bash
# Docker
curl -X POST http://localhost:8082/api/v1/users/{followerID}/follow/{followedID}

# Local
curl -X POST http://localhost:8080/api/v1/users/{followerID}/follow/{followedID}
```

### 6. Obtener Timeline de Usuario
```bash
# Docker
curl -X GET http://localhost:8082/api/v1/users/{userID}/timeline

# Local
curl -X GET http://localhost:8080/api/v1/users/{userID}/timeline
```

## Comandos Útiles de Docker

### Ver logs de la aplicación
```bash
docker compose logs app
```

### Ver logs de PostgreSQL
```bash
docker compose logs postgres
```

### Detener los servicios
```bash
docker compose down
```

### Detener y eliminar volúmenes (elimina datos)
```bash
docker compose down -v
```

### Reconstruir la aplicación
```bash
docker compose build app
docker compose up -d
```

## Estructura de la Base de Datos

La base de datos se inicializa automáticamente con las siguientes tablas:

- **users**: Almacena información de usuarios
- **tweets**: Almacena los tweets de los usuarios
- **followers**: Relación de seguimiento entre usuarios

## Configuración

La aplicación utiliza las siguientes variables de entorno:

- `DATABASE_URL`: Cadena de conexión a PostgreSQL
- `PORT`: Puerto donde corre la aplicación (8082 en Docker, 8080 en local)

### Reiniciar desde cero
```bash
docker compose down -v
docker compose up -d
```

## Consideraciones Técnicas

Para la estructura del proyecto utilicé arquitectura hexagonal.
La aplicación posee los servicios básicos del back-end que le permitirían a un usuario publicar tweets, seguir a otros usuarios y ver el timeline de tweets.

Para simplificar se implementó una base de datos in memory, sin embargo en el documento de arquitectura general de una aplicación escalable se especifica el tipo de base de datos que usaría.
También se implementó una DB PostgreSQL que funciona completamente con Docker.

## Desarrollo

Para desarrollo local, la aplicación tiene un fallback a base de datos en memoria si no se proporciona `DATABASE_URL`.
