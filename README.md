## Rbac API

### Endpoint
`POST /rbac`
`POST /extend_rbac`

```bash
curl -X POST -H "Content-Type: application/json" -H "Authorization: " -d '{ 
  "username": "ftest", 
  "namespace": "test", "email": "fateme.pourghaemi@gmail.com" 
}' localhost:8080/rbac

curl -X POST -H "Content-Type: application/json" -H "Authorization: " -d '{ 
  "username": "ftest",
  "namespaces": ["test1","test2"], "email": "fateme.pourghaemi@gmail.com"
}' localhost:8080/extend_rbac
```
