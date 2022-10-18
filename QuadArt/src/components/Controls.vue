<script setup lang="ts">
import iro from '@jaames/iro'
import { ref, onMounted } from 'vue';

import VueSlider from 'vue-slider-component'
import 'vue-slider-component/theme/default.css'

interface Props {
    leafSize?: number
    roundedCorner?: number
    errorThreshold?: number
}
const props = defineProps<Props>();

const emit = defineEmits([
    "backgroundColor",
    "parameter",
    "control",
    "imageControl",
])

function initColorPicker() {
    // iro.js
    const picker = document.getElementById('picker')
    const width = picker.clientWidth - 40;
    const colorPicker = iro.ColorPicker('#picker', {
        width: width,
        layoutDirection: 'horizontal',
    });
    colorPicker.on('input:end', (color) => {
        emit("backgroundColor", color.hexString)
    });
}

const slider = {
    // contained: true,
    tooltip: 'focus',
    useKeyboard: false,
    lazy: true,
    dragOnClick: true,
}

const roundedCornerSlider = {
    ...slider,
    min: 0,
    max: 50,
    marks: {
        0: "0px",
        50: "50px",
    },
}
const roundedCornerVal = ref(roundedCornerSlider.min);
function handleRoundedCornerChange(val) {
    emit('parameter', 'roundedCorner', val);
}

const leafSizeSlider = {
    ...slider,
    min: 4,
    max: 64,
    interval: 4,
    marks: {
        4: 4,
        8: 8,
        16: 16,
        32: 32,
        64: 64,
    },
}
const leafSizeVal = ref(leafSizeSlider.min);
function handleLeafSizeChange(val) {
    emit('parameter', 'leafSize', val);
}

const errorThresholdSlider = {
    ...slider,
    min: 50,
    max: 2500,
    interval: 50,
    marks: {
        50: 50,
        2500: 2500,
    },
}
const errorThresholdVal = ref(errorThresholdSlider.min);
function handleErrorThresholdChange(val) {
    emit('parameter', 'errorThreshold', val);
}

function imageUploadClick() {
    document.getElementById('image-upload-input').click();
}

function handleImageUpload(file:File) {
    const imageTypeRe = /image.*/;
    if (!file.type.match(imageTypeRe)) {
        alert('please choose image file');
        return;
    }

    const reader = new FileReader();
    reader.onload = readerEvent => {
        emit('imageControl', 'upload', readerEvent.target.result);
        // img.src = readerEvent.target.result;
    }
    // use readAsDataURL for now
    // https://stackoverflow.com/a/31743665
    reader.readAsDataURL(file)
}

function clampSlider(n:number, sliderCfg:{min:number, max:number}):number {
    return Math.min(Math.max(n, sliderCfg.min), sliderCfg.max);
}

onMounted(() => {
    document.getElementById('image-upload-input').onchange = e => {
        const file = (e.target as HTMLInputElement).files.item(0);
        handleImageUpload(file);
    };

    initColorPicker();

    roundedCornerVal.value = clampSlider(props.roundedCorner || 0, roundedCornerSlider);
    leafSizeVal.value = clampSlider(props.leafSize || 0, leafSizeSlider);
    errorThresholdVal.value = clampSlider(props.errorThreshold || 0, errorThresholdSlider);
});
</script>

<template>
    <div class="menu">
        <div class="item">
            <h4 class="header">BG-Color Picker</h4>
            <div id="picker"></div>
        </div>
        <div class="item">
            <h4 class="header">Rounded Corner</h4>
            <div class="slider-group">
                <div class="slider">
                    <vue-slider
                        v-bind="roundedCornerSlider"
                        v-model="roundedCornerVal"
                        @change="handleRoundedCornerChange"
                    >
                    </vue-slider>
                </div>
            </div>
        </div>
        <div class="item">
            <h4 class="header">Leaf Size</h4>
            <div class="slider-group">
                <div class="slider">
                    <vue-slider
                        v-bind="leafSizeSlider"
                        v-model="leafSizeVal"
                        @change="handleLeafSizeChange"
                    >
                    </vue-slider>
                </div>
            </div>
        </div>
        <div class="item">
            <h4 class="header">Error Threshold</h4>
            <div class="slider-group">
                <div class="slider">
                    <vue-slider
                        v-bind="errorThresholdSlider"
                        v-model="errorThresholdVal"
                        @change="handleErrorThresholdChange"
                    >
                    </vue-slider>
                </div>
            </div>
        </div>
        <div class="item">
            <h4 class="header">Control</h4>
            <div class="four buttons">
                <button title="Start" @click.prevent="$emit('control', 'start')">
                    <font-awesome-icon icon="fa-solid fa-play"></font-awesome-icon>
                </button>
                <button title="Pause" @click.prevent="$emit('control', 'pause')">
                    <font-awesome-icon icon="fa-solid fa-pause"></font-awesome-icon>
                </button>
                <button title="Step" @click.prevent="$emit('control', 'step')">
                    <font-awesome-icon icon="fa-solid fa-forward-step"></font-awesome-icon>
                </button>
                <button title="Reset" @click.prevent="$emit('control', 'reset')">
                    <font-awesome-icon icon="fa-solid fa-repeat"></font-awesome-icon>
                </button>
            </div>
        </div>
        <div class="item">
            <h4 class="header">Image Control</h4>
            <div class="two buttons">
                <button title="Upload" @click.prevent="imageUploadClick">
                    <font-awesome-icon icon="fa-solid fa-upload"></font-awesome-icon>
                </button>
                <input id="image-upload-input" type="file" style="display:none;">
                <button title="Save" @click.prevent="$emit('imageControl', 'save')">
                    <font-awesome-icon icon="fa-solid fa-file-arrow-down"></font-awesome-icon>
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* @import "./base.css"; */

.menu {
    margin: 0;
    overflow-y: auto;
    border-right: 1px solid var(--color-border);
    border-top: .5px solid var(--color-border);
    border-bottom: .5px solid var(--color-border);;
}

.item {
    width: 100%;
    padding: .75rem;
    padding-right: 1rem;
    border-top: .5px solid var(--color-border);
    border-bottom: .5px solid var(--color-border);
}
.item:first-of-type {
    border-top: 0;
}
.item:last-of-type {
    border-bottom: 1px solid var(--color-border);
}

.header {
    margin-bottom: 0.5rem;
    font-weight: 800;
}

/* button */
button {
    margin: 0;
    border: 0;
    padding: 0.4rem;
}
button {
    border: 1px solid transparent;
    background-color: var(--color-background-mute);
    color: var(--vt-c-text-light-2);
}
button:hover {
    color: var(--vt-c-text-light-1);
    border: 1px solid var(--color-border-hover);
}

/* buttons */
.buttons {
    align-content: center;
    margin: 0;
    padding: 0;
    border: 1px solid var(--color-border);
}
.buttons > button {
    border: 1px solid transparent;
    border-left: .5px solid var(--color-border);
    border-right: .5px solid var(--color-border);
}
.buttons > button:first-of-type {
    border-left: .5px solid transparent;
}
.buttons > button:last-of-type {
    border-right: .5px solid transparent;
}
.buttons > button:hover {
    border: 1px solid var(--color-border-hover);
}
/* TODO: add radius to buttons first and last child */

.four.buttons > button {
    width: 25%;
}

.two.buttons > button {
    width: 50%;
}

.slider-group {
    display: flex;
    align-items: start;
    justify-content: space-between;
    height: 2.5rem;
}

.slider-group > .slider {
    width: calc(100% - 2rem);
    margin: 0 auto 0 auto;
}
</style>
