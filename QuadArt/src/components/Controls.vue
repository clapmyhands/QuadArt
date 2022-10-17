<script setup lang="ts">
defineEmits([
    "controlChange",
])
</script>

<template>
    <div class="menu">
        <div class="item"></div>
        <div class="item"></div>
        <div class="item">
            <h4 class="header">Control</h4>
            <div class="four buttons">
                <button title="Start" @click.prevent="$emit('controlChange', 'start')">
                    <font-awesome-icon icon="fa-solid fa-play"></font-awesome-icon>
                </button>
                <button title="Pause" @click.prevent="$emit('controlChange', 'pause')">
                    <font-awesome-icon icon="fa-solid fa-pause"></font-awesome-icon>
                </button>
                <button title="Step" @click.prevent="$emit('controlChange', 'step')">
                    <font-awesome-icon icon="fa-solid fa-forward-step"></font-awesome-icon>
                </button>
                <button title="Reset" @click.prevent="$emit('controlChange', 'reset')">
                    <font-awesome-icon icon="fa-solid fa-repeat"></font-awesome-icon>
                </button>
            </div>
        </div>
        <div class="item"></div>
    </div>
</template>

<style scoped>
/* @import "./base.css"; */

.menu {
    margin: 0;
    /* display: grid; */
    border-right: 1px solid var(--color-border);
}

.item {
    width: 100%;
    padding: 1rem;
    border-top: 1px solid var(--color-border);
}

.item:first-of-type {
    border-top: 0px;
    /* 1px solid var(--vt-c-divider-light-2); */
}

.header {
    margin-bottom: 0.5rem;
    font-weight: 800;
}

/* button */
button {
    margin: 0;
    border: 0;
    padding: 0.5rem;
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
