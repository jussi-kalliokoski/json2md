# json2md

[![CI status](https://github.com/jussi-kalliokoski/json2md/workflows/CI/badge.svg)](https://github.com/jussi-kalliokoski/json2md/actions)

A CLI tool for formatting single-dimensional (nested arrays/objects not supported at the moment) JSON arrays of objects as markdown tables.

For example:

```bash
cat << EOF | json2md
[
    {
        "null": null,
        "bool": true,
        "number": "12345.6789",
        "big number": 12345678901234567891234,
        "string": "hello world"
    }
]
EOF
```

Outputs the following table:

```md
| null   | bool | number     | big number              | string      |
|--------|------|------------|-------------------------|-------------|
| <null> | true | 12345.6789 | 12345678901234567891234 | hello world |
```

On macOS, if you want to swap the JSON stored in your clipboard with markdown, run the following:

```bash
pbpaste | json2md | pbcopy
```
