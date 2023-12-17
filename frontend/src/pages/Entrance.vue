<script setup lang="ts">
import { useBlackJack } from '@/composables/use-black-jack'
import { ERoute } from '@/enums/route'
import { EWsRoute } from '@/enums/ws-route'
import { router } from '@/router'
import { ref } from 'vue'

const { wsSend, onLoginSuccess } = useBlackJack()

onLoginSuccess.subscribe((token) => {
  console.log('token', token)
  localStorage.setItem('access_token', token)
  router.push({ name: ERoute.Lobby })
})

const userInput = ref({
  username: '',
  password: ''
})

enum Mode {
  Register = 'Register',
  Login = 'Login'
}
const currentMode = ref(Mode.Login)
const changeMode = (mode: Mode) => {
  currentMode.value = mode
}

const onSubmit = () => {
  if (!userInput.value.username || !userInput.value.password) {
    alert('請完整輸入帳密')
    return
  }

  const req = {
    username: userInput.value.username,
    password: userInput.value.password
  }

  if (currentMode.value === Mode.Login) {
    wsSend(EWsRoute.Login, req)
  } else if (currentMode.value === Mode.Register) {
    wsSend(EWsRoute.Register, req)
  }
}
</script>

<template>
  <div class="bg-gray-800 flex items-center justify-center h-screen">
    <div class="bg-gray-700 text-white p-8 rounded-md w-96">
      <h2 class="text-xl mb-4">Please {{ currentMode }} to your account</h2>
      <div class="mb-4">
        <label for="username" class="block mb-2">Username</label>
        <input
          v-model="userInput.username"
          type="text"
          id="username"
          class="bg-gray-800 border border-gray-600 p-2 rounded w-full"
          placeholder="Username"
        />
      </div>
      <div class="mb-4">
        <label for="password" class="block mb-2">Password</label>
        <input
          v-model="userInput.password"
          type="password"
          id="password"
          class="bg-gray-800 border border-gray-600 p-2 rounded w-full"
          placeholder="Password"
        />
      </div>
      <button
        @click="onSubmit"
        v-if="currentMode === Mode.Login"
        class="bg-gradient-to-r from-orange-400 to-pink-500 text-white py-2 px-4 rounded w-full mb-4"
      >
        Login
      </button>
      <button
        v-if="currentMode === Mode.Register"
        class="bg-gradient-to-r from-orange-400 to-pink-500 text-white py-2 px-4 rounded w-full mb-4"
      >
        Register
      </button>
      <div class="text-center">
        <span class="text-gray-300">Don't have an account?</span>
        <button
          @click="changeMode(Mode.Register)"
          class="text-white border border-pink-500 py-2 px-4 rounded"
        >
          REGISTER
        </button>
      </div>
    </div>
  </div>
</template>
