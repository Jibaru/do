let {
    var1 = 1;
    var2 = "hello";
    var3 = true;
    var4 = false;
    var5 = var1;
}

do {
    method = "GET";
    url = "http://example.com/:id";
    params = {
        "id": "$var1"
    };
    headers = {
        "Content-Type": "application/json",
        "X-Message": "$var2",
        "X-Var5": var5
    };
    body = `{
        "var1": $var1,
        "var2": "$var2",
        "var3": $var3,
        "var4": $var4,
        "var5": $var5
    }`;
}
