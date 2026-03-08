# goKeyValueStore


## Input using MQTT:

Post on topic: "topic":"keyvalue/{prefix}"

```json
  {"data":{"key":"exampleKey","value":"exampleValue"}}
```

## Input using REST API:

### Post
http://{ipaddress:port}/key/{prefix}/{key}
"body" = value as a string

### Get
http://{ipaddress:port}/key/{prefix}/{key}
value is returned as a string