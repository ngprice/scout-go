# Scout-Go

This is a lightweight server for the game "Scout", by Kei Kajino, written in Go.

The full rules for "Scout" can be found here: https://www.hras.cz/user/related_files/new_edition_scout_rules_eng.pdf

This repo is part of the scout-rl training project; the model training repo is here: https://github.com/ngprice/scout-rl


## Setup

Setup the local environment with the Makefile:

* `make deps`
    * Installs dependencies for Go runtime
* `make build`
    * Builds Go binary
* `make test`
    * Runs unit tests
* `make run`
    * Launches a local game server listening on :50051


## API

The game server communicates via GRPC/protobuf. The .proto file defining the service is here: /proto/scout.proto

### Service Endpoints

`ScoutService` exposes this interface:
```
service ScoutService {
  rpc CreateGame      (CreateGameRequest)      returns (CreateGameResponse);
  rpc PlayerAction    (PlayerActionRequest)    returns (PlayerActionResponse);
  rpc GetGameState    (GetGameStateRequest)    returns (GetGameStateResponse);
  rpc GetPlayerState  (GetPlayerStateRequest)  returns (GetPlayerStateResponse);
  rpc GetValidActions (GetValidActionsRequest) returns (GetValidActionsResponse);
}
```

## Protocol Documentation
<a name="top"></a>

### Table of Contents

- [proto/scout.proto](#proto_scout-proto)
    - [Action](#scout-Action)
    - [Card](#scout-Card)
    - [CreateGameRequest](#scout-CreateGameRequest)
    - [CreateGameResponse](#scout-CreateGameResponse)
    - [Game](#scout-Game)
    - [GetGameStateRequest](#scout-GetGameStateRequest)
    - [GetGameStateResponse](#scout-GetGameStateResponse)
    - [GetPlayerStateRequest](#scout-GetPlayerStateRequest)
    - [GetPlayerStateResponse](#scout-GetPlayerStateResponse)
    - [GetValidActionsRequest](#scout-GetValidActionsRequest)
    - [GetValidActionsResponse](#scout-GetValidActionsResponse)
    - [Player](#scout-Player)
    - [PlayerActionRequest](#scout-PlayerActionRequest)
    - [PlayerActionResponse](#scout-PlayerActionResponse)
    - [PlayerState](#scout-PlayerState)
  
    - [Action.ActionType](#scout-Action-ActionType)
  
    - [ScoutService](#scout-ScoutService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_scout-proto"></a>
<p align="right"><a href="#top">Top</a></p>

### proto/scout.proto



<a name="scout-Action"></a>

#### Action



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| action_type | [Action.ActionType](#scout-Action-ActionType) |  |  |
| scout_take_index | [int32](#int32) |  |  |
| scout_put_index | [int32](#int32) |  |  |
| show_first_index | [int32](#int32) |  |  |
| show_length | [int32](#int32) |  |  |






<a name="scout-Card"></a>

#### Card



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value1 | [int32](#int32) |  |  |
| value2 | [int32](#int32) |  |  |






<a name="scout-CreateGameRequest"></a>

#### CreateGameRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| num_players | [int32](#int32) |  |  |






<a name="scout-CreateGameResponse"></a>

#### CreateGameResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |






<a name="scout-Game"></a>

#### Game



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| num_players | [int32](#int32) |  |  |
| active_player_index | [int32](#int32) |  |  |
| active_set | [Card](#scout-Card) | repeated |  |
| active_set_player_index | [int32](#int32) |  |  |
| consecutive_scouts | [int32](#int32) |  |  |
| round | [int32](#int32) |  |  |
| complete | [bool](#bool) |  |  |
| player_states | [PlayerState](#scout-PlayerState) | repeated |  |






<a name="scout-GetGameStateRequest"></a>

#### GetGameStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |






<a name="scout-GetGameStateResponse"></a>

#### GetGameStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game | [Game](#scout-Game) |  |  |






<a name="scout-GetPlayerStateRequest"></a>

#### GetPlayerStateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| player_index | [int32](#int32) |  |  |






<a name="scout-GetPlayerStateResponse"></a>

#### GetPlayerStateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player | [Player](#scout-Player) |  |  |






<a name="scout-GetValidActionsRequest"></a>

#### GetValidActionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| player_index | [int32](#int32) |  |  |






<a name="scout-GetValidActionsResponse"></a>

#### GetValidActionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mask | [bool](#bool) | repeated |  |






<a name="scout-Player"></a>

#### Player



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| index | [int32](#int32) |  |  |
| score | [int32](#int32) |  |  |
| hand | [Card](#scout-Card) | repeated |  |
| can_reverse_hand | [bool](#bool) |  |  |
| can_scout_and_show | [bool](#bool) |  |  |






<a name="scout-PlayerActionRequest"></a>

#### PlayerActionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| game_id | [string](#string) |  |  |
| player_index | [int32](#int32) |  |  |
| action | [Action](#scout-Action) |  |  |






<a name="scout-PlayerActionResponse"></a>

#### PlayerActionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| err | [bool](#bool) |  |  |
| errMsg | [string](#string) |  |  |






<a name="scout-PlayerState"></a>

#### PlayerState



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| player_index | [int32](#int32) |  |  |
| hand_size | [int32](#int32) |  |  |
| score | [int32](#int32) |  |  |





 


<a name="scout-Action-ActionType"></a>

#### Action.ActionType


| Name | Number | Description |
| ---- | ------ | ----------- |
| ActionScout | 0 |  |
| ActionScoutReverse | 1 |  |
| ActionShow | 2 |  |
| ActionScoutAndShow | 3 |  |
| ActionScoutAndShowReverse | 4 |  |
| ActionReverseHand | 5 |  |


 

 


<a name="scout-ScoutService"></a>

#### ScoutService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateGame | [CreateGameRequest](#scout-CreateGameRequest) | [CreateGameResponse](#scout-CreateGameResponse) |  |
| PlayerAction | [PlayerActionRequest](#scout-PlayerActionRequest) | [PlayerActionResponse](#scout-PlayerActionResponse) |  |
| GetGameState | [GetGameStateRequest](#scout-GetGameStateRequest) | [GetGameStateResponse](#scout-GetGameStateResponse) |  |
| GetPlayerState | [GetPlayerStateRequest](#scout-GetPlayerStateRequest) | [GetPlayerStateResponse](#scout-GetPlayerStateResponse) |  |
| GetValidActions | [GetValidActionsRequest](#scout-GetValidActionsRequest) | [GetValidActionsResponse](#scout-GetValidActionsResponse) |  |

 



### Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

