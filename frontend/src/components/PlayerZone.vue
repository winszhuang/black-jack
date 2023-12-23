<!-- eslint-disable func-call-spacing -->
<script setup lang="ts">
import { EPlayerState } from '@/enums/player-status'
import Card from './Card.vue'
import { computed } from 'vue'

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

const canReady = computed(() => props.userState === EPlayerState.Wait)
const canHit = computed(() => props.userState === EPlayerState.Play)
const canStand = computed(() => canHit.value)
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
        <button
          :disabled="!canReady"
          class="btn-lg btn-danger disabled:opacity-50"
          id="ready"
          @click="onReady"
        >
          Ready
        </button>
        <button
          :disabled="!canHit"
          class="btn-lg btn-success disabled:opacity-50"
          id="hit"
          @click="onHit"
        >
          Hit
        </button>
        <button
          :disabled="!canStand"
          class="btn-lg btn-warning disabled:opacity-50"
          id="stand"
          @click="onStand"
        >
          Stand
        </button>
      </div>
    </div>
  </div>
</template>
