# Inline RBAC for structs

## Usage

Define access with struct tags:

```go
type User struct {
  DisplayName string `access:"admin:*, owner: read update, *: read" json:"display_name"`
  Password    string `access:"admin:*, owner: update" json:"password"`
  ACL         *ACL   `access:"admin:*, owner: update, *: read" access_field_name:"acl"`
}
```

## Testing
    go get -t ./...
    go test