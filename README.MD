# Block User Service

This is a Golang application managing tasks using a MongoDB + Redis and Gin as Rest-API framework.


docker-compose up --build

`curl http://localhost:8080/blackli--location --request POST 'http://localhost:8080/blacklist' --header 'Content-Type: application/json' --data-raw '{
"UserID": 123,
"Token": "Well this has been some fun already!",
"BlockType":"user"}'
`  
`curl http://localhost:8080/blacklist/check/123/user`  


