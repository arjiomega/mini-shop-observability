# mini-shop-observability

```
go tool oapi-codegen -config config.yaml openapi.yaml
go tool oapi-codegen
```
```
goose -dir migrations create create_products_table sql
goose -dir migrations sqlite3 products.db up
```

```HTTP
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Coffee"}'
>>> {"id":1,"name":"Coffee"}

curl http://localhost:8080/products/1
>>> {"id":1,"name":"Coffee"}
```

docker build -f product-service/Dockerfile -t mini-shop-product-service .

protoc \
  -I proto \
  --go_out=proto \
  --go_opt=module=mini-shop/proto \
  --go-grpc_out=proto \
  --go-grpc_opt=module=mini-shop/proto \
  proto/product/product.proto
