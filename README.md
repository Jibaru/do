# Do

Do is a Go program that parses custom `.do` files and executes HTTP requests based on the parsed content.

# Motivation

I created `do` because I needed a way to version HTTP requests and distribute them across teams in a simple, plaintext format. `.do` files provide a human-readable and easily shareable way to define HTTP requests, making it convenient for collaboration and version control.

## Usage

Execute your filename.do files:

```
do filename.do
```

## File Format

The .do file format consists of two main sections: let and do.

### Let Section

The let section is used to declare variables that can be referenced in the do section. Variables are declared in the following format:

```do
let {
    variableName = "value";
    // Add more variables here...
}
```

### Do Section

The do section specifies the HTTP request to be executed. It contains various fields such as method, URL, params, query, headers, and body. Here's an example of how to define the do section:

```do
do {
    method = "POST";
    url = "https://api.example.com/users/:id";
    params = {"id": "$userId"};
    query = {"token": "$token"};
    headers = {"Authorization": "Bearer $token", "Content-Type": "application/json"};
    body = `{"name": "John Doe", "email": "john@example.com"}`
}
```

In the do section, variables defined in the let section can be referenced using the syntax `"$variableName"`.
