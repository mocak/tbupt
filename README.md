# TBUPT (Toggle Backend Unattended Programming Test)

TBUPT is an implementation of a REST API to simulate a deck of cards

## Requirements

- Go ^1.17 or Docker

## Installation

### by go package

To install this package, simply go get it:

``
go get github.com/mocak/tbupt
``

### Makefile Commands

If you downloaded the source files, you can start and test app by commands below.

| Command      | Description                         |
|--------------|-------------------------------------|
| start        | Starts server                       |
| test         | Runs tests                          |
| docker-start | Crates and starts docker container  |
| docker-stop  | Stops and  removes docker container |
| docker-test  | Runs tests on docker container      |

## Usage

### Create Deck

#### Request:

URL:

``
POST localhost:3000?cards=AS,5D
``

Body:

```
{
    "shuffled": true 
}
```

Response:

```
{
    "DeckID": "1812b565-ec8f-44ff-b7bf-b266da50cbeb",
    "Shuffled": true,
    "Remaining": 2
}
```

### Draw Card

URL:

``
POST localhost:3000/<deck_id>/draw
``

Body:

```
{
    "count": 1 
}
```

Response:

```
[
    {
        "value": "5",
        "suit": "DIAMONDS",
        "code": "5D"
    }
]
```

### Open Deck

URL:

``
PUT localhost:3000/<deck_id>/draw
``

Body:

```
{
    "count": 1 
}
```

Response:

```
{
    "deck_id": "1812b565-ec8f-44ff-b7bf-b266da50cbeb",
    "shuffled": true,
    "remaining": 1,
    "cards": [
        {
            "value": "ACE",
            "suit": "SPADES",
            "code": "AS"
        }
    ]
}
```

## About this Solution

- Using memory as data storage, sql implementation can be done easily by implementing Storage interfaces.
- There is validation (and normalization) layer above storage to keep the storage dumb as possible. Multiple similar layers can be added easily by interface chaining if required.
- Project structure started as MVC and can be converted to other designs (domain driven, package oriented...) when scope started to become clearer.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[GNU](https://choosealicense.com/licenses/gpl-3.0/)