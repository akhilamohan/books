## Prerequisites
* go 1.19 and make should be installed on the setup where unit tests are executed
* make and docker should be installed on the setup to run application

## Run unit tests

```
make run-uts
```

## Run the application

```
make docker-up
```
Application will be listening on port 3000

## Sample requests

### Add book 
Adds given book to the list.

#### Request
```
POST /api/books
```
```
curl -i -X POST -H  'Accept: application/json' http://localhost:3000/api/books -d '{"title": "title1", "author": "author1", "publication-date": "1991-10-04", "id": "id-123"}'
```
#### Response
```
HTTP/1.1 201 Created
Location: /api/books/id-123
Date: Fri, 10 Mar 2023 11:06:38 GMT
Content-Length: 0
```

### Get books 
Returns all books in the list
#### Request
```
GET /api/books
```
```
curl -i -X GET -H 'Accept: application/json' http://localhost:3000/api/books 
```
#### Response
```
HTTP/1.1 200 OK
Date: Fri, 10 Mar 2023 11:09:15 GMT
Content-Length: 86
Content-Type: text/plain; charset=utf-8

[{"title":"title1","author":"author1","publication-date":"1991-10-04","id":"id-123"}]
```

### Get book with id 
Return book with given id

#### Request
```
GET /api/book/<id>
```
```
curl -i -X GET -H 'Accept: application/json' http://localhost:3000/api/book/id-123
```
#### Response
```
HTTP/1.1 200 OK
Date: Fri, 10 Mar 2023 11:19:04 GMT
Content-Length: 84
Content-Type: text/plain; charset=utf-8

{"title":"title1","author":"author1","publication-date":"1991-10-04","id":"id-123"}
```

### Update book 
Update book based on id

#### Request
```
PUT /api/book/<id>
```
```
curl -i -X PUT -H 'Accept: application/json' http://localhost:3000/api/book/id-123 -d '{"title": "title1", "author": "author1", "publication-date": "1991-10-01", "id": "id-123"}' 
```
#### Response
```
HTTP/1.1 200 OK
Location: /api/books/id-123
Date: Fri, 10 Mar 2023 11:25:48 GMT
Content-Length: 0
```
### Delete book
Delete book of given id
#### Request
```
DELETE /api/book/<id>
```
```
curl -i -X DELETE -H 'Accept: application/json' http://localhost:3000/api/book/id-123 
```
#### Response
```
HTTP/1.1 200 OK
Date: Fri, 10 Mar 2023 11:31:59 GMT
Content-Length: 0
```

## Shut-down the application

```
make docker-down
```

