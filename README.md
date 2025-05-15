# plex-pusher

## Help
```
$ ./plex-pusher --help
NAME:
   plex-pusher - plex event pusher

USAGE:
   plex-pusher [global options] command [command options]

VERSION:
   0.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen value  Specify the host:port to run pusher (default: "0.0.0.0:9876") [%PLEX_PUSHER_LISTEN%]
   --token value   WxPusher app token [%PLEX_PUSHER_TOKEN%]
   --uid value     WxPusher app uid [%PLEX_PUSHER_UID%]
   --accesslog     Enable accesslog (default: false) [%PLEX_PUSHER_ACCESSLOG%]
   --help, -h      show help
   --version, -v   print the version
```