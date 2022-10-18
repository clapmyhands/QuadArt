<script setup lang="ts">
import iro from '@jaames/iro'
import { onMounted } from 'vue';

const emit = defineEmits([
    "backgroundColor",
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

onMounted(() => {
    document.getElementById('image-upload-input').onchange = e => {
        const file = (e.target as HTMLInputElement).files.item(0);
        handleImageUpload(file);
    };

    initColorPicker();
});
</script>

<template>
    <div class="menu">
        <div class="item">
            <h4 class="header">BG-Color Picker</h4>
            <div id="picker"></div>
        </div>
        <div class="item"></div>
        <div class="item"></div>
        <div class="item"></div>
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
    border-top: .5px solid var(--color-border);
    border-bottom: .5px solid var(--color-border);
}
.item:first-of-type {
    border-top: 0;
    /* border-top: 1px solid var(--color-border); */
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
</style>
