# A lightweight docker container to listen for new emails(IMAP) and forward to telegram(RestAPI) chat on behalf of bot. 

### Example of usage:
```
docker build -t mail2telegram
```

```
docker run -d \
  --name=mail2telegram \
  -e TELEGRAM_CHAT_ID=<CHAT ID> \
  -e TELEGRAM_BOT_TOKEN=<TOKEN> \
  -e MAIL_IMAP_HOST=imap.gmail.com:993 \
  -e MAIL_USER=<somemail@gmail.com> \
  -e MAIL_PASSWORD=<password> \
  -e MAIL_POLL_INTERVAL_SEC=30 \ # Will be used in case IMAP server doesn't support IDLE
  --restart unless-stopped \
  mail2telegram
```

