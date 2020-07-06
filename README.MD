# README

## Example Commands

New User
```
curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/api/auth/sign-up

curl -X POST -d '{"email" : "'$(date +%Y%m%d%H%M%S)'@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/api/auth/sign-up

curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' charkadog.herokuapp.com/auth/sign-up
```

```
curl -X GET -H 'Content-Type: application/json' localhost:8000/api/user/confirm/1havh6c0qc1uk334bzu0nwhcykgcrch1
```

Get Current User
```
curl -X GET -d '' -H 'Content-Type: application/json' -H "Authorization: Bearer ${access_token}" localhost:8000/api/users/current
```

Login
```
access_token=$(curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/api/auth/login 2> /dev/null | jq ".token.accessToken" | tr -d '"')

curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/api/auth/login | jq

access_token_java=$(curl -X POST -d '{"email" : "admin@admin.com", "password" : "password"}' -H 'Content-Type: application/json' localhost:8080/auth/login | jq '.token.access_token' | tr -d '"')

access_token=$(curl -X POST -d '{"email" : "dono@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' charkadog.herokuapp.com/auth/login 2> /dev/null | jq ".token.accessToken" | tr -d '"')


curl -X POST -d '{"email" : "20200520234451@dono.com", "password" : "secret"}' -H 'Content-Type: application/json' localhost:8000/api/auth/login 


```

Get Current User
```
curl -X GET -d '' -H 'Accept: application/json, text/plain, */*' -H "Authorization: Bearer ${access_token}" localhost:8000/api/users/current

curl -X GET -d '' -H 'Accept: application/json, text/plain, */*' -H "Authorization: Bearer ${access_token_java}" localhost:8080/users/current
```


Add contact
```
curl -X POST -d '{"name" : "bob", "phone" : "98908432"}' -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' localhost:8000/api/contacts/new
```

Get contacts
```
curl -X GET -d '{"email" : "dono@dono.com", "password" : "secret"}' -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' localhost:8000/api/me/contacts
```

Add card transaction
```
curl -X POST -d '{"dateTime":"2020-04-25T11:39:41.422Z","amount":{"value":10000,"scale":2},"currencyCode":"ZAR","reference":"simulation","merchantName":"The Coders Bakery","merchantCity":"Cape Town","merchantCountryCode":"ZA","merchantCountryName":"South Africa","merchantCategoryCode":"bakeries","merchantCategoryName":"Bakeries","email":"donohutcheon@gmail.com","firstName":"Donovan","lastName":"Hutcheon"}' -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' localhost:8000/api/card-transactions/new

curl -X POST -d '{"dateTime":"2020-04-25T11:39:41.422Z","amount":{"value":10000,"scale":2},"currencyCode":"ZAR","reference":"simulation","merchantName":"The Coders Bakery","merchantCity":"Cape Town","merchantCountryCode":"ZA","merchantCountryName":"South Africa","merchantCategoryCode":"bakeries","merchantCategoryName":"Bakeries","email":"donohutcheon@gmail.com","firstName":"Donovan","lastName":"Hutcheon"}
' -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' charkadog.herokuapp.com/card-transactions/new

```

Get card transactions
```
curl -X GET -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' localhost:8000/api/me/card-transactions
curl -X GET -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' 'localhost:8000/api/me/card-transactions?from=9&count=3&sortField=amount&sortDir=desc' | jq

curl -X GET -H "Authorization: Bearer ${access_token}" -H 'Content-Type: application/json' charkadog.herokuapp.com/me/card-transactions
```

#### Original blog post https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

##Postgres
```
docker run -it --network host --rm postgres bash
```

## MariaDB CLI
```
docker run -it --network host --rm mariadb mysql -hdonovanh -uroot -pcharka

docker run -it --network host --rm mariadb mysql -heporqep6b4b8ql12.chr7pe7iynqr.eu-west-1.rds.amazonaws.com -ueb3fsl6m7o2tiaxo -ph0du2iazqmit1p2h fclpddz9oxai1nec

??? docker run -it --network side-project_default --rm mariadb mysql -hdonovanh -uroot -pcharka
```

## Heroku Config Vars
Staging:
```
heroku config:set token_password=CheckThisIsTheJwtPassword --app charkadog
heroku config:set BIND_ADDRESS=0.0.0.0 --app charkadog
heroku config:set ENVIRONMENT=staging --app charkadog
heroku config:set URL="https://charkadog.herokuapp.com" --app charkadog
```

Production:
```
heroku config:set token_password=YoThisIsTheJwtPassword --app heenadog
heroku config:set BIND_ADDRESS=0.0.0.0 --app heenadog
heroku config:set ENVIRONMENT=prod --app heenadog
heroku config:set URL="https://heenadog.herokuapp.com" --app heenadog
```


## Heroku URL

https://charkadog.herokuapp.com
https://heenadog.herokuapp.com
