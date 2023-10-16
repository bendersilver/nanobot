```js
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: chat not found"
}
```

```js
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: group chat was migrated to a supergroup chat",
  "parameters": {
    "migrate_to_chat_id": -123456789
  }
}
```

```js
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: invalid file id"
}
```

```js
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: message is not modified"
}
```

```js
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: message text is empty"
}
```

```js
{
  "ok": false,
  "error_code": 400,
  "description": "[Error]: Bad Request: user not found"
}
```

```js
{
  "ok": false,
  "error_code": 400,
  "description": "Bad Request: wrong parameter action in request"
}
```

```js
{
  "ok": false,
  "error_code": 409,
  "description": "Conflict: terminated by other long poll or webhook"
}
```

```js
{
  "ok": false,
  "error_code": 403,
  "description": "Forbidden: bot was blocked by the user"
}
```

```js
{
  "ok": false,
  "error_code": 403,
  "description": "Forbidden: bot can't send messages to bots"
}
```

```js
{
  "ok": false,
  "error_code": 403,
  "description": "Forbidden: bot was kicked from the group chat"
}
```

```js
{
  "ok": false,
  "error_code": 403,
  "description": "Forbidden: user is deactivated"
}
```

```js
{
  "ok": false,
  "error_code": 429,
  "description": "Too Many Requests: retry after X",
  "parameters": { "retry_after": 123 }
}
```

```js
{
  "ok": false,
  "error_code": 401,
  "description": "Unauthorized"
}
```

```js
{
  "ok": false,
  "error_code": 409,
  "description": "Conflict: can't use getUpdates method while webhook is active; use deleteWebhook to delete the webhook first"
}
```
