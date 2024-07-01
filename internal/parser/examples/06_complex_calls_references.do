let {
    var1 = env("VAR_VAL", "ANOTHER_VAR_VAL");
    var2 = env(var1 "default2");
    var3 = var2;
    var4 = var3;
}

do {
    method = "POST";
    url = "http://example.com";
    headers = {
        "Content-Type": "application/json"
    };
}
