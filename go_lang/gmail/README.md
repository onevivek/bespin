# Nylas Lite Gmail

## Quickstart
- replace dummy access_token env variable in `docker-compose.yml`:
    - get client id and client secret from [1pass](https://start.1password.com/open/i?a=ZOXIFG7ZCNHUVNWQZIOLDFS45I&v=qkx6k5ffvmlr3d5fvaagows2em&i=z7f66jojsbe4jjrv5dbka7yulq&h=nylas.1password.com)
    - generate access token from https://developers.google.com/oauthplayground/
- create `/ssh` directory and download a ssh private key named `id_rsa` from [1pass](https://nylas.1password.com/vaults/r3hxsekmjk56ab33omvzg36s5i/allitems/3ils42fufoced5l3v6lwss3puy)   
- run `make compose-up` 

It should output:
```
...
nylas-lite-gmail    | {"level":"info","msg":"started sync","time":"2021-03-30T09:34:01.325583818"}
nylas-lite-gmail    | {"level":"info","msg":"fetched 300 messages","time":"2021-03-30T09:34:09.757291532"}
nylas-lite-gmail    | {"level":"info","msg":"stopped sync","time":"2021-03-30T09:34:09.757329147"}
nylas-lite-gmail exited with code 0
```

If you get account from redis it should have `last_run_at`, `provider_cursor_id` and `status` populated:

```
[tihomirjovicic@localhost gmail]$ make compose-exec 
docker container exec -it gmail_debug_1 sh
/ # redis-cli -h redis
redis:6379> hgetall account:1
1) "status"
2) "done"
3) "last_run_at"
4) "1617196643"
5) "provider_cursor_id"
6) "7314817"
redis:6379> 
```

----

## Environment Variables

| Name  | Required | Values | Default | Description |
| --- | --- |--- | --- | --- | 
| `GMAIL_ACCESS_TOKEN` | yes | string | - | Gmail Access Token | 
| `GMAIL_REDIS_ADDRESS` | yes | string | redis:6379 | Redis server address | 
| `GMAIL_ACCOUNT_ID` | yes | string | - | Account Id | 
| `GMAIL_HISTORIC_SYNC_DAYS` | no | int | 0 | How many days of history to pull. If zero, it will pull from account's historyId | 
| `GMAIL_BATCH_TIMEOUT` | no | int | 10 | Http timeout for batch get call | 
| `GMAIL_DISABLE_BATCHING` | no | bool | false | Make all message requests individually instead of batching | 

----

https://developers.google.com/gmail/api/guides/

https://developers.google.com/gmail/api/reference/rest/

https://developers.google.com/oauthplayground
