#### QA

### GIT 
git checkout prod
git pull
git pull origin qa
git add . 
git commit -m "merged with qa"
git push

# heroku authorizations:create
# $env:HEROKU_API_KEY = "#####"; 
docker login --username=_ --password=$(heroku auth:token) registry.heroku.com
heroku git:remote -a planets-golang-api
git remote -v

docker build --rm -f "Dockerfile.prod" -t planets-golang-api:v1 "."
docker tag planets-golang-api:v1 registry.heroku.com/planets-golang-api/web
docker push registry.heroku.com/planets-golang-api/web

heroku container:release web
heroku open 
heroku logs -t

