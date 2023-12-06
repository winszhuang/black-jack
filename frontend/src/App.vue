<script setup lang="ts">
import PlayerZone from '@/components/PlayerZone.vue'
import { EMsgCode } from '@/enums/msg-code'
import { useWs } from '@/composables/use-ws'
import { notify } from '@/utils/toast'
import { ref } from 'vue'

const ws = useWs<EMsgCode>({
  dev: 'ws://localhost:8080/ws',
  prod: 'ws://localhost:8080'
})

const myId = ref('')
const playersDetail = ref<PlayerDetail[]>([])

ws.on(EMsgCode.ClientJoin, (res) => {
  notify(res.message)
  myId.value = res.data as string
})

ws.on(EMsgCode.BroadcastJoin, (res) => {
  notify(res.message)
  console.log('res.message', res.message)
  console.log('res.data', res.data)
})

ws.on(EMsgCode.BroadcastLeave, (res) => {
  notify(res.message)
})

ws.on(EMsgCode.BroadcastGameStart, (res) => {
  notify(res.message)
})

ws.on(EMsgCode.BroadcastGameOver, (res) => {
  notify(res.message)
})

ws.on(EMsgCode.UpdatePlayersDetail, (res) => {
  console.log(res.data)
  playersDetail.value = res.data as PlayerDetail[]
})

const onReady = () => ws.send(EMsgCode.ClientReady, '123')
const onHit = () => ws.send(EMsgCode.ClientHit)
const onStand = () => ws.send(EMsgCode.ClientStand)
</script>

<template>
  <h1 class="title">BLACK JACK</h1>
  <div class="main">
    <h2><span id="command">Gambling Time</span></h2>
    <div class="row1">
      <PlayerZone
        v-for="player in playersDetail"
        :key="player.id"
        :user-id="player.id"
        :name="player.id"
        :cards="player.deck"
        :user-state="player.state"
        :is-me="player.id === myId"
      >
      </PlayerZone>
    </div>
    <div class="row2">
      <div class="buttons">
        <button class="btn-lg btn-danger" id="ready" @click="onReady">Ready</button>
        <button class="btn-lg btn-success" id="hit" @click="onHit">Hit</button>
        <button class="btn-lg btn-warning" id="stand" @click="onStand">Stand</button>
      </div>
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
