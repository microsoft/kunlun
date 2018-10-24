Converts any terraform templates in the directory `templates/` into managable Go source code.

```
go-bindata -ignore "terraform.tf*" -o templates.go -pkg=webserver templates/
```
