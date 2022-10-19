// move to semantic ui or something
const p = {
    leafSize: 12,
    leafSizeMin: 4,
    leafSizeMax: 64,
    leafSizeStep: 4,
    roundedCorner: 0,
    roundedCornerMin: 0,
    roundedCornerMax: 50,
    errorThreshold: 420,
    errorThresholdMin: 50,
    errorThresholdMax: 2500,
    errorThresholdStep: 50,
    updateViewFreq: 200,
    updateModelFreq: 1,
}
Object.seal(p);

let canvas = null;
let context = null;
let img = new Image();
let svg = null;
let quads = [];
let updateModelTimer = null;
let updateViewTimer = null;
let id = 0;
let running = true;
let errorInfo = '0';

function start() {
    running = true;
    updateView();
    updateModel();
}

function stop() {
    running = false;
}

function restart() {
    running = true;
    reset();
    updateView();
    updateModel();
}

function reset() {
    const ratio = img.width / img.height;
    let width = img.width;
    let height = img.height;
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
    canvas.width = width;
    canvas.height = height;
    context.drawImage(img, 0, 0, width, height);

    quads = [];
    quads.push(createQuadNode(0, 0, width, height, '#000000'));
    redraw();
}

function rgbToHex(r, g, b) {
    return '#' + ((1 << 24) + (r << 16) + (g << 8) + (b)).toString(16).slice(1);
}

function se(c1, c2) {
    const rErr = c1.r - c2.r;
    const gErr = c1.g - c2.g;
    const bErr = c1.b - c2.b;
    return 0.2989 * (rErr*rErr) + 0.5870 * (gErr*gErr) + 0.1140 * (bErr*bErr);
}

function calcAverageColor(imageData) {
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

function calcColorMSE(imageData, {r, g, b}) {
    let mse = 0;
    for (let i = 0; i < imageData.length; i+=4) {
        mse += se(
            {r: r, g: g, b: b},
            {r: imageData[i], g: imageData[i+1], b: imageData[i+2]}
        )
    }
    return Math.sqrt(mse);
}

function createQuadNode(x, y, width, height, prevColor) {
    id++;
    const data = context.getImageData(x, y, width, height).data;
    const {avgR, avgG, avgB} = calcAverageColor(data);
    const mse = calcColorMSE(data, {r: avgR, g: avgG, b: avgB})
    const colorHex = rgbToHex(avgR, avgG, avgB);

    return {
        id: id,
        error: mse,
        color: colorHex,
        prevColor: prevColor,
        x: x, y: y, width: width, height: height,
        leaf: width < p.leafSize || height < p.leafSize,
    };
}

// needed so d3 know which id to remove
function key(n) {
    return n.id;
}

function redraw(highlight) {
    const iterations = (quads.length - 1) / 3;
    d3.select('#info')
        .text(
            'Iterations: ' + iterations +
            ' - Shapes: ' + quads.length +
            ' - Error: ' + errorInfo);

    // d3 join data to rect
    rect = svg.selectAll('rect').data(quads, key);
    // remove old
    rect.exit().remove();
    // create new one
    rect.enter()
        .append('rect')
        .attr('x', (n) => {return n.x+0.25})
        .attr('y', (n) => {return n.y+0.25})
        .attr('rx', p.roundedCorner)
        .attr('width', (n) => {return n.width-0.5})
        .attr('height', (n) => {return n.height-0.5})
        .attr('fill', (n) => {return highlight? '#ffffff': n.prevColor})
        .transition().duration(500).styleTween('fill', (n) => {
            return d3.interpolate(highlight? '#ffffff': n.prevColor, n.color);
        });
}

function split(node) {
    if (node.leaf) {
        return
    }
    const idx = quads.indexOf(node);
    quads.splice(idx, 1);

    const halfW = node.width / 2;
    const halfH = node.height / 2;
    quads.push(
        createQuadNode(node.x, node.y, halfW, halfH, node.color),
        createQuadNode(node.x + halfW, node.y, halfW, halfH, node.color),
        createQuadNode(node.x, node.y + halfH, halfW, halfH, node.color),
        createQuadNode(node.x + halfW, node.y + halfH, halfW, halfH, node.color)
    )
}

function step() {
    // can use quadtree instead of filtering array
    const nonLeafNodes = quads.filter((n) => {return !n.leaf})
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
    if (maxE.error < p.errorThreshold) {
        running = false;
        return
    }
    // console.log(maxE.error);
    split(maxE);
    errorInfo = maxE.error.toPrecision(5).toString();
}

function updateModel() {
    clearTimeout(updateModelTimer);
    if (running) {
        step();
        updateModelTimer = setTimeout(updateModel, p.updateModelFreq);
    }
}

function updateView() {
    clearTimeout(updateViewTimer);
    if (running) {
        redraw();
        updateViewTimer = setTimeout(updateView, p.updateViewFreq);
    }
}

function initSlider() {
    $('.ui.slider').slider();

    const roundedCornerInput = $('#rounded-corner-slider-input');
    roundedCornerInput.val(p.roundedCorner +'px');
    $('#rounded-corner-slider').slider({
        min: p.roundedCornerMin,
        max: p.roundedCornerMax,
        onMove: (d) => {
            roundedCornerInput.val(d + 'px');
        },
        onChange: (d) => {
            p.roundedCorner = d;
            restart();
        }
    });

    const leafSizeInput = $('#leaf-size-slider-input');
    leafSizeInput.val(p.leafSize +'px');
    $('#leaf-size-slider').slider({
        min: p.leafSizeMin,
        max: p.leafSizeMax,
        step: p.leafSizeStep,
        start: p.leafSize,
        onMove: (d) => {
            leafSizeInput.val(d + 'px');
        },
        onChange: (d) => {
            p.leafSize = d;
            restart();
        }
    });

    const errorThresholdInput = $('#error-threshold-slider-input');
    errorThresholdInput.val(p.errorThreshold);
    $('#error-threshold-slider').slider({
        min: p.errorThresholdMin,
        max: p.errorThresholdMax,
        step: p.errorThresholdStep,
        start: p.errorThreshold,
        labelDistance: 2000,
        onMove: (d) => {
            errorThresholdInput.val(d);
        },
        onChange: (d) => {
            p.errorThreshold = d;
            restart();
        }
    });
}

function initColorPicker() {
    // iro.js
    const width = document.getElementById('picker').parentElement.clientWidth - 40;
    const colorPicker = new iro.ColorPicker('#picker', {
        width: 260,
        layoutDirection: 'horizontal',
    });
    colorPicker.on('input:end', (color) => {
        svg.style('background-color', color.hexString);
    });
}

function handleImageUpload(file) {
    const imageTypeRe = /image.*/;
    if (!file.type.match(imageTypeRe)) {
        alert('please choose image file');
        return;
    }

    const reader = new FileReader();
    reader.onload = readerEvent => {
        img.src = readerEvent.target.result;
    }
    // use readAsDataURL for now
    // https://stackoverflow.com/a/31743665
    reader.readAsDataURL(file)
}

function initImageUploader() {
    document.getElementById('image-upload-input').onchange = e => {
        const file = e.target.files[0];
        handleImageUpload(file);
    };
}

function initDragAndDrop() {
    const dropArea = document.getElementById('main');
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(
            eventName,
            e => {
                e.preventDefault();
                e.stopPropagation();
            },
            false,
        );
    });

    ['dragenter', 'dragover'].forEach(eventName => {
        dropArea.addEventListener(eventName, _ => dropArea.classList.add('highlight'), false)
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropArea.addEventListener(eventName, _ => dropArea.classList.remove('highlight'), false)
    });

    dropArea.addEventListener('drop', e => {
        // if a file
        const files = e.dataTransfer.files
        if (files.length > 0) {
            handleImageUpload(files[0]);
            return;
        }
    });
}

function save() {
    const svgE = document.getElementById('target');
    const {width: w, height: h} = svgE.getBBox();

    // clone svg and turn to blobURL
    const cloneE = svgE.cloneNode(true);
    cloneE.setAttribute('xmlns', xmlns="http://www.w3.org/2000/svg");  // needed otherwise chrome will not work
    const outerHTML = cloneE.outerHTML;
    const blob = new Blob([outerHTML], {type:'image/svg+xml;charset=utf-8'});
    const URL = window.URL || window.webkitURL || window;
    const blobURL = URL.createObjectURL(blob);

    // create new image, canvas, and draw blobURL to canvas
    let copyImg = new Image();
    copyImg.onload = () => {
        let tmpCanvas = document.createElement('canvas');
        tmpCanvas.width = w;
        tmpCanvas.height = h;

        let tmpContext = tmpCanvas.getContext('2d');
        tmpContext.drawImage(copyImg, 0, 0, w, h);

        // save canvas
        let png = tmpCanvas.toDataURL();
        let download = function(href, name){
            const link = document.createElement('a');
            link.download = name;
            // link.style.display = "none";
            document.body.append(link);
            link.href = href;
            link.click();
            link.remove();
        }
        download(png, "image.png");
    }
    copyImg.src = blobURL;
}

window.onload = () => {
    canvas = document.getElementById('source');
    context = canvas.getContext('2d');
    svg = d3.select('svg');

    img.src = './tmp.jpg';
    img.onload = () => {
        reset();
        updateView();
        updateModel();
    };

    initSlider();
    initColorPicker();
    initImageUploader();
    $('.sidebar').sidebar({
        closable: false,
        dimPage: false,
    });

    initDragAndDrop();
}
