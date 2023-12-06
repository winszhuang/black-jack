<!-- eslint-disable func-call-spacing -->
<script setup lang="ts">
import { EPlayerState } from '@/enums/player-status'
import Card from './Card.vue'

const props = defineProps<{
  userId: string
  name: string
  isMe: boolean
  cards: Card[]
  userState: EPlayerState
  onReady: () => void
  onHit: () => void
  onStand: () => void
}>()
</script>

<template>
  <div :id="userId" :class="isMe && ' bg-lime-50'" class="relative">
    <h2 :class="isMe && 'text-lime-300'">
      {{ props.name }}:
      <span :id="`score-${props.userId}`">0</span>
    </h2>
    <div class="text-2xl">state: {{ EPlayerState[userState] }}</div>
    <div :id="`box-${props.userId}`" class="flex flex-wrap">
      <template v-for="c in props.cards" :key="c">
        <Card :card="c.name"></Card>
      </template>
    </div>
    <div class="row2 absolute bottom-0 right-0 left-0" v-if="isMe">
      <div class="buttons">
        <button class="btn-lg btn-danger" id="ready" @click="onReady">Ready</button>
        <button class="btn-lg btn-success" id="hit" @click="onHit">Hit</button>
        <button class="btn-lg btn-warning" id="stand" @click="onStand">Stand</button>
      </div>
    </div>
  </div>
</template>
