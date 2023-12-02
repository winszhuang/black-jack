<script setup lang="ts">
import PlayerZone from './components/PlayerZone.vue'
import { ECard } from './constants/card'

const me = 'jiljiljil'

type PlayerData = {
  userId: string,
  name: string,
  cards: Card[]
}

const players: PlayerData[] = [
  {
    userId: 'jiljiljil',
    name: 'wins',
    cards: [
      { type: ECard.AC, value: 11 }
    ]
  },
  {
    userId: 'opilhilj',
    name: 'tina',
    cards: [
      { type: ECard._4S, value: 4 }
    ]
  },
  {
    userId: 'whukhui',
    name: 'reg',
    cards: [
      { type: ECard._5H, value: 5 }
    ]
  }
]

// import Card from '@/components/Card.vue'
// import { ECard } from '@/constants/card.ts'

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
                <button class="btn-lg btn-danger" id="deal" onmouseover="tink.play()">Ready</button>
                <button class="btn-lg btn-success" id="hit" onmouseover="tink.play()" > Hit </button>
                <button class="btn-lg btn-warning" id="stand" onmouseover="tink.play()">Stand</button>
            </div>
        </div>
        <div class="row3">
            <h3>Scoreboard: </h3>
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
