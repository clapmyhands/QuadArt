<script setup lang="ts">
import {onMounted, ref, reactive} from 'vue'
import Controls from './components/Controls.vue'
import Canvas from './components/Canvas.vue'

const canvas = ref(null);

const running = ref(false);
const param = reactive({
  leafSize: 8,  // default 12
  roundedCorner: 0,  // default 0
  errorThreshold: 1000,  // default 420
  backgroundColor: "#f2f2f2",
});
const imgSrc = ref('/tmp.jpg');

function toggleRunning(state?: boolean) {
  running.value = state == undefined? !running.value: state;
}

onMounted(() => {
  toggleRunning(true);
})

function handleBackgroundColor(color) {
  param.backgroundColor = color;
}

function handleControl(controlName) {
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
};

function handleImageControl(imageControlName, ...args) {
  switch(imageControlName) {
    case "upload":
      if (args.length != 1) {
        return;  // unexpected
      }
      imgSrc.value = args[0] as string;
      break;
    case "save":
      canvas.value.save();
      break;
    default:
      return
  }
};

</script>

<template>
  <Controls
    @background-color="handleBackgroundColor"
    @control="handleControl"
    @image-control="handleImageControl"
  />
  <main>
    <Canvas
      ref="canvas"
      :img-src='imgSrc'
      :param="param"
      :running="running"
      @error-threshold-reached="toggleRunning(false)"
    />
  </main>
</template>

<style scoped>
main {
  padding: 1rem;
}
</style>
