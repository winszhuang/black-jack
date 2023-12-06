<script setup lang="ts">
const props = defineProps<{
  card: string
}>()

const imageUrl = new URL(`/src/assets/${props.card}.png`, import.meta.url).href
</script>

<template>
  <img :src="imageUrl" alt="" class="card-img" />
</template>

<style scoped>
.card-img {
  height: 20rem;
  width: 13rem;
  padding: 1rem;
}
</style>
