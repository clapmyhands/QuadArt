<script setup lang="ts">
import {onMounted, ref, reactive} from 'vue'
import Controls from './components/Controls.vue'
import Canvas from './components/Canvas.vue'

const canvas = ref(null)

const param = reactive({
  leafSize: 8,  // default 12
  roundedCorner: 0,  // default 0
  errorThreshold: 1000,  // default 420
  running: false,
});
const imgSrc = ref('/tmp.jpg');

function toggleRunning(state?: boolean) {
  param.running = state == undefined? !param.running: state;
}

onMounted(() => {
  toggleRunning(true);
})

function handleControlChange(controlName) {
  switch(controlName) {
    case "start":
      toggleRunning(true);
      break;
    case "pause":
      toggleRunning(false);
      break;
    case "step":
      canvas.value.step();
      canvas.value.redraw(true);
      break;
    case "reset":
      canvas.value.reset();
      break;
    default:
      return
  }
}

</script>

<template>
  <Controls @control-change="handleControlChange"/>
  <main>
    <Canvas ref="canvas" :img-src='imgSrc' :param="param" @error-threshold-reached="toggleRunning(false)"/>
  </main>
</template>

<style scoped>
main {
  padding: 1rem;
}
</style>
