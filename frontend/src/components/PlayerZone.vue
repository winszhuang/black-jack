<script setup lang="ts">
import { EPlayerState } from '@/enums/player-status'
import Card from './Card.vue'

const props = defineProps<{
  userId: string
  name: string
  isMe: boolean
  cards: Card[]
  userState: EPlayerState
}>()
</script>

<template>
  <div :id="userId">
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
  </div>
</template>
