import { EOperationCode } from '@/enums/msg-code'
import { useWs } from './use-ws'
import { ref, watch } from 'vue'
import { EWsRoute } from '@/enums/ws-route'
import { EMessageType } from '@/enums/message'
import { Subject } from 'rxjs'

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
const loginSuccess$ = new Subject<string>()
const connectSuccess$ = new Subject<{ is_login: boolean }>()
const joinRoomSuccess$ = new Subject<string>()

export function useBlackJack() {
  ws ??= genWsInstance()

  ws.on(EWsRoute.WsConnected, (res) => {
    connectSuccess$.next(res.data)
  })

  ws.on(EWsRoute.Login, (res) => {
    pushNotify(res.message, res.success ? EMessageType.Success : EMessageType.Error)
    if (res.success) {
      loginSuccess$.next(res.data)
    }
  })

  ws.on(EWsRoute.Register, (res) => {
    console.log(res)
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
  })
  ws.on(EWsRoute.GetRoomsInfo, (res) => {
    pushNotify(res.message, res.success ? EMessageType.Success : EMessageType.Error)
    if (!res.success) return

    rooms.value = res.data
    console.log(rooms.value)
  })

  // ws.on(EOperationCode.ClientJoin, (res) => {
  //   pushNotify(res.message)
  //   myId.value = res.data as string
  // })

  // ws.on(EOperationCode.BroadcastJoin, (res) => {
  //   pushNotify(res.message)
  //   console.log('res.message', res.message)
  //   console.log('res.data', res.data)
  // })

  // ws.on(EOperationCode.BroadcastLeave, (res) => {
  //   pushNotify(res.message)
  // })

  // ws.on(EOperationCode.BroadcastGameStart, (res) => {
  //   pushNotify(res.message)
  // })

  // ws.on(EOperationCode.BroadcastGameOver, (res) => {
  //   pushNotify(res.message)
  // })

  // ws.on(EOperationCode.BroadcastReStart, (res) => {
  //   pushNotify(res.message)
  // })

  // ws.on(EOperationCode.UpdatePlayersDetail, (res) => {
  //   console.log(res.data)
  //   playersDetail.value = res.data as PlayerDetail[]
  // })

  function pushNotify(message: string, type: EMessageType) {
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

  function wsSend(route: EWsRoute, data?: any) {
    return ws?.send(route, data)
  }

  //   const onReady = () => wsSend(EOperationCode.ClientReady)
  //   const onHit = () => wsSend(EOperationCode.ClientHit)
  //   const onStand = () => wsSend(EOperationCode.ClientStand)

  return {
    wsSend,

    messageList,
    playersDetail,
    myId,
    rooms,

    onLoginSuccess: loginSuccess$,
    onConnectSuccess: connectSuccess$,
    onJoinRoomSuccess: joinRoomSuccess$

    // onReady,
    // onHit,
    // onStand
  }
}
