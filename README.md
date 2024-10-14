# Music libraty API
Music library API can add, update and delete songs from library. It allows user to get songs filtered with query parameters, pagination is also supported
# Usage
1. Create .env file in root directory of project
```bash
# Database
DB_PORT='5432'
DB_HOST='db'
DB_NAME='music_lib'
DB_USER='anton'
DB_PASSWORD='1111'
# Web server
PORT='8080'
# External music info API
BASE_URL='http://host.docker.internal:8088'
TIMEOUT='10'
```
2. docker compose up --build -d
3. Swagger documentation on http://localhost:[PORT]/swagger
