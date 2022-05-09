## First Commit

1. Ran: "go get -u github.com/gofiber/fiber/v2 go.mongodb.org/mongo-driver/mongo github.com/joho/godotenv" and "go get -u github.com/gofiber/helmet/v2"
2. Setup a basic server.

## Second Commit

1. Configured the database connection.
2. Configured env variables.
3. Configured a root router.

## Third Commit

1. Created a User model.
2. Created nested api endpoints

## Fourth Commit -IMPORTANT

1. Updated the user model.
2. Created a temp object model.
3. Created a TTL index in database to delete temp objects after 300 seconds.
4. To create TTL index: go to db atlas, go to index tab of collection you want the index in, click create index, in options type {expireAfterSeconds: <time in seconds>} and in fields type { "fieldName": "1" }
