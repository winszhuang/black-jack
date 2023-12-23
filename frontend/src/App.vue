<script setup lang="ts">
import { useBlackJack } from '@/composables/use-black-jack'
import MessageList from '@/components/MessageList.vue'
import { router } from '@/router'
import { ERoute } from '@/enums/route'
import { useRoute } from 'vue-router'

const route = useRoute()
const { messageList, onConnectSuccess, myId } = useBlackJack()

onConnectSuccess.subscribe((info) => {
  console.log('loginInfo', info)
  info.is_login ? handleIsLogin(info) : handleNotLogin()
})

const handleIsLogin = (info: LoginData) => {
  myId.value = info.user_id
  if (route.name === ERoute.Entrance) {
    router.push({ name: ERoute.Lobby })
  }
}

const handleNotLogin = () => {
  localStorage.removeItem('access_token')
  if (route.name !== ERoute.Entrance) {
    router.push({ name: ERoute.Entrance })
  }
}
</script>

<template>
  <MessageList :message-list="messageList" />
  <RouterView />
</template>
