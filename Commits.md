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

## Fifth Commit

1. Created a request struct for the initial registration route.
2. Made it so we use omitempty in every model's bson tag so that if a field(s) is empty, our db won't change the field(s) to empty/null.
3. Began work on the initiate registration route.

## Sixth Commit

1. Refactored the function that connects to the database.
2. Created a success response and a validate email helper function.
3. Complete implemented the initiate registration route.

## Seventh Commit

1. Installed SendGrid.
2. Created a function to send a registration email.

## Eight Commit

1. Integrated Twilio and created a SendRegistrationText function.

## Ninth Commit -IMPORTANT

1. Created a js function to generate 256 bit cryptographic strings to be used as auth token secrets.
2. Installed jwt and created a function to create an access and refresh token given a user id.
3. Created a function to generate the 6 digit user verification codes rather than using rand.
4. Integrated most, but not all, of the finalize registration function.

## Tenth Commit

1. Completely integrated the finalize registration route.
2. Installed bcrypt and created two functions: one to hash password and one to compare.

## Eleventh Commit

1. Implemented login and token login routes.

## Twelveth Commit

1. Implemented entire password reset process.
2. Added functions to send reset email and text.

## Thirteenth Commit

1. Created a profile model.
2. Created logic to create profile when user registers.
3. Made it so a user's profile is returned when they login.

## Fourteenth Commit

1. Changed error messages in login to add a bit more security.
2. Swapped around the position of some logic in intial registration route.
3. Made the Code request body field a string which is converted to an int to avoid logic flaw if the code provided is actually 0 and not just empty.

## Heroku

1. Started Heroku integration.
2. To connect project to existing app: heroku login, heroku git:remote -a appname

## Heroku

1. Everything works
