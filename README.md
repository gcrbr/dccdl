# dccdl
 Small golang program to download DCC/XDCC files

## Usage
```console
./dccdl -i <irc server ip> -p <port, default 6667> -b <bot name> -c <channel, without #> -f <file id, without #> -s (only add -s if you want to download the file)
```
![Example file](example.png)

You can find files from the many XDCC search engines on the web

## Building
Clone the repo and then do:
```console
go build .
```

## Disclaimer
This tool is provided for authorized use only. Any unauthorized use will not get any support