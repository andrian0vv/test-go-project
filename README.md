# test-go-project

## Run service
```bash
make up
```

## Stop service
```bash
make down
```

## Create user
```curl
curl -X POST --location "http://127.0.0.1:8080/api/users" \
    -H "Content-Type: application/json" \
    -d "{
          \"first_name\": \"ivan\",
          \"last_name\": \"ivanov\",
          \"age\": 18,
          \"is_married\": true,
          \"password\": \"12345678\"
        }"
```

## Create order
```curl
curl -X POST --location "http://127.0.0.1:8080/api/orders" \
    -H "Content-Type: application/json" \
    -d "{
          \"user_id\": 1,
          \"products\": [
            {
              \"id\": 1,
              \"quantity\": 1
            }
          ]
        }"
```