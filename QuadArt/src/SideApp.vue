<script setup lang="ts">
import {onMounted, reactive} from 'vue'
import Canvas from './components/Canvas.vue'

const param = reactive({
  leafSize: 8,  // default 12
  roundedCorner: 0,  // default 0
  errorThreshold: 1000,  // default 420
  running: false,
})

function toggleRunning(state?: boolean) {
  param.running = state == undefined? !param.running: state;
}

onMounted(() => {
  toggleRunning(true);
})



</script>

<template>
  <main>
    <Canvas img-src='/tmp.jpg' :param="param" @error-threshold-reached="toggleRunning(false)"/>
  </main>
  <button @click="_ => toggleRunning(true)">Start</button>
  <button @click="_ => toggleRunning(false)">Stop</button>
</template>

<style scoped>

</style>
