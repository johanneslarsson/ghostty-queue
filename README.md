# ghostty-queue
As the scandinavian I am, I like standing in queues. But this queue is digital and I have no idea how long it is.
This discord bot will help me keep track of how many people are ahead of me in the queue and also give me some
indication on the progress.

## Bot must be invited to the server (discord calls them guilds)
For the bot to work it has to be invited to the server, by someone with manage server privilege
https://discord.com/oauth2/authorize?client_id=1267905755797786686&guild_id=1005603569187160125&scope=bot&permissions=0
the application is called `gqueue`. 
When its added it says it has the power to add commands, but I have no plans to do that.

It don't need any special permissions, all it needs is the intent `SERVER MEMBERS INTENT`, something I as app owner
can enable myself. 

## Usage
Plan is to run this daily to produce `list.md` and a copy in `archive` with todays date of this repository.
```
GUILD_ID=1005603569187160125 BOT_TOKEN=<secret> go run .
```

## Example
I`ve created an example on how it would look [example](example.md)
