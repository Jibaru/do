let {
    var1 = "param";
    the_file = file("/path/to/file");
}

// Multipart form data
do {
    method = "POST";
    url = "http://localhost:8080/upload";
    headers = {
        "Content-Type": "multipart/form-data"
    };
    body = {
        "key1": "value1",
        "key2": the_file,
        "key3": var1
    };
}
