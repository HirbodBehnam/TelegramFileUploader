# Telegram File Uploader

Upload big files to Telegram via a bot

## Usage

At first obtain `APP_ID` and `APP_HASH` from [my telegram](https://my.telegram.org/). Then create a bot
using [bot father](https://t.me/BotFather) and set the environment variable `BOT_TOKEN` as your bot token. Also
set `RECEIVER_ID` as your username.

Before using this program make sure that you have started the bot from your Telegram account.

Also make sure that ffmpeg and ffprobe are installed and in your path.

### Command line arguments help:

```
Usage: uploader [-d|-p|-v|-a|-m] file_to_upload
Flags:                                                                        
        -d: Upload simple document                                            
        -p: Upload as picture                                                 
        -v: Upload as video                                                   
        -a: Upload as voice                                                   
        -m: Upload as music                                                   
Metadata of the file will be found out using ffmpeg and ffprobe               
                                                                              
Environment variables which must be set: APP_ID APP_HASH BOT_TOKEN RECEIVER_ID
```

You can also not define the flag to upload the file as document.

## Known Issues

* For now, the waveform of voices won't work. I might just create a random waveform to make it show something on big
  files.