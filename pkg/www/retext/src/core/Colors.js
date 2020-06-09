let allColors = ["#001F3F","#0074D9","#7FDBFF","#39CCCC","#3D9970","#2ECC40","#01FF70","#FFDC00","#FF851B","#FF4136","#F012BE","#B10DC9","#85144B","#fff","#aaa","#ddd","#111"];

export function getColor(index) {
    index = parseInt(index);
    while(index >= allColors.length) {
        allColors.push(createNewColor(allColors.length));
    }

    return allColors[index];
}

function createNewColor(index) {
    return generateColor(index);
}

const rgbToHex = (r, g, b) => '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('')
// https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/ for hsvToRgb() & generateColor()
// HSV values in [0..1[
// returns [r, g, b] values from 0 to 255
function hsvToRgb(h, s, v) {
    let h_i = Math.floor(h * 6);
    let f = h * 6 - h_i;
    let p = v * (1 - s);
    let q = v * (1 - f * s);
    let t = v * (1 - (1 - f) * s);
    let r, g, b;
    if (h_i === 0) {
        [r, g, b] = [q, v, p];
    }
    if (h_i === 1) {
        [r, g, b] = [p, v, t];
    }
    if (h_i === 2) {
        [r, g, b] = [p, q, v];
    }
    if (h_i === 3) {
        [r, g, b] = [v, t, p];
    }
    if (h_i === 4) {
        [r, g, b] = [t, p, v];
    }
    if (h_i === 5) {
        [r, g, b] = [v, p, q];
    }
    [r, g, b] = [Math.floor(r * 256), Math.floor(g * 256), Math.floor(b * 256)];
    return rgbToHex(r, g, b);
}

let goldenRatioConjugate = 0.618033988749895

export function generateColor(h) {
    h = h * goldenRatioConjugate + goldenRatioConjugate;
    h %= 1;
    return hsvToRgb(h, 0.7, 0.95);
}

// https://stackoverflow.com/a/35970186
export function invertColor(hex, bw) {
    if (hex.indexOf('#') === 0) {
        hex = hex.slice(1);
    }
    // convert 3-digit hex to 6-digits.
    if (hex.length === 3) {
        hex = hex[0] + hex[0] + hex[1] + hex[1] + hex[2] + hex[2];
    }
    if (hex.length !== 6) {
        throw new Error('Invalid HEX color.');
    }
    let r = parseInt(hex.slice(0, 2), 16),
        g = parseInt(hex.slice(2, 4), 16),
        b = parseInt(hex.slice(4, 6), 16);
    if (bw) {
        // http://stackoverflow.com/a/3943023/112731
        return (r * 0.299 + g * 0.587 + b * 0.114) > 186
            ? '#000000'
            : '#FFFFFF';
    }
    // invert color components
    r = (255 - r).toString(16);
    g = (255 - g).toString(16);
    b = (255 - b).toString(16);
    // pad each with zeros and return
    return "#" + r.padStart(2, '0') + g.padStart(2, '0') + b.padStart(2, '0');
}
