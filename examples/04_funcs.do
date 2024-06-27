let {
    var1 = env("NO_EXIST", "default");
    path = env("PATH", "/path");
}

do {
    method = env("ENV_METHOD", "GET");
    url = "https://jsonplaceholder.typicode.com/todos/:id";
    query = {"id": 1};
}
