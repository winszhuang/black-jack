<script setup lang="ts">
import { EMessageType } from '@/enums/message'

const props = defineProps<{ messageList: MessageItem[] }>()
</script>

<template>
  <Teleport to="body">
    <TransitionGroup
      name="list"
      tag="ul"
      class="absolute inset-0 pointer-events-none top-28 flex flex-col items-center gap-y-4"
    >
      <li
        v-for="message in props.messageList"
        :key="message.id"
        :class="message.type === EMessageType.Success ? 'bg-blue-500' : 'bg-red-600'"
        class="bg-blue-500 text-white font-semibold py-2 px-4 rounded shadow-md h-10"
      >
        {{ message.text }}
      </li>
    </TransitionGroup>
  </Teleport>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(-30px);
}
</style>
