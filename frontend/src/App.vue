<script setup lang="ts">
import PlayerZone from '@/components/PlayerZone.vue'
import { ECard } from '@/enums/card'
import { EMsgCode } from '@/enums/msg-code'
import { useWs } from '@/composables/use-ws'
import { notify } from '@/utils/toast'

const ws = useWs<EMsgCode>({
  dev: 'ws://localhost:8080/ws',
  prod: 'ws://localhost:8080'
})
ws.on(EMsgCode.OneJoin, (res) => {
  notify(res.message)
  console.log('res.message', res.message)
  console.log('res.data', res.data)
})

ws.on(EMsgCode.BroadcastJoin, (res) => {
  notify(res.message)
  console.log('res.message', res.message)
  console.log('res.data', res.data)
})

ws.on(EMsgCode.UpdatePlayersDetail, (res) => {
  console.log(res.data)
})
const me = 'jiljiljil'

type PlayerData = {
  userId: string
  name: string
  cards: Card[]
}

const players: PlayerData[] = [
  {
    userId: 'jiljiljil',
    name: 'wins',
    cards: [{ type: ECard.AC, value: 11 }]
  },
  {
    userId: 'opilhilj',
    name: 'tina',
    cards: [{ type: ECard._4S, value: 4 }]
  },
  {
    userId: 'whukhui',
    name: 'reg',
    cards: [{ type: ECard._5H, value: 5 }]
  }
]
// import Card from '@/components/Card.vue'
// import { ECard } from '@/constants/card.ts'

const onReady = () => ws.send(EMsgCode.OneReady)
const onHit = () => ws.send(EMsgCode.OneHit)
const onStand = () => ws.send(EMsgCode.OneStand)
</script>

<template>
  <h1 class="title">BLACK JACK</h1>
  <div class="main">
    <h2><span id="command">Gambling Time</span></h2>
    <div class="row1">
      <PlayerZone
        v-for="player in players"
        :key="player.userId"
        :user-id="player.userId"
        :name="player.name"
        :cards="player.cards"
        :is-me="player.userId === me"
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
