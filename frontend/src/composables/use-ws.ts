/* eslint-disable @typescript-eslint/no-explicit-any */

import { WsRequest, WsResponse } from '@/types/ws'

// useWs.ts
const eventList: Array<{ msg_code: number; callback: (data: WsResponse) => void }> = []
let ws: WebSocket

export function useWs<E extends number>(
  config = {
    dev: (import.meta.env.VITE_API_URL as string).replace('http', 'ws'),
    prod: (import.meta.env.VITE_API_URL as string).replace('https', 'wss')
  }
) {
  const isDev = import.meta.env.DEV
  ws ??= new WebSocket(isDev ? config.dev : config.prod)

  ws.onmessage = async (e) => {
    const jsonStr = e.data as string
    console.warn('receive data: ', jsonStr)

    let data: WsResponse
    try {
      data = JSON.parse(jsonStr) as WsResponse
    } catch (error) {
      throw Error(error as any)
    }

    for (const ev of eventList) {
      if (ev.msg_code === data.msg_code) {
        return ev.callback(data)
      }
    }
  }

  function on(msgCode: E, callback: (data: WsResponse) => void) {
    eventList.push({
      msg_code: msgCode as number,
      callback
    })
  }

  function onOpen(callback = () => console.log('ws open connection')) {
    ws.onopen = () => callback()
  }

  function onClose(callback = () => console.log('ws close connection')) {
    ws.onclose = () => callback()
  }

  function send(msgCode: E, data?: WsRequest) {
    ws.send(
      JSON.stringify({
        msgCode,
        data
      })
    )
  }

  return {
    send,
    on,
    onOpen,
    onClose,
    ws
  }
}
