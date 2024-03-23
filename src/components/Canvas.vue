<script lang="ts">
class QuadNode {
    error: number;
    color: string;
    prevColor: string;
    x: number;
    y: number;
    w: number;
    h: number;
    leaf: any;
    leafSize: number;

    constructor(context2D: CanvasRenderingContext2D, x: number, y: number, w: number, h: number, prevColor: string, leafSize: number) {
        const data = context2D.getImageData(x, y, w, h).data;
        const color = this.calcAverageColor(data);
        const mse = this.calcColorMSE(data, color)
        const colorHex = this.colorToHexString(color);

        this.error = mse;
        this.color = colorHex;
        this.prevColor = prevColor ?? "#000000";
        this.x = x
        this.y = y
        this.w = w
        this.h = h
        this.leaf = w < leafSize || h < leafSize
        this.leafSize = leafSize;
    }

    private colorToHexString(color: Color): string {
        return '#' + ((1 << 24) + (color.r << 16) + (color.g << 8) + (color.b)).toString(16).slice(1);
    }

    private calcAverageColor(imageData: Uint8ClampedArray): Color {
        let cumR = 0, cumG = 0, cumB = 0;
        for (let i = 0; i < imageData.length; i += 4) {
            cumR += imageData[i];
            cumG += imageData[i + 1];
            cumB += imageData[i + 2];
        }
        const area = imageData.length / 4;
        const avgR = Math.round(cumR / area);
        const avgG = Math.round(cumG / area);
        const avgB = Math.round(cumB / area);
        return { r: avgR, g: avgG, b: avgB }
    }

    private calcColorDiff(c1: Color, c2: Color): number {
        const rErr = c1.r - c2.r;
        const gErr = c1.g - c2.g;
        const bErr = c1.b - c2.b;
        return 0.2989 * (rErr * rErr) + 0.5870 * (gErr * gErr) + 0.1140 * (bErr * bErr);
    }

    private calcColorMSE(imageData: Uint8ClampedArray, color: Color): number {
        let mse = 0;
        for (let i = 0; i < imageData.length; i += 4) {
            mse += this.calcColorDiff(
                color,
                { r: imageData[i], g: imageData[i + 1], b: imageData[i + 2] }
            )
        }
        return Math.sqrt(mse);
    }

    Split(context2D: CanvasRenderingContext2D): QuadNode[] {
        if (this.leaf) {
            return []
        }

        const halfW = this.w / 2;
        const halfH = this.h / 2;
        const x1 = this.x;
        const x2 = this.x + halfW;
        const y1 = this.y;
        const y2 = this.y + halfH;
        
        return [
            new QuadNode(context2D, x1, y1, halfW, halfH, this.color, this.leafSize),
            new QuadNode(context2D, x2, y1, halfW, halfH, this.color, this.leafSize),
            new QuadNode(context2D, x1, y2, halfW, halfH, this.color, this.leafSize),
            new QuadNode(context2D, x2, y2, halfW, halfH, this.color, this.leafSize),
        ]
    }
}

function rgbToHex(r: number, g: number, b: number) {
    return '#' + ((1 << 24) + (r << 16) + (g << 8) + (b)).toString(16).slice(1);
}

interface Color {
    r: number
    g: number
    b: number
}

function se(c1: Color, c2: Color) {
    const rErr = c1.r - c2.r;
    const gErr = c1.g - c2.g;
    const bErr = c1.b - c2.b;
    return 0.2989 * (rErr * rErr) + 0.5870 * (gErr * gErr) + 0.1140 * (bErr * bErr);
}

function calcAverageColor(imageData: Uint8ClampedArray) {
    let cumR = 0, cumG = 0, cumB = 0;
    for (let i = 0; i < imageData.length; i += 4) {
        cumR += imageData[i];
        cumG += imageData[i + 1];
        cumB += imageData[i + 2];
    }
    const area = imageData.length / 4;
    const avgR = Math.round(cumR / area);
    const avgG = Math.round(cumG / area);
    const avgB = Math.round(cumB / area);
    return { avgR, avgG, avgB }
}

function calcColorMSE(imageData: Uint8ClampedArray, color: Color) {
    let mse = 0;
    for (let i = 0; i < imageData.length; i += 4) {
        mse += se(
            color,
            { r: imageData[i], g: imageData[i + 1], b: imageData[i + 2] }
        )
    }
    return Math.sqrt(mse);
}

function convertRemToPixel(rem: number): number {
    return rem * parseFloat(getComputedStyle(document.documentElement).fontSize);
}
</script>
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import type {Ref} from 'vue'
import { select, interpolate } from 'd3';
import { computed } from '@vue/reactivity';

interface Parameter {
    leafSize: number  // default 12
    roundedCorner: number  // default 0
    errorThreshold: number  // default 420
    backgroundColor: string
}

interface Props {
    imgSrc: string
    param: Parameter
    running?: boolean
}

const props = withDefaults(defineProps<Props>(), {
    param: () => <Parameter>{
        leafSize: 12,
        roundedCorner: 0,
        errorThreshold: 420,
        running: false,
        backgroundColor: '#ffffff',
    },
    running: false
});

// const emit = defineEmits(['errorThresholdReached']);

const img = new Image();
const canvas: Ref<HTMLCanvasElement | null> = ref(null);
let context: CanvasRenderingContext2D | null = null;
let svg = null;
let quads = ref([]);
let tmp : Ref<QuadNode[]> = ref([]);
let id = 0;
let errorVal = ref('0');

let updateModelTimer: number | null = null;
let updateViewTimer: number | null = null;
const updateViewFreq = 200;
const updateModelFreq = 1;

const errorInfo = computed(() => {
    const iterations = Math.max((quads.value.length - 1) / 3, 0);
    return 'Iterations: ' + iterations +
        ' - Shapes: ' + quads.value.length +
        ' - Error: ' + errorVal.value;
})

function createQuadNode(x: number, y: number, width: number, height: number, prevColor: string) {
    id++;  // dont use global id?
    const data = context.getImageData(x, y, width, height).data;  // dont use global context?
    const { avgR, avgG, avgB } = calcAverageColor(data);
    const mse = calcColorMSE(data, { r: avgR, g: avgG, b: avgB })
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
    const canvasE = document.getElementsByClassName('canvas');
    if (canvasE.length <= 0) {
        return
    }

    const parentE = canvasE[0].parentElement as HTMLElement;
    const maxWidth = parentE.clientWidth - convertRemToPixel(2);
    const maxHeight = document.documentElement.clientHeight - convertRemToPixel(3);  // 2 padding + 1 text

    let width = img.width;
    let height = img.height;
    const ratio = width / height;
    if (width > maxWidth) {
        width = maxWidth
        height = Math.round(width / ratio);
    }
    if (height > maxHeight) {
        height = maxHeight
        width = Math.round(height * ratio);
    }
    console.log(parentE.clientWidth, parentE.clientHeight);
    console.log(width, height);

    svg.attr('viewBox', '0 0 ' + width + ' ' + height)
        .attr('width', width)
        .attr('height', height)
        .style('background-color', props.param.backgroundColor);
    context.canvas.width = width;
    context.canvas.height = height;
    context.drawImage(img, 0, 0, width, height);

    quads.value = [];
    quads.value.push(createQuadNode(0, 0, width, height, '#000000'));
    errorVal.value = '0';

    redraw();
}

function split(node) {
    if (node.leaf) {
        return
    }
    const idx = quads.value.indexOf(node);
    quads.value.splice(idx, 1);

    const halfW = node.width / 2;
    const halfH = node.height / 2;
    const x1 = node.x;
    const x2 = node.x + halfW;
    const y1 = node.y;
    const y2 = node.y + halfH;
    quads.value.push(
        createQuadNode(x1, y1, halfW, halfH, node.color),
        createQuadNode(x2, y1, halfW, halfH, node.color),
        createQuadNode(x1, y2, halfW, halfH, node.color),
        createQuadNode(x2, y2, halfW, halfH, node.color)
    );
}

function step() {
    // can use quadtree instead of filtering array
    const nonLeafNodes = quads.value.filter(n => !n.leaf)
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
        // emit('errorThresholdReached');
        return
    }
    split(maxE);
    errorVal.value = maxE.error.toPrecision(5);
}

function tmpStep() {
    const nodes = tmp.value
    if (nodes.length <= 0) {
        return
    }
    // use heap
    let maxE = nodes[0];
    for (let i = 1; i < nodes.length; i++) {
        const e = nodes[i];
        if (e.error > maxE.error) {
            maxE = e
        }
    }
    if (maxE.error < props.param.errorThreshold) {
        // emit('errorThresholdReached');
        return
    }
    maxE.Split(context!)
    errorVal.value = maxE.error.toPrecision(5);
}

function redraw(highlight: boolean = false) {
    const rect = svg.selectAll('rect').data(quads.value, n => n.id);
    rect.exit().remove();
    rect.enter()
        .append('rect')
        .attr('x', n => { return n.x + 0.25 })
        .attr('y', n => { return n.y + 0.25 })
        .attr('rx', props.param?.roundedCorner || 0)
        .attr('width', n => { return n.width - 0.5 })
        .attr('height', n => { return n.height - 0.5 })
        .attr('fill', n => { return highlight ? '#ffffff' : n.prevColor })
        .transition().duration(500).styleTween('fill', n => {
            return interpolate(highlight ? '#ffffff' : n.prevColor, n.color);
        });
}

function updateModel() {
    if (updateModelTimer !== null) {
        clearTimeout(updateModelTimer);
    }
    if (props.running) {
        step();
        updateModelTimer = setTimeout(updateModel, updateModelFreq);
    }
}

function updateView() {
    if (updateViewTimer !== null) {
        clearTimeout(updateViewTimer);
    }
    if (props.running) {
        redraw();
        updateViewTimer = setTimeout(updateView, updateViewFreq);
    }
}

function save() {
    const svgE = document.getElementById('target') as HTMLOrSVGElement as SVGElement as SVGGraphicsElement;
    const { width: w, height: h } = svgE.getBBox();

    // clone svg and turn to blobURL
    svgE.setAttribute('xmlns', "http://www.w3.org/2000/svg");  // needed otherwise chrome will not work
    const outerHTML = svgE.outerHTML;
    const blob = new Blob([outerHTML], { type: 'image/svg+xml;charset=utf-8' });
    const URL = window.URL || window.webkitURL;
    const blobURL = URL.createObjectURL(blob);

    // create new image, canvas, and draw blobURL to canvas
    let copyImg = new Image();
    copyImg.onload = () => {
        let tmpCanvas = document.createElement('canvas');
        tmpCanvas.width = w;
        tmpCanvas.height = h;

        let tmpContext = tmpCanvas.getContext('2d');
        tmpContext!.drawImage(copyImg, 0, 0, w, h);

        // save canvas
        let png = tmpCanvas.toDataURL();
        let download = function (href: string, name: string) {
            const link = document.createElement('a');
            link.download = name;
            document.body.append(link);
            link.href = href;
            link.click();
            link.remove();
        }
        download(png, "image.png");
    }
    copyImg.src = blobURL;
}

onMounted(() => {
    // will be read frequently so use software canvas for performance
    context = canvas.value?.getContext('2d', { willReadFrequently: true }) ?? null;
    svg = select('#target');

    img.onload = () => reset()
    img.src = props.imgSrc;
});

watch(() => props.imgSrc, (newImgSrc) => img.src = props.imgSrc);

watch(() => props.param.backgroundColor, (newColor) => svg.style('background-color', newColor));

// interestingly on HMR, this does not trigger because the running state is defined on parent
watch(() => props.running, (newValue) => {
    // console.log("running status changed: "+ props.running);
    if (newValue) {
        updateView();
        updateModel();
    } else {
        redraw(true);
    }
});

defineExpose({
    step,
    reset,
    redraw,
    save,
});
</script>

<template>
    <div class="canvas">
        <canvas id="source" ref="canvas"></canvas>
        <svg id="target"></svg>
        <div id="content">
            <p id="info">{{ errorInfo || 'placeholder' }}</p>
        </div>
    </div>
</template>

<style scoped>
.canvas {
    width: fit-content;
    text-align: center;
    margin: 0;
    padding: 0;
}

#source {
    display: none;
}

#target {
    outline: 5px solid black;
}
</style>
