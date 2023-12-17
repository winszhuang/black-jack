import { EOperationCode } from '@/enums/msg-code'
import { EWsRoute } from '@/enums/ws-route'

type WsRequestLoginData = {
  username: string
  password: string
}
type WSRegisterReqData = {
  username: string
  password: string
}

type WSJoinRoomReqData = {
  room_id: string
}

type WSLeaveRoomReqData = {
  room_id: string
}

type WSPlayGameReqData = {
  opcode: EOperationCode
  gametype: number
  gamedata: any
}

type WsRequest = {
  data: any
  route: EWsRoute
}

type WsResponse = {
  route: EWsRoute
  data: any
  success: boolean
  error_code: string
  message: string
}

type WsEvents = {
  [EOperationCode.SomeOneJoin]: null
  [EOperationCode.SomeOneLeave]: { data: string }
  [EOperationCode.SomeOneReady]: { data: string }
  [EOperationCode.GameStart]: { data: string }
  [EOperationCode.SomeOneHit]: { data: string }
  [EOperationCode.SomeOneStand]: { data: string }
  [EOperationCode.GameOver]: { data: string }
  [EOperationCode.UpdatePlayersDetail]: { data: string }
}

type ResponseClientJoin = string
type ResponseBroadcastJoin = string
// type ResponseUpdatePlayersDetail = {

// }
