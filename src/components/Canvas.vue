<script lang="ts">
interface Color {
    r: number
    g: number
    b: number
}

class QuadNode {
    static counter: number = 0;
    id: number;
    error: number;
    color: string;
    prevColor: string;
    x: number;
    y: number;
    w: number;
    h: number;
    leaf: any;
    leafSize: number;

    static GenerateID(): number {
        return QuadNode.counter++;
    }

    static ResetID() {
        QuadNode.counter = 0;
    }

    constructor(
        context2D: CanvasRenderingContext2D,
        x: number, y: number, w: number, h: number,
        leafSize: number,
        prevColor: string = "#000000",
    ) {
        const data = context2D.getImageData(x, y, w, h).data;
        const color = this.calcAverageColor(data);
        const mse = this.calcColorMSE(data, color)
        const colorHex = this.colorToHexString(color);

        this.id = QuadNode.GenerateID();
        this.error = mse;
        this.color = colorHex;
        this.prevColor = prevColor;
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
            new QuadNode(context2D, x1, y1, halfW, halfH, this.leafSize, this.color),
            new QuadNode(context2D, x2, y1, halfW, halfH, this.leafSize, this.color),
            new QuadNode(context2D, x1, y2, halfW, halfH, this.leafSize, this.color),
            new QuadNode(context2D, x2, y2, halfW, halfH, this.leafSize, this.color),
        ]
    }
}

function convertRemToPixel(rem: number): number {
    return rem * parseFloat(getComputedStyle(document.documentElement).fontSize);
}
</script>
<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import type { Ref } from 'vue'
import { select, interpolate } from 'd3';
import { Selection } from 'd3';
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
let svg: Selection<SVGGElement, QuadNode, HTMLElement, any> | null = null;
let tmp: Ref<QuadNode[]> = ref([]);
let errorVal = ref('0');

let updateModelTimer: number | null = null;
let updateViewTimer: number | null = null;
const updateViewFreq = 200;
const updateModelFreq = 1;

const errorInfo = computed(() => {
    const iterations = Math.max((tmp.value.length - 1) / 3, 0);
    return 'Iterations: ' + iterations +
        ' - Shapes: ' + tmp.value.length +
        ' - Error: ' + errorVal.value;
})

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
    // console.log(parentE.clientWidth, parentE.clientHeight);
    // console.log(width, height);

    if (context === null) {
        console.log("context is null???")
        return
    }

    svg?.attr('viewBox', '0 0 ' + width + ' ' + height)
        .attr('width', width)
        .attr('height', height)
        .style('background-color', props.param.backgroundColor);

    // reset Canvas context
    context.canvas.width = width;
    context.canvas.height = height;
    context.drawImage(img, 0, 0, width, height);

    // reset QuadNode stoere
    tmp.value = [];
    tmp.value.push(new QuadNode(context, 0, 0, width, height, props.param.leafSize));

    // reset errorVal
    errorVal.value = '0';

    redraw();
}

function step() {
    // can use quadtree instead of filtering array
    const nonLeafNodes = tmp.value.filter(n => !n.leaf)
    if (nonLeafNodes.length <= 0) {
        return
    }
    // use heap
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

    tmp.value.splice(tmp.value.indexOf(maxE), 1);
    const newNodes = maxE.Split(context!);
    for (let i = 0; i < newNodes.length; i++) {
        const node = newNodes[i];
        // if (!node.leaf) {
            // tmp.value.push(node);
        // }
        tmp.value.push(node);
    }
    errorVal.value = maxE.error.toPrecision(5);
}

function redraw(highlight: boolean = false) {
    if (svg === null) {
        return
    }
    const rect: Selection<SVGRectElement, QuadNode, SVGGElement, any> = svg
        .selectAll<SVGRectElement, QuadNode>('rect')
        .data(tmp.value, n => n.id);
    rect.exit().remove();
    rect.enter()
        .append('rect')
        .attr('x', n => { return n.x + 0.25 })
        .attr('y', n => { return n.y + 0.25 })
        .attr('rx', props.param?.roundedCorner || 0)
        .attr('width', n => { return n.w - 0.5 })
        .attr('height', n => { return n.h - 0.5 })
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
        const saveCanvas = async (canvas: HTMLCanvasElement, suggestedName: string) => {
            // Feature detection. The API needs to be supported
            // and the app not run in an iframe.
            const supportsFileSystemAccess = 'showSaveFilePicker' in window && (() => {
                try {
                    return window.self === window.top;
                } catch {
                    return false;
                }
            })();
            // If the File System Access API is supported…
            if (supportsFileSystemAccess) {
                try {
                    // Show the file save dialog.
                    const handle = await window.showSaveFilePicker({
                        suggestedName,
                    });
                    // Write the blob to the file.
                    const writable = await handle.createWritable();
                    canvas.toBlob(async (blob) => {
                        await writable.write(blob!)
                        await writable.close();
                    })
                    return;
                } catch (err) {
                    // Fail silently if the user has simply canceled the dialog.
                    if (err instanceof DOMException && err.name !== 'AbortError') {
                        console.error(err.name, err.message);
                        return;
                    }
                    return;
                }
            }

            const link = document.createElement('a');
            link.download = suggestedName;
            document.body.append(link);
            link.href = canvas.toDataURL();
            link.click();
            link.remove();
        };
        saveCanvas(tmpCanvas, "image.png")
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

watch(() => props.imgSrc, (newImgSrc) => img.src = newImgSrc);

watch(() => props.param.backgroundColor, (newColor) => svg?.style('background-color', newColor));

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
