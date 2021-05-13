#### QA

### GIT 
git checkout qa
git pull
git pull origin master
git add . 
git commit -m "merged with master (dev)"
git push

# heroku authorizations:create
# $env:HEROKU_API_KEY = "#####"; 
docker login --username=_ --password=$(heroku auth:token) registry.heroku.com
heroku git:remote -a planets-golang-api-qa
git remote -v

docker build --rm -f "Dockerfile.qa" -t planets-golang-api-qa:v1 "."
docker tag planets-golang-api-qa:v1 registry.heroku.com/planets-golang-api-qa/web
docker push registry.heroku.com/planets-golang-api-qa/web

heroku container:release web
heroku open 
heroku logs -t

