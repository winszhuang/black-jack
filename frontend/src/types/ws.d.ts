import { EMsgCode } from '@/enums/msg-code'

type WsRequest = {
  data: any
  msg_code: EMsgCode
}

type WsResponse = {
  msg_code: EMsgCode
  data: any
  success: boolean
  error_code: string
  message: string
}

type WsEvents = {
  [EMsgCode.SomeOneJoin]: null
  [EMsgCode.SomeOneLeave]: { data: string }
  [EMsgCode.SomeOneReady]: { data: string }
  [EMsgCode.GameStart]: { data: string }
  [EMsgCode.SomeOneHit]: { data: string }
  [EMsgCode.SomeOneStand]: { data: string }
  [EMsgCode.GameOver]: { data: string }
  [EMsgCode.UpdatePlayersDetail]: { data: string }
}

type ResponseClientJoin = string
type ResponseBroadcastJoin = string
// type ResponseUpdatePlayersDetail = {

// }
