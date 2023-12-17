<script setup lang="ts">
import { useBlackJack } from './composables/use-black-jack'
import MessageList from '@/components/MessageList.vue'

const { messageList } = useBlackJack()
</script>

<template>
  <MessageList :message-list="messageList" />
  <RouterView />
</template>
