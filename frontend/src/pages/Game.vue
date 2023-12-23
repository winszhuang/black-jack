<script setup lang="ts">
import PlayerZone from '@/components/PlayerZone.vue'
import { useBlackJack } from '@/composables/use-black-jack'
import { notify } from '@/utils/toast'

const { playersDetail, myId, messageList, onSomeOneJoinRoom, onReady, onHit, onStand } =
  useBlackJack()

onSomeOneJoinRoom.subscribe((userData) => {
  notify(`玩家[${userData.name}]進入房間 ~`)
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
    <div class="row3">
      <h3>Scoreboard:</h3>
      <table>
        <tr>
          <th>Wins</th>
          <th>Losses</th>
          <th>Draws</th>
        </tr>
        <tr>
          <td><span id="wins">0</span></td>
          <td><span id="losses">0</span></td>
          <td><span id="draws">0</span></td>
        </tr>
      </table>
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
