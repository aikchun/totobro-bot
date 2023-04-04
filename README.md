# totobro-bot
A no-nonsense helper to fetch Singapore lottery results and details of the next draw.

# Getting started
You will need:
1. Create a bot or get a bot token through the [BotFather](https://t.me/botfather)
2. Docker
3. ngrok
4. Postman or other REST API client

## Get dependencies
```
go get
```

## Change your env
Copy `.dbenv_example` into `.dbenv`. Change the credentials as necessary.

Copy `.env_example` into `.env`. 

Put the token you have gotten in Getting started (1) as `BOT_TOKEN`

`DATABASE_URL` should have the credentials listed in your `.dbenv`.

`FETCH_NEXT_DRAW_TASK_MOCK_ID` can be any int64 value.

`FETCH_NEXT_DRAW_TASK_USER_FIRST_NAME` can be any string value.

### Database
Run `docker-compose up -d`
This will start create a postgres db using docker.

#### Migrations
Run:

```bash
docker run -v /absolute/path/to/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://{POSTGRES_HOST}:{POSTGRES_PASSWORD}@localhost:5432/{POSTGRES_DB}\?sslmode=disable up
```

### Build and serve your app
Run:
```
go build main.go bot.go
./main
```
This will make the app serve on `8080`

#### Run ngrok
```
ngrok http 8080
```

This would output several URL but we are only interested in this.
```
...
Web Interface                 http://127.0.0.1:4040
Forwarding                    https://xxxxx-xxxx-xxx.ap.ngrok.io -> http://localhost:8080
```

Copy the `https` link.
Open up your web browser or your favorite REST API client key in this and make this request.

```
GET https://api.telegram.org/bot{token}/setWebhook?url={https://xxxxx-xxxx-xxx.ap.ngrok.io}
```

You should see: `{"ok":true,"result":true,"description":"Webhook was set"}`

And you are good to go.

To start off, go talk to your bot and send `/start` to it.
