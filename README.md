# gpt-go
## Written By KYLiN, nexk1n 

--- 
### Project look like 
![](./img/project_look_like.png)


You can chat with ai
### Setup

You need add a file  `.env`

File format

```sh
SERVER_API_URL=http://{SERVER_DOMAIN}/{root}
SERVER_URL=http://{SERVER_DOMAIN}
```
> Recommend use our server to setup the API: [Book-To-Comic](https://github.com/KeithLin724/Book-To-Comics)


### Run gpt-go
```sh
# run the server
go run .
```

### Run in container
```sh
# build the project
docker-compose up --build -d

# shutdown the server
docker-compose down

# view the log 
docker-compose logs -f

# after build to run 
docker-compose up
```