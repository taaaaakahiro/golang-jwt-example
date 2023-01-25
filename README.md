# golang-jwt-example


## Run
```
$ make run
```

## Request
```
$ curl -X GET localhost:9000/user/all -H "accept: application/json"
$ curl -X POST localhost:9000/user/login -H "accept: application/json" -d"{\"login_id\":\"user1\",\"password\":\"pass1\"}"
```

## Docs
 - https://pkg.go.dev/github.com/dgrijalva/jwt-go  
 - https://qiita.com/akubi0w1/items/dee1000699a3e9d9b2e3  
 - https://qiita.com/suzuki0430/items/f8308db9220c7bf19fa2  
 
## jwt
 - https://jwt.io/
 - https://zenn.dev/mfykmn/articles/eeaeb9a05130b8
 - jwt token
```sh
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VySWQxIiwiZXhwIjozMjUwNTc5NTg3NSwibmJmIjoxNjc0NjQ2ODc3LCJpYXQiOjE2NzQ2NDY4NzcsImp0aSI6IjAxR1FNQlBYNzdCOFQzMFBYUE5RWFg4Mk1XIn0.zSOH26U02bu2XnX_TWJfwPUmztTZr95lEbakpyEra90 #unlimited jwt token
```

