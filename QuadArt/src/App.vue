<script setup lang="ts">
import {onMounted, ref, reactive} from 'vue'
import Controls from './components/Controls.vue'
import Canvas from './components/Canvas.vue'
import imgUrl from '@/assets/tmp.jpg'

const canvas = ref(null);

const running = ref(false);
const param = reactive({
  leafSize: 12,  // default 12
  roundedCorner: 0,  // default 0
  errorThreshold: 500,  // default 420
  backgroundColor: "#ffffff",
});

const imgSrc = ref(imgUrl);

function toggleRunning(state?: boolean) {
  running.value = state == undefined? !running.value: state;
}

onMounted(() => {
  toggleRunning(true);
})

function handleBackgroundColor(color) {
  param.backgroundColor = color;
}

function handleParameterChange(paramName, val) {
  switch(paramName) {
    case "roundedCorner":
      param.roundedCorner = val
      break;
    case "leafSize":
      param.leafSize = val
      break;
    case "errorThreshold":
      param.errorThreshold = val
      break;
    default:
  }
  canvas.value.reset();
  toggleRunning(true);
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

function handleImageUpload(file:File) {
    const imageTypeRe = /image.*/;
    if (!file.type.match(imageTypeRe)) {
        alert('please choose image file');
        return;
    }

    const reader = new FileReader();
    reader.onload = readerEvent => {
        imgSrc.value = readerEvent.target.result as string
    }
    // use readAsDataURL for now
    // https://stackoverflow.com/a/31743665
    reader.readAsDataURL(file)
}

function handleImageControl(imageControlName, ...args) {
  switch(imageControlName) {
    case "upload":
      if (args.length != 1) {
        return;  // unexpected
      }
      handleImageUpload(args[0] as File)
      break;
    case "save":
      canvas.value.save();
      break;
    default:
      return
  }
};

const highlight = ref(false);
function handleDrop(e:DragEvent) {
    highlight.value = false;
    const files = e.dataTransfer.files
    if (files.length > 0) {
        handleImageUpload(files[0]);
    }
}
</script>

<template>
  <Controls
    @background-color="handleBackgroundColor"
    @control="handleControl"
    @image-control="handleImageControl"
    @parameter="handleParameterChange"
    :rounded-corner="param.roundedCorner"
    :leaf-size="param.leafSize"
    :error-threshold="param.errorThreshold"
    />
  <main
    :class="{ highlight: highlight }"
    @dragenter.stop.prevent="highlight=true"
    @dragover.stop.prevent="highlight=true"
    @dragleave.stop.prevent="highlight=false"
    @drop.stop.prevent="handleDrop"
    >
    <Canvas
      ref="canvas"
      :img-src='imgSrc'
      :param="param"
      :running="running"
      @error-threshold-reached="toggleRunning(false)"
      ></Canvas>
  </main>
</template>

<style scoped>
main {
  padding: 1rem;
}
</style>
