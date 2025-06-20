# Article CMS

This repository contains a cms service built using golang

## Installation Using Docker Compose

After cloning this repository, please create a .env file and copy-paste the content from .env.example file

And then, execute this command at the project root directory:

```bash
docker compose up -d
```

## Features

- Swagger provided: All endpoints can be tried by accessing the swagger docs: http://localhost:9000/swagger/index.html
- Unit Test: each services has average of more than 90% coverage 
- Authentication using JWT
- Dockerized
- Integration testing 

## DB Schema
Here is the db schema used in this service:
![ERD CMS](erd-cms.png)