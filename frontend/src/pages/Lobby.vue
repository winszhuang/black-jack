<script setup lang="ts">
import { useBlackJack } from '@/composables/use-black-jack'
import { ERoute } from '@/enums/route'
import { EWsRoute } from '@/enums/ws-route'
import { router } from '@/router'

const { rooms, onConnectSuccess, onJoinRoomSuccess, wsSend } = useBlackJack()

onConnectSuccess.subscribe((info) => {
  if (!info.is_login) {
    router.push({ name: ERoute.Entrance })
  }
})

// #TODO 這邊為何不會過來?????
onJoinRoomSuccess.subscribe((roomId) => {
  console.warn('roomId', roomId)
  router.push(`/${ERoute.Game}/${roomId}`)
})

function clickRoom(roomId: string) {
  wsSend(EWsRoute.JoinRoom, { room_id: roomId })
}
</script>

<template>
  <div class="p-10">
    <h2 class="text-4xl mb-8 font-bold">選擇房間</h2>
    <div class="flex gap-8">
      <button
        @click="clickRoom(room.id)"
        v-for="room in rooms"
        :key="room.id"
        class="border rounded-lg p-7"
      >
        <h3 class="text-5xl mb-4">{{ room.name }}</h3>
        <img
          :src="
            room.image ||
            'https://t4.ftcdn.net/jpg/04/73/25/49/360_F_473254957_bxG9yf4ly7OBO5I0O5KABlN930GwaMQz.jpg'
          "
          alt=""
        />
      </button>
    </div>
  </div>
</template>
