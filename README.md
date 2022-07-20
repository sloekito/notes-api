### Design
- Used a popular library called mux
- Used go-memdb
- A note can have a title and text

### Run Code
```
 go run cmd/notes-api/main.go
```

### Test Code
```
curl -X POST \
  http://localhost:8000/notes \
  -d '{
    "Title": "Summer Time",
    "Text": "And the livin'\'' is easy"
}'
```

```
curl -X GET \
  http://localhost:8000/notes \
```

```
curl -X PUT \
  http://localhost:8000/notes/<id> \
  -d '{
    "Title": "Winter Time",
    "Text": "And the livin'\'' is not easy"
  }'
```

```
curl -X DELETE \
  http://localhost:8000/notes/685e8f52-5ec3-46ec-8d6b-5043fce4002d 
```


### Improvement Ideas
#### Unit tests

- Given more time this would be the first thing to focus on
- Can split database calls to repository layer
- Or make mock objects to simulate memdb
- Can also test by running in docker

#### Get Notes endpoint
- Can support pagination for get notes endpoint
- Implement limits on title and text length.
- Implement other ways to query notes
- Limit number of items returned in GET query
- Another endpoint can also be implemented to get a single or multiple notes by ID or other queries

#### Create notes endpoint
- Only support single note creation

#### Update Notes endpoint
- At the moment does not support single field updates. For example we have to update both the title and text fields to get the desired result.
- A PATCH endpoint may be more appropriate for some of those operations.

#### Delete Note endpoint
- Intentionally return 200 OK whether an operation is successful or not, but can be changed to fit requirements


####  Overall system improvements
- Restructure code to support unit tests and integration tests
- Better error handling, 4xx, 5xx
- Database limits
- Documentation
- Input validation


### How to deploy and scale
- Compile code, run in docker
- Can be run in Kubernetes or other container management systems 
- Database should be separated into persistent layer if running at scale
- Secure behind authentication and authorization (can be a JWT auth)