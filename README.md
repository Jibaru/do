<div style="align: center;">
<img src="screenshots/logo.png" width="180px"/>
</div>

<div style="align: center;">
<i>HTTP Requests in plain text</i>
</div>

<br/>

Do is a Go program that parses custom `.do` files and executes HTTP requests based on the parsed content.

# Motivation

I created `do` because I needed a way to version HTTP requests and distribute them across teams in a simple, plaintext format. `.do` files provide a human-readable and easily shareable way to define HTTP requests, making it convenient for collaboration and version control.

## Usage

Execute your `filename.do` files:

```
do filename.do
```

## Example

```do
let {
    userId = 23;
    token = "the-http-token";
    deleted = false;
    // Add more variables here...
}

do {
    method = "POST";
    url = "https://api.example.com/users/:id";
    params = {"id": "$userId"};
    query = {
        "actives": true,
        "deleted": false
    };
    headers = {
        "Authorization": "Bearer $token",
        "Content-Type": "application/json"
    };
    body = `{
        "name": "John Doe",
        "email": "john@example.com"
    }`;
}
```

## File Format

The .do file format consists of two main sections: let and do.

### Let Section

The let section is used to declare variables that can be referenced in the do section. Variables are declared in the following format:

```do
let {
    userId = 23;
    token = "the-http-token";
    deleted = false;
    // Add more variables here...
}
```

You should use the `variable = value;` format to define a variable with its value. Values has specific types that are defined below:

### Types

| Type   | Description                                                                   | Example   |
| ------ | ----------------------------------------------------------------------------- | --------- |
| int    | Integer value                                                                 | 12        |
| float  | Decimal value                                                                 | 92.3      |
| bool   | Boolean value (true or false)                                                 | true      |
| string | A character sequence. You can use **\`** to wrap a string that contains **"** | "example" |

There are another values that only should be accepted in do section:

| Type | Description                     | Example                                      |
| ---- | ------------------------------- | -------------------------------------------- |
| map  | A collection of key-value pairs | {"key1": 12, "key2": false, "key3": "hello"} |

### Do Section

The do section specifies the HTTP request to be executed. It contains various fields such as method, URL, params, query, headers, and body. Here's an example of how to define the do section:

```do
do {
    method = "POST";
    url = "https://api.example.com/users/:id";
    params = {"id": "$userId"};
    query = {
        "actives": true,
        "deleted": false
    };
    headers = {
        "Authorization": "Bearer $token",
        "Content-Type": "application/json"
    };
    body = `{
        "name": "John Doe",
        "email": "john@example.com"
    }`;
}
```

In the do section, variables defined in the let section can be referenced using the syntax `"$variableName"`.
Take in note that `body` is not a map, it is a string.

You should use the `variable = value;` format to define a variable with its value. Values has specific types that are defined below:

### Fields

| Field   | Type   | Description                                                                             | Required | Example                               |
| ------- | ------ | --------------------------------------------------------------------------------------- | -------- | ------------------------------------- |
| method  | string | The http request method                                                                 | Yes      | "POST"                                |
| url     | string | The url to request. It should accept replace params that should be used with `:` + name | Yes      | "https://example.com"                 |
| params  | map    | The params to replace in the url (without `:`).                                         | No       | {"id": 26}                            |
| query   | map    | The query params for the request.                                                       | No       | {"active": false, "order": "asc"}     |
| headers | map    | The headers for the request.                                                            | No       | {"Authorization": "application/json"} |
| body    | string | The body for the request.                                                               | No       | \`{"name": "john"}\`                  |

## VS-Code do language support

You can add support for `.do` files using the following extension:

[Download Do Language Support Visual Studio Code Extension](https://marketplace.visualstudio.com/items?itemName=jibaru.do-language-support)

## Roadmap

- [x] Add support for variables
- [ ] Add support for load env variables
- [ ] Add support for prompt variables
- [ ] Add support for displaying beauty response
- [ ] Add support for comments
- [ ] Add support for importing files
