# The Mr Poll Bot Codebase

This the codebase for the Mr Poll Discord bot and integrated API that the website frontend communicates with.\
This is written in the [Go Programming Language](https://go.dev), you can download Go by clicking [here](https://go.dev/dl/). \
After downloading Go you will have to install the required packages for this Project:
```shell
go mod tidy
```

## Disclaimer
This is an unfinished project and is likely not the actual code that the Mr Poll bot is running at this time. The aim of this project is to replace Mr Poll's code with a system that is stable and reliable while having a low memory footprint, hence why Go was chosen.

## Packages
If you're interested in looking through the packages that this project uses:
- [Disgo - Discord API Wrapper & More](https://github.com/disgoorg/disgo)
- [Gorm - Database ORM](https://gorm.io)
- [Env - Environment Variable Manager](https://github.com/gofor-little/env)

## Setup

Before you can start the bot and API, you will need to configure your `.env` file.\
You can copy the `template.env` file and rename it to `.env`, adding the values to the variables.\
\
You can then deploy the slash commands for the bot to Discord by using the script:
```shell
go run deploy/deploy.go
```

## Operation

The bot will require a PostgreSQL database in order to operate, you can set one up or use the docker compose file to make one for you.\
Install docker on your machine, then run the command:
```shell
docker compose -f docker-compose.dev.yml up
```

Now that you have a database running, you can start the bot:
```shell
go build
./bot
```
This bot uses a http server to retrieve interaction requests from Discord, you will need to portforward so Discord can send you requests.
The easiest way to do this is to just use [ngrok](https://ngrok.com), with the command:
```shell
ngrok http 3001
```
You can copy the link it gives you and set it as the bot's interaction endpoint on the [dev portal](https://discord.dev)