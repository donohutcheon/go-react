# README

## Example Commands

New User
```
curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/user/new
```

Get Current User
```
curl -X GET -d '' -H 'Content-Type: application/json' -H "Authorization: Bearer ${access_token}" localhost:8000/users/current
```

Login
```
access_token=$(curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/auth/login 2> /dev/null | jq ".token.access_token" | tr -d '"')

access_token_java=$(curl -X POST -d '{"email" : "admin@admin.com", "password" : "password"}' -H 'Content-Type: application/json' localhost:8080/auth/login | jq '.token.access_token' | tr -d '"')
```

Get Current User
```
curl -X GET -d '' -H 'Accept: application/json, text/plain, */*' -H "Authorization: Bearer ${access_token}" localhost:8000/users/current

curl -X GET -d '' -H 'Accept: application/json, text/plain, */*' -H "Authorization: Bearer ${access_token_java}" localhost:8080/users/current
```


Add contact
```
curl -X POST -d '{"name" : "bob", "phone" : "98908432"}' -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' localhost:8000/contacts/new
```

Get contacts
```
curl -X GET -d '{"email" : "dono@dono.com", "password" : "secret"}' -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' localhost:8000/me/contacts
```

#### Original blog post https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

MariaDB CLI
```
docker run -it --network host --rm mariadb mysql -hdonovanh -uroot -pcharka

docker run -it --network host --rm mariadb mysql -heporqep6b4b8ql12.chr7pe7iynqr.eu-west-1.rds.amazonaws.com -ueb3fsl6m7o2tiaxo -ph0du2iazqmit1p2h fclpddz9oxai1nec'

??? docker run -it --network side-project_default --rm mariadb mysql -hdonovanh -uroot -pcharka
```

Heroku Config Vars
```
heroku config:set token_password=thisIsTheJwtPasswordMofo --app charkadog
heroku config:set BIND_ADDRESS=0.0.0.0 --app charkadog
```