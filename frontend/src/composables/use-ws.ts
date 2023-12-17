/* eslint-disable @typescript-eslint/no-explicit-any */

import { EWsRoute } from '@/enums/ws-route'
import { WsResponse } from '@/types/ws'

// useWs.ts
const eventList: Array<{ route: EWsRoute; callback: (data: WsResponse) => void }> = []
let ws: WebSocket

export function useWs<E extends EWsRoute>(url: string | URL, protocols?: string | string[]) {
  ws ??= new WebSocket(url, protocols)

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
      if (ev.route === data.route) {
        return ev.callback(data)
      }
    }
  }

  function on(route: E, callback: (data: WsResponse) => void) {
    eventList.push({
      route: route as EWsRoute,
      callback
    })
  }

  function send(route: E, data?: any) {
    ws.send(
      JSON.stringify({
        route,
        data
      })
    )
  }

  return {
    send,
    on,
    originWs: ws
  }
}
