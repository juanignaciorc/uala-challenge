# Microblogging-pltf

## Instrucciones para levantar el proyecto

1. Clona el repositorio en tu máquina local.
2. Dentro del directorio raíz del proyecto `microblogging-pltf` navega hasta la carpeta `cmd`  y ejecuta el archivo main.go utilizando el comando ```go run main.go``` .

## Documentación

En la carpeta docs se encuentra un png que ilustra de arquitectura general de la plataforma de microblogging como solución escalable y optimizada para lecturas a implementar.

## Endpoints

A continuación se detallan los endpoints disponibles para la interacción con la aplicación:

1. **Ping**
    - Método: `GET`
    - URL: `localhost:8080/ping`
    - Descripción: Endpoint para verificar la conexión con el servidor.

2. **Crear Usuario**
    - Método: `POST`
    - URL: `localhost:8080/api/v1/users`
    - Headers:
        - Content-Type: `application/json`
    - Body:
      ```json
      {
          "name": "juan",
          "email": "juan@hotmail.com"
      }
      ```
    - Descripción: Crea un nuevo usuario con el nombre y correo electrónico proporcionados.

3. **Obtener Usuario por ID**
    - Método: `GET`
    - URL: `localhost:8080/api/v1/users/{userID}`
    - Descripción: Obtiene la información de un usuario específico mediante su ID.

4. **Publicar Tweet**
    - Método: `POST`
    - URL: `localhost:8080/api/v1/users/{userID}/tweet`
    - Headers:
        - Content-Type: `application/json`
    - Body:
      ```json
      {
          "message": "juan tweet"
      }
      ```
    - Descripción: Permite a un usuario publicar un nuevo tweet.

5. **Seguir a un Usuario**
    - Método: `POST`
    - URL: `localhost:8080/api/v1/users/{followerID}/follow/{followedID}`
    - Descripción: Establece que el usuario con ID `followerID` sigue al usuario con ID `followedID`.

6. **Obtener Timeline de Usuario**
    - Método: `GET`
    - URL: `localhost:8080/api/v1/users/{userID}/timeline`
    - Descripción: Obtiene el timeline de tweets del usuario especificado por su ID. Tener en cuenta que para simplificar solo se obtienen los tweets de los usuarios que sigue el usuario especificado.


## Consideraciones 

Para la estructura del proyecto utilicé arquitectura hexagonal.
La aplicación posee los servicios básicos del back-end que le permitirian a un usuario publicar tweets, seguir a otros usuarios y ver el timeline de tweets.

Para simplificar se implemento una data base in memory sin embargo en el documento de arquitectura general de una aplicación escalable se especifica el tipo de base de datos que usaria.
Tambien se implemento una DB PostgreSQL aunque no quedo productiva ni finalizada.  