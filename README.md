# stocksbot
A telegram bot to report stocks and currency conversion rate values, written in Go.
  
## How to customise it
The code is configured to send messages in times defined by me, suiting my needs.
For this, I used `robfig/cron` library to create cron schedules. The configuration
is the same as a `crontab`. If you need help to create a `cron`, here is a [link](https://crontab.guru/#*_*_*_*_*) to help you.

Also, if you want to add another stock or replace the current one, just do the following:
```golang
applStock, err := quote.Get("APPL")
if err != nil {
    return nil, err
}
values = append(values, types.NewFinance(applStock))
```
In this example, I used Apple stock symbol to fetch information about their stocks.

For currencies, the code is almost the same:
```golang
eurBrl, err := forex.Get("EURBRL=X")
if err != nil {
    return nil, err
}
values = append(values, types.NewFinance(&eurBrl.Quote))
``` 
The code above will fetch information about how much does a `EUR` costs in `BRL`.

Don't forget to update the container timezone accordingly. The default is set to my location `America/Sao_Paulo`.
All you need is to update this line from the `Dockerfile`:
```bash
ENV TZ=America/Sao_Paulo
```

## How to use it 
### If you don't have a bot and/or group
Before executing the code, you need to set up a bot and have a telegram group.
You will need a group ID in which the bot will send messages to and a token for your bot.
To get the group ID, just add the [IDBot](https://telegram.me/getidsbot) to the group and send a `/getgroupid` message.
For the bot, just follow [this tutorial](https://core.telegram.org/bots) and at the end of the set up, a token will be generated for you.

Last step is to consume the group ID and the bot token. Just replace these lines on `main.go`:
```golang
telegramToken  = "<my-bot-token>"
telegramChatID = int64(<my-group-id>)
```

### How to run
- The simple way (if you have Golang already installed in your machine)
  - `go run main.go`
- The simple way with Docker (yeah, you don't have to install any dependencies)
  - `docker build -t stocksbot.` then
  - `docker run --detach stocksbot`

That's it. Feel free to provide feedback and/or report issues. 
