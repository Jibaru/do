let {
    token = "token"
    id = 12
}

do {
    method = "GET"
    url = "http://localhost:8080/api/todos/:id"
    params = {"id": "$id"}
    query = {"id": "$id"}
    headers = {"Authorization": "$token", "Content-Type": "application/json"}
    body = {"extra": "something"}
}
