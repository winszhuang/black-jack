<script setup lang="ts">
import { useBlackJack } from '@/composables/use-black-jack'
import { ERoute } from '@/enums/route'
import { EWsRoute } from '@/enums/ws-route'
import { router } from '@/router'
import { notify } from '@/utils/toast'
import { ref } from 'vue'

const { wsSend, onLoginSuccess, onRegisterSuccess } = useBlackJack()

onLoginSuccess.subscribe((loginInfo) => {
  localStorage.setItem('access_token', loginInfo.token)
  router.push({ name: ERoute.Lobby })
})

onRegisterSuccess.subscribe(() => {
  notify('註冊成功! 直接幫你登入哦')
  setTimeout(() => {
    const req = {
      username: userInput.value.username,
      password: userInput.value.password
    }
    wsSend(EWsRoute.Login, req)
  }, 1000)
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
  <div class="bg-primary-400 flex items-center justify-center h-screen text-white-100 font-thin">
    <div class="min-w-[380px]">
      <div class="px-8 pt-10 pb-2 rounded-t-2xl bg-primary-200 text-center">
        <h2 class="text-4xl font-sans font-light text-white-50">
          {{ currentMode === Mode.Login ? 'Login' : 'Register ' }}
        </h2>
      </div>
      <div class="p-8 rounded-b-2xl bg-primary-200">
        <div class="mb-4">
          <label for="username" class="block mb-2">Username</label>
          <input
            v-model="userInput.username"
            type="text"
            id="username"
            class="bg-primary-400 border border-gray-600 p-2 rounded w-full"
            placeholder="Username"
          />
        </div>
        <div class="mb-4">
          <label for="password" class="block mb-2">Password</label>
          <input
            v-model="userInput.password"
            type="password"
            id="password"
            class="bg-primary-400 border border-gray-600 p-2 rounded w-full"
            placeholder="Password"
          />
        </div>
        <button
          @click="onSubmit"
          v-if="currentMode === Mode.Login"
          class="mt-3 py-2 px-4 rounded w-full mb-4 bg-yellow-100 text-dark-900 font-normal text-xl"
        >
          Login
        </button>
        <button
          @click="onSubmit"
          v-if="currentMode === Mode.Register"
          class="mt-3 py-2 px-4 rounded w-full mb-4 bg-yellow-100 text-dark-900 font-normal text-xl"
        >
          Register
        </button>
        <div class="mt-10 text-center">
          <span class="text-gray-300">{{
            currentMode === Mode.Login ? `Don't have an account ?` : 'Already have an account ?'
          }}</span>
          <button
            v-if="currentMode === Mode.Login"
            @click="changeMode(Mode.Register)"
            class="py-2 px-3 text-yellow-50 underline"
          >
            REGISTER
          </button>
          <button v-else @click="changeMode(Mode.Login)" class="py-2 px-3 text-yellow-50 underline">
            LOGIN
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
