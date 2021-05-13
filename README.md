# API 
## This is an example of Rest API (CRUD) using GoLang, Gin, MongoDB
#
# Before Running
## It is needed to create a config.yml with the same format as config-example.yml file. 
## MONGO_HOST value is Mongo's Query String
## To run a specific enviroment it is needed to pass the env (dev, qa, prod) as parameter.
## Example: go run main.go -env=prod
#
# Deploy
## The aplication is already deployed 
## QA: https://planets-golang-api-qa.herokuapp.com/swagger/index.html
## PROD:  https://planets-golang-api.herokuapp.com/swagger/index.html
## The aplication has two deploy PowerShell scripts, QA and PROD.
## QA: deployQAHeroku.ps1
## PROD: deployPRODHeroku.ps1
## To deploy it is needed to configure HEROKU_API_KEY on yours machine env variables.
## The scripts update the environment git branch automatically.
#
# Documentation
## A Swagger can be found and used to test the API with the path: {{host}}/swagger/index.html
#
# Tests
## Inside the folder 'tests', there is a Postman collection and 3 enviroments (dev, qa, prod). It can be used to test the API using Postman. Running all the tests in sequence should not fail any test.
##