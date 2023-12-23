<script setup lang="ts">
import PlayerZone from '@/components/PlayerZone.vue'
import { useBlackJack } from '@/composables/use-black-jack'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const {
  playersDetail,
  rooms,
  myId,
  messageList,
  onSomeOneJoinRoom,
  onSomeOneStand,
  onSomeOneWin,
  onReady,
  onHit,
  onStand,
  pushNotify
} = useBlackJack()

const route = useRoute()
const roomName = computed(() => {
  const room = rooms.value.find((room) => room.id === route.params.roomId)
  return room?.name
})

onSomeOneJoinRoom.subscribe((userData) => {
  pushNotify(`玩家[${userData.user_name}]進入房間 ~`)
})

onSomeOneStand.subscribe((userData) => {
  pushNotify(`玩家[${userData.user_name}]停止要牌`)
})

onSomeOneWin.subscribe((users) => {
  const userNames = users.reduce((curr, prev) => curr + prev.user_name + ' ', '')
  if (users.length === 0) {
    pushNotify('沒有任何玩家獲勝 ~')
  } else {
    pushNotify(`玩家[${userNames}]獲得勝利!!`)
  }
})
</script>

<template>
  <MessageList :message-list="messageList" />
  <div class="bg-primary-400 h-screen">
    <header class="bg-primary-300">
      <h1 class="text-white-50 font-light text-3xl px-8 py-3">{{ roomName || 'ROOM' }}</h1>
    </header>
    <div class="p-10">
      <div class="bg-primary-100 rounded-xl">
        <div class="flex flex-wrap flex-row p-4 border-b-2 min-h-[700px] gap-6">
          <PlayerZone
            v-for="player in playersDetail"
            :key="player.id"
            :user-id="player.id"
            :name="player.name"
            :cards="player.deck"
            :user-state="player.state"
            :is-me="player.id === myId"
            :on-ready="onReady"
            :on-hit="onHit"
            :on-stand="onStand"
          >
          </PlayerZone>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
