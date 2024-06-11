let {
    token = "token";
    id = 12;
    date = "2020-01-18"
}

do {
    method = "GET";
    url = "http://localhost:8080/api/todos/:id";
    params = {"id": "$id"};
    query = {
        "after": "$date"
    };
    headers = {
        "Authorization": "Bearer $token",
        "Content-Type": "application/json"
    };
    body = `{"extra": "something"}`;
}
