import { ERoute } from '@/enums/route'
import { RouteRecordRaw } from 'vue-router'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: ERoute.Entrance,
    component: () => import('@/pages/Entrance.vue')
  },
  {
    path: `/${ERoute.Lobby}`,
    name: ERoute.Lobby,
    component: () => import('@/pages/Lobby.vue')
  },
  {
    path: `/${ERoute.Game}/:roomId`,
    name: ERoute.Game,
    component: () => import('@/pages/Game.vue')
  },
  {
    path: '/:catchAll(.*)*',
    redirect: () => {
      alert('cant find this route!!')

      return { name: ERoute.Entrance }
    }
  }
]
