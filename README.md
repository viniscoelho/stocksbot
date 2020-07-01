# stocksbot
A telegram bot to report stocks and currency conversion rate values., written in Go.
  
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
The code represents the currency you want to convert from, `EUR`, to another one, `BRL`.

Don't forget to update the container timezone accordingly. The default is set to my location `America/Sao_Paulo`.
All you need is to update this line from the `Dockerfile`:
```bash
ENV TZ=America/Sao_Paulo
```

## How to use it 
- The simple way (if you have Golang already installed in your machine)
  - `go run main.go`
- The simple way with Docker (yeah, you don't have to install any dependencies)
  - `docker build .` then
  - `docker run --detach --name stocksbot <image-name>`