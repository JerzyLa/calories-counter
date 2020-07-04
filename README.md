
### Start mysql docker

```bash
docker run --name=mysql-local -p 3306:3306 -e MYSQL_ROOT_HOST=% -e MYSQL_ROOT_PASSWORD=my-secret-password -d mysql/mysql-server:5.7
```

### Init database

run script `sql/init_db.sql`

### Setup env variables

`TOKEN_SECRET` `API_APP_ID` `API_KEY` `MYSQL_DB_SOURCE`

### Start server application

````bash
make start
````
