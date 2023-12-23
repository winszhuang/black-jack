import { EOperationCode } from '@/enums/msg-code'
import { useWs } from './use-ws'
import { ref } from 'vue'
import { EWsRoute } from '@/enums/ws-route'
import { EMessageType } from '@/enums/message'
import { Subject } from 'rxjs'
import { WSPlayGameReqResData } from '@/types/ws'

const genWsInstance = () => {
  const isDev = import.meta.env.DEV
  const websocketUrl = isDev
    ? 'ws://localhost:8080/ws'
    : (import.meta.env.VITE_API_URL as string).replace('https', 'wss')
  const token = localStorage.getItem('access_token') || 'no'
  return useWs<EWsRoute>(websocketUrl, ['access_token', token])
}

let ws: ReturnType<typeof useWs>
const playersDetail = ref<PlayerDetail[]>([])
const myId = ref('')
const messageList = ref<MessageItem[]>([])
const rooms = ref<Room[]>([])
let messageCounter = 0

// #NOTICE 這些放在use裡面會有問題
const loginSuccess$ = new Subject<{ token: string; user_id: string }>()
const registerSuccess$ = new Subject<void>()
const connectSuccess$ = new Subject<LoginData>()
const joinRoomSuccess$ = new Subject<string>()
const updatePlayerDetail$ = new Subject<PlayerDetail[]>()

const someOneJoin$ = new Subject<{ id: string; name: string }>()

export const useBlackJack = () => {
  ws ??= genWsInstance()

  ws.on(EWsRoute.WsConnected, (res) => {
    const loginData = res.data as LoginData
    if (loginData.is_login && loginData.user_id) {
      myId.value = loginData.user_id!
    }
    connectSuccess$.next(res.data)
  })

  ws.on(EWsRoute.Login, (res) => {
    pushNotify(res.message, res.success ? EMessageType.Success : EMessageType.Error)
    if (res.success) {
      myId.value = res.data
      loginSuccess$.next(res.data)
    }
  })

  ws.on(EWsRoute.Register, (res) => {
    pushNotify(res.message, res.success ? EMessageType.Success : EMessageType.Error)
    if (res.success) {
      registerSuccess$.next(res.data)
    }
  })
  ws.on(EWsRoute.JoinRoom, (res) => {
    pushNotify(res.message, res.success ? EMessageType.Success : EMessageType.Error)
    if (res.success) {
      joinRoomSuccess$.next(res.data)
    }
  })
  ws.on(EWsRoute.LeaveRoom, (res) => {
    console.log(res)
  })
  ws.on(EWsRoute.PlayBlackJack, (res) => {
    if (!res.success) return
    console.log(res)
    const gameData = res.data as WSPlayGameReqResData
    switch (gameData.opcode) {
      case EOperationCode.BroadcastJoin:
        someOneJoin$.next(gameData.gamedata)
        break
      case EOperationCode.UpdatePlayersDetail:
        updatePlayerDetail$.next(gameData.gamedata)
        break
    }
  })
  ws.on(EWsRoute.GetRoomsInfo, (res) => {
    pushNotify(res.message, res.success ? EMessageType.Success : EMessageType.Error)
    if (!res.success) return

    rooms.value = res.data
    console.log(rooms.value)
  })

  updatePlayerDetail$.subscribe((detail) => {
    playersDetail.value = detail
  })

  const pushNotify = (message: string, type: EMessageType) => {
    const messageItem: MessageItem = {
      text: message,
      id: `${performance.now()}-${messageCounter++}`,
      type
    }

    messageList.value.push(messageItem)

    setTimeout(() => {
      const index = messageList.value.findIndex((m) => m.id === messageItem.id)
      if (index !== -1) {
        messageList.value.splice(index, 1)
      }
    }, 3000)
  }

  const wsSend = (route: EWsRoute, data?: any) => {
    return ws?.send(route, data)
  }

  const clientSend = (opcode: EOperationCode) => {
    wsSend(EWsRoute.PlayBlackJack, {
      opcode,
      gametype: 1,
      gamedata: {}
    } as WSPlayGameReqResData)
  }

  const onReady = () => clientSend(EOperationCode.ClientReady)
  const onHit = () => clientSend(EOperationCode.ClientHit)
  const onStand = () => clientSend(EOperationCode.ClientStand)

  return {
    wsSend,

    messageList,
    playersDetail,
    myId,
    rooms,

    onLoginSuccess: loginSuccess$,
    onRegisterSuccess: registerSuccess$,
    onConnectSuccess: connectSuccess$,
    onJoinRoomSuccess: joinRoomSuccess$,

    onSomeOneJoinRoom: someOneJoin$,
    onUpdatePlayerDetail: updatePlayerDetail$,

    onReady,
    onHit,
    onStand
  }
}
