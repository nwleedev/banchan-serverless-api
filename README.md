# Banchan Serverless API

The backend api of [banchanapp.com](https://banchanapp.com) built with Golang &amp; Serverless Framework.

## Libraries Used

- SQLBoiler for querying database
- Serverless Framework

## Features

- Query posts by page parameter
- Query post by id
- Query posts by tags
- Search posts by keywords

You can generate Go structs and methods by this command.

```
sqlboiler psql
```

Database schemas and columns should be defined first.
