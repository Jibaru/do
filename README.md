<div align="center">
    <img src="screenshots/logo.png" width="180px"/>
    <p style="font-style: italic">HTTP Requests in plain text</p>
</div>

<br/>

Do is a Go program that parses custom `.do` files and executes HTTP requests based on the parsed content.

# Motivation

I created `do` because I needed a way to version HTTP requests and distribute them across teams in a simple, plaintext format. `.do` files provide a human-readable and easily shareable way to define HTTP requests, making it convenient for collaboration and version control.

## Installation

### With golang

Execute:

```
go install github.com/jibaru/do/cmd/do@latest
```

Also you can rename to dohttp.

Unix:

```
mv ~/go/bin/do ~/go/bin/dohttp
```

Windows (CMD):

```
move %USERPROFILE%\go\bin\do.exe %USERPROFILE%\go\bin\dohttp.exe
```

## Usage

Execute your `filename.do` files:

```
do filename.do
```

## Example

```do
let {
    base = "https://api.example.com";
    userId = 23;
    token = "the-http-token";
    page = 1;
    deleted = false;
}

do {
    method = "POST";
    url = "$base/users/:id";
    params = {"id": userId};
    query = {
        "p": page,
        "deleted": deleted
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
    url = "$base/users/:id";
    params = {"id": userId};
    query = {
        "p": page,
        "deleted": deleted
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

In the do section, variables defined in the let section can be referenced using their names.
You can insert variables in strings using the syntax `"$variableName"`.
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

## Output

The output of the `do` command will be the request + response in a json format.

```json
{
   "do_file": {
      "let": {
         "variables": {
           "token": "token-value"
         }
      },
      "do": {
         "method": "POST",
         "url": "https://www.fakepage.com/keys/:id",
         "params": {
            "id": 1
         },
         "query": {
            "limit": 1
         },
         "headers": {
           "Authorization": "Bearer token-value"
         },
         "body": "value"
      }
   },
   "response": {
      "status_code": 200,
      "body": "{\"key\": 123}",
      "headers": {
         "Content-Type": "application/json; charset=utf-8",
      }
   },
   "error": null
}
```

The `request` shows the parsed request from the .do file.
The `response` shows the response from the request if everything works well.
The `error` shows the error if parsing the .do file or executing the request fails. It is only a string.

If you want to use the response into another program, make sure validate error is null before trying to parse the response and request.

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
