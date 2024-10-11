# GoBlackjack
This contains a Go version of the Blackjack game as an API.

## Environment
The `.env` file must be created for this application to run.

### Env File Settings
Edit the `.env` file to set the following:

| Setting             | Description                                                       | Example   |
|---------------------|-------------------------------------------------------------------|-----------|
| SERVER_HOST         | The host name the service will listen on                          | localhost |
| SERVER_PORT         | The port the service will listen on                               | 3000      |
| SERVER_READ_TIMEOUT | The number of seconds before a read timeout will occur            | 10        |
| STAGE_STATUS        | If this is set to `dev` then it will not attempt graceful closure | dev       |

## Usage

This API contains 5 interactions.

1. [Deal](#deal)
2. [Hit](#hit)
3. [Stay](#stay)
4. [Stats](#stats)
5. [History](#history)

> **NOTE:** Depending on your shell, you may need to remove the quotes around the URLS in the CURL commands

### Deal<a id="deal"></a>
This endpoint takes no parameters, and will either start a new game if one does not exist for the device making the call,
or will retrieve any game currently in progress for the device.

The returned data contains the players cards and the token to play the game.

> **NOTE:** The device ID is a hash of the user agent and the client IP

```shell
curl -X 'POST' 'http://localhost:3000/deal'
```
```json
{
  "token":"190324df-34c2-4c07-97a5-1a06a21c9f6d",
  "device":"d08d4747b78e17c5459e8744604b90b35e669426f9c9d8e5b161b8828711c1ba",
  "cards":["J♦","A♥"],
  "dealerCards":[],
  "handValue":21,
  "dealerValue":0,
  "status":"playing"
}
```

### Hit<a id="hit"></a>
This endpoint optionally takes the game token as a parameter, and will draw another card for the players hand.
If the token is not specified then the device ID will be used instead to find the game.

The returned data contains the players cards and the token to play the game.

```shell
curl -X 'POST' 'http://localhost:3000/hit?token=game-token-goes-here'

curl -X 'POST' 'http://localhost:3000/hit'
```
```json
{
  "token":"2203e6c9-7383-48d9-9002-f441520a7791",
  "device":"d08d4747b78e17c5459e8744604b90b35e669426f9c9d8e5b161b8828711c1ba",
  "cards":["3♥","10♦","3♦"],
  "dealerCards":[],
  "handValue":16,
  "dealerValue":0,
  "status":"playing"
}
```

### Stay<a id="stay"></a>
This endpoint optionally takes the game token as a parameter, and will pass the turn to the dealer who will draw cards.
Both hands will be evaluated and a winner will be chosen.
If the token is not specified then the device ID will be used instead to find the game.

The returned data contains the players and the dealers cards, their relative values and the token to play the game.
However the game is over at this point and an new /deal call must be made to start a new game.

```shell
curl -X 'POST' 'http://localhost:3000/stay?token=game-token-goes-here'

curl -X 'POST' 'http://localhost:3000/stay'
```
```json
{
  "token":"2203e6c9-7383-48d9-9002-f441520a7791",
  "device":"d08d4747b78e17c5459e8744604b90b35e669426f9c9d8e5b161b8828711c1ba",
  "cards":["J♦","A♥"],
  "dealerCards":["6♣","J♣","6♠"],
  "handValue":21,
  "dealerValue":22,
  "status":"Dealer Bust"
}
```

### Stats<a id="stats"></a>
This endpoint takes no parameters and will return the win, loss and draw count for the device making the call.

```shell
curl 'http://localhost:3000/stats'
```
```json
{
  "wins":4,
  "loses":2,
  "draws":1
}
```

### History<a id="history"></a>
This endpoint takes no parameters and will return the game history for the device making the call, as an array of responses.

```shell
curl 'http://localhost:3000/history'
```

```json
[
  {"token":"6c359eb8-16bb-406a-93ff-6fbdaf1e5519",device:"d08d4747b78e17c5459e8744604b90b35e669426f9c9d8e5b161b8828711c1ba","cards":["6♣","5♦","10♥"],"dealerCards":["7♠","5♥","Q♣"],"handValue":21,"dealerValue":22,"status":"Dealer Bust"},
  {"token":"de3db63b-4363-4c33-80cc-3ff51f02ea81",device:"d08d4747b78e17c5459e8744604b90b35e669426f9c9d8e5b161b8828711c1ba","cards":["9♠","5♥","7♥"],"dealerCards":["Q♥","5♠","10♥"],"handValue":21,"dealerValue":25,"status":"Dealer Bust"},
  {"token":"420b767b-9506-47cc-a1e8-ed11d513fd30",device:"d08d4747b78e17c5459e8744604b90b35e669426f9c9d8e5b161b8828711c1ba","cards":["4♣","10♣","9♦"],"dealerCards":["6♣","3♥"],"handValue":23,"dealerValue":9,"status":"Bust"}
]
```
