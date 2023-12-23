<script setup lang="ts">
import PlayerZone from '@/components/PlayerZone.vue'
import { useBlackJack } from '@/composables/use-black-jack'

const {
  playersDetail,
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
  <h1 class="title">BLACK JACK</h1>
  <div class="main">
    <h2><span id="command">Gambling Time</span></h2>
    <div class="row1">
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
