<script lang="ts">
function rgbToHex(r:number, g:number, b:number) {
    return '#' + ((1 << 24) + (r << 16) + (g << 8) + (b)).toString(16).slice(1);
}

function se(c1:{r:number, g:number, b:number}, c2:{r:number, g:number, b:number}) {
    const rErr = c1.r - c2.r;
    const gErr = c1.g - c2.g;
    const bErr = c1.b - c2.b;
    return 0.2989 * (rErr*rErr) + 0.5870 * (gErr*gErr) + 0.1140 * (bErr*bErr);
}

function calcAverageColor(imageData:Uint8ClampedArray) {
    let cumR = 0, cumG = 0, cumB = 0;
    for (let i = 0; i < imageData.length; i+=4) {
        cumR += imageData[i];
        cumG += imageData[i+1];
        cumB += imageData[i+2];
    }
    const area = imageData.length / 4;
    const avgR = Math.round(cumR / area);
    const avgG = Math.round(cumG / area);
    const avgB = Math.round(cumB / area);
    return {avgR, avgG, avgB}
}

function calcColorMSE(imageData:Uint8ClampedArray, {r, g, b}) {
    let mse = 0;
    for (let i = 0; i < imageData.length; i+=4) {
        mse += se(
            {r: r, g: g, b: b},
            {r: imageData[i], g: imageData[i+1], b: imageData[i+2]}
        )
    }
    return Math.sqrt(mse);
}
</script>
<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue';
import * as d3 from 'd3';
import { computed } from '@vue/reactivity';

interface Parameter {
  leafSize?: number,  // default 12
  roundedCorner?: number,  // default 0
  errorThreshold?: number,  // default 420
  running?: boolean,
}
interface Props {
    imgSrc: string,
    param?: Parameter
}
const props = withDefaults(defineProps<Props>(), {
    param: () => <Parameter>{
        leafSize: 12,
        roundedCorner: 0,
        errorThreshold: 420,
        running: false,
    }
});

const emit = defineEmits(['errorThresholdReached']);

// define quad class
// class QuadNode {
//     id: number
//     x: number
//     y: number
//     w: number
//     h: number
//     prevColor: string
//     color: string
//     error: number

//     constructor(x:number, y:number, w:number, h:number, prevColor: string) {
//         this.id = id++
//         this.x = x;
//         this.y = y;
//         this.w = w;
//         this.h = h;
//         this.prevColor = prevColor;
//     }

//     isLeaf():boolean {
//         return this.w < props.param.leafSize || this.h < props.param.leafSize
//     }
// }

const img = new Image();
const canvas = ref(null);
let context = null;
let svg = null;
let quads = reactive([]);
let id = 0;
let errorVal = ref(0);

let updateModelTimer = null;
let updateViewTimer = null;
const updateViewFreq = 200;
const updateModelFreq = 1;

const errorInfo = computed(() => {
    const iterations = Math.max((quads.length - 1) / 3, 0);
    return 'Iterations: ' + iterations +
            ' - Shapes: ' + quads.length +
            ' - Error: ' + errorVal.value.toString();
})

function split(node) {
    if (node.leaf) {
        return
    }
    const idx = quads.indexOf(node);
    quads.splice(idx, 1);

    const halfW = node.width / 2;
    const halfH = node.height / 2;
    const x1 = node.x;
    const x2 = node.x + halfW;
    const y1 = node.y;
    const y2 = node.y + halfH;
    quads.push(
        createQuadNode(x1, y1, halfW, halfH, node.color),
        createQuadNode(x2, y1, halfW, halfH, node.color),
        createQuadNode(x1, y2, halfW, halfH, node.color),
        createQuadNode(x2, y2, halfW, halfH, node.color)
    );
}

function step() {
    // can use quadtree instead of filtering array
    const nonLeafNodes = quads.filter(n => !n.leaf)
    if (nonLeafNodes.length <= 0) {
        return
    }
    let maxE = nonLeafNodes[0];
    for (let i = 1; i < nonLeafNodes.length; i++) {
        const e = nonLeafNodes[i];
        if (e.error > maxE.error) {
            maxE = e
        }
    }
    if (maxE.error < props.param.errorThreshold) {
        emit('errorThresholdReached');
        return
    }
    split(maxE);
    errorVal.value = maxE.error.toPrecision(5);
}

function redraw(highlight:boolean = false) {
    const rect = svg.selectAll('rect').data(quads, n => n.id);
    rect.exit().remove();
    rect.enter()
        .append('rect')
        .attr('x', n => {return n.x+0.25})
        .attr('y', n => {return n.y+0.25})
        .attr('rx', props.param?.roundedCorner || 0)
        .attr('width', n => {return n.width-0.5})
        .attr('height', n => {return n.height-0.5})
        .attr('fill', n => {return highlight? '#ffffff': n.prevColor})
        .transition().duration(500).styleTween('fill', n => {
            return d3.interpolate(highlight? '#ffffff': n.prevColor, n.color);
        });
}

function createQuadNode(x:number, y:number, width:number, height:number, prevColor:string) {
    id++;  // dont use global id?
    const data = context.getImageData(x, y, width, height).data;  // dont use global context?
    const {avgR, avgG, avgB} = calcAverageColor(data);
    const mse = calcColorMSE(data, {r: avgR, g: avgG, b: avgB})
    const colorHex = rgbToHex(avgR, avgG, avgB);

    return {
        id: id,
        error: mse,
        color: colorHex,
        prevColor: prevColor,
        x: x, y: y, width: width, height: height,
        leaf: width < props.param.leafSize || height < props.param.leafSize,
    };
}

function reset() {
    const ratio = img.width / img.height;
    let width = img.width;
    let height = img.height;
    // TODO: take care where both are bigger
    if (width > 1024) {
        width = 1024
        height = Math.round(width / ratio);
    } else if (height > 1024) {
        height = 1024
        width = Math.round(height * ratio);
    }

    svg.attr('viewBox', '0 0 ' + width + ' ' + height)
        .attr('width', width)
        .attr('height', height);
    context.canvas.width = width;
    context.canvas.height = height;
    context.drawImage(img, 0, 0, width, height);

    quads = [];
    quads.push(createQuadNode(0, 0, width, height, '#000000'));
    redraw();
}

function updateModel() {
    clearTimeout(updateModelTimer);
    if (props.param.running) {
        step();
        updateModelTimer = setTimeout(updateModel, updateModelFreq);
    }
}

function updateView() {
    clearTimeout(updateViewTimer);
    if (props.param.running) {
        redraw();
        updateViewTimer = setTimeout(updateView, updateViewFreq);
    }
}

onMounted(() => {
    // will be read frequently so use software canvas for performance
    context = canvas.value.getContext('2d', {willReadFrequently: true});
    svg = d3.select('svg');

    img.src = props.imgSrc;
    img.onload = () => {
        reset();
    }
});

// interestingly on HMR, this does not trigger because the running state is defined on parent
watch(() => props.param.running, (newValue) => {
    console.log("running status changed: "+ props.param.running);
    if (newValue) {
        updateView();
        updateModel();
    } else {
        redraw(true);
    }
});
</script>

<template>
    <div class="canvas">
         <!-- class="ui content segment"> -->
        <canvas id="source" ref="canvas"></canvas>
        <svg id="target"></svg>
        <div id="content">
            <p id="info">{{ errorInfo || 'placeholder'}}</p>
        </div>
    </div>
</template>

<style scoped>
.canvas {
    width: fit-content;
    align-items: center;
    text-align: center;
}

#source {
    display: none;
    /* max-width: 70%; */
    /* outline: 5px solid black; */
}

#target {
    /* display: none; */
    /* max-width: 70%; */
    outline: 5px solid black;
}
</style>
