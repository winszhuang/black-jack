<script setup lang="ts">
import { useBlackJack } from '@/composables/use-black-jack'
import { ERoute } from '@/enums/route'
import { EWsRoute } from '@/enums/ws-route'
import { router } from '@/router'

const { rooms, onJoinRoomSuccess, wsSend } = useBlackJack()
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
  <div class="bg-primary-400 h-screen">
    <header class="bg-primary-300">
      <h2 class="text-white-50 font-light text-3xl px-8 py-3">LOBBY</h2>
    </header>
    <section class="p-10">
      <h3 class="text-white-100 mb-4">
        Feel Free To Select The Room
        <tr></tr>
        And, Play Fun ~
      </h3>
      <div class="flex gap-8">
        <button
          @click="clickRoom(room.id)"
          v-for="room in rooms"
          :key="room.id"
          class="border-0 rounded-2xl p-7 bg-primary-100 text-left"
        >
          <img
            :src="
              room.image ||
              'https://media.istockphoto.com/id/1420518918/vector/blackjack-icon-flat-style-vector-illustration-isolated-on-white-background.jpg?s=612x612&w=0&k=20&c=a3WpuAt-cEzO1rNAPasLfi6FgwTTepKNm10194yOEUg='
            "
            class="rounded-lg mb-2"
            alt=""
          />
          <h3 class="text-2xl mb-4 text-white-50 font-thin">{{ room.name }}</h3>
        </button>
      </div>
    </section>
  </div>
</template>
