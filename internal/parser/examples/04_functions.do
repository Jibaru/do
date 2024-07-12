let {
    var1 = env("NO_EXIST", "default");
    path = env("ALREADY_EXISTS", "default2");
    var3 = date("ISO8601");
    var4 = uuid();
}

do {
    method = env("ENV_METHOD", "GET");
    url = "https://jsonplaceholder.typicode.com/todos/:id";
    query = {"id": var1, "id2": "$var1"};
}
