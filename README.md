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

## Config Properties og Default Values


| **Property**               | **Default Value**       |
|----------------------------|-------------------------|
| `mqtt_address`             | `localhost:1883`        |
| `mqtt_username`            | `mqtt-user`             |
| `mqtt_password`            | `mqtt-password`         |
| `gokeyvaluestore_url`      | `http://localhost:9091` |
| `own_rest_address`         | `:9092`                 |

## REST API:
Welcome to the Runner REST Service

- `/health/` - Health check endpoint
- **GET** `/training/{uid}` - Retrieve a training by UID
- **POST** `/training` - Create a new training
- **PUT** `/training/{uid}` - Update an existing training
- **DELETE** `/training/{uid}` - Delete a training
- **GET** `/training` - Retrieve all trainings
