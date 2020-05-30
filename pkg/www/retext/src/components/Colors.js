export class ColorPalette {
    static _allColors = [    /*.category:nth-child(17n+1){background-color:#001F3F}.category:nth-child(17n+2){background-color:#0074D9}.category:nth-child(17n+3){background-color:#7FDBFF}.category:nth-child(17n+4){background-color:#39CCCC}.category:nth-child(17n+5){background-color:#3D9970}.category:nth-child(17n+6){background-color:#2ECC40}.category:nth-child(17n+7){background-color:#01FF70}.category:nth-child(17n+8){background-color:#FFDC00}.category:nth-child(17n+9){background-color:#FF851B}.category:nth-child(17n+10){background-color:#FF4136}.category:nth-child(17n+11){background-color:#F012BE}.category:nth-child(17n+12){background-color:#B10DC9}.category:nth-child(17n+13){background-color:#85144B}.category:nth-child(17n+14){background-color:#fff}.category:nth-child(17n+15){background-color:#aaa}.category:nth-child(17n+16){background-color:#ddd}.category:nth-child(17n+17){background-color:#111}.navy{color:#001F3F}.blue{color:#0074D9}.aqua{color:#7FDBFF}.teal{color:#39CCCC}.olive{color:#3D9970}.green{color:#2ECC40}.lime{color:#01FF70}.yellow{color:#FFDC00}.orange{color:#FF851B}.red{color:#FF4136}.fuchsia{color:#F012BE}.purple{color:#B10DC9}.maroon{color:#85144B}.white{color:#fff}.silver{color:#ddd}.gray{color:#aaa}.black{color:#111}.border--navy{border-color:#001F3F}.border--blue{border-color:#0074D9}.border--aqua{border-color:#7FDBFF}.border--teal{border-color:#39CCCC}.border--olive{border-color:#3D9970}.border--green{border-color:#2ECC40}.border--lime{border-color:#01FF70}.border--yellow{border-color:#FFDC00}.border--orange{border-color:#FF851B}.border--red{border-color:#FF4136}.border--fuchsia{border-color:#F012BE}.border--purple{border-color:#B10DC9}.border--maroon{border-color:#85144B}.border--white{border-color:#fff}.border--gray{border-color:#aaa}.border--silver{border-color:#ddd}.border--black{border-color:#111}.fill-navy{fill:#001F3F}.fill-blue{fill:#0074D9}.fill-aqua{fill:#7FDBFF}.fill-teal{fill:#39CCCC}.fill-olive{fill:#3D9970}.fill-green{fill:#2ECC40}.fill-lime{fill:#01FF70}.fill-yellow{fill:#FFDC00}.fill-orange{fill:#FF851B}.fill-red{fill:#FF4136}.fill-fuchsia{fill:#F012BE}.fill-purple{fill:#B10DC9}.fill-maroon{fill:#85144B}.fill-white{fill:#fff}.fill-gray{fill:#aaa}.fill-silver{fill:#ddd}.fill-black{fill:#111}.stroke-navy{stroke:#001F3F}.stroke-blue{stroke:#0074D9}.stroke-aqua{stroke:#7FDBFF}.stroke-teal{stroke:#39CCCC}.stroke-olive{stroke:#3D9970}.stroke-green{stroke:#2ECC40}.stroke-lime{stroke:#01FF70}.stroke-yellow{stroke:#FFDC00}.stroke-orange{stroke:#FF851B}.stroke-red{stroke:#FF4136}.stroke-fuchsia{stroke:#F012BE}.stroke-purple{stroke:#B10DC9}.stroke-maroon{stroke:#85144B}.stroke-white{stroke:#fff}.stroke-gray{stroke:#aaa}.stroke-silver{stroke:#ddd}.stroke-black{stroke:#111}
    .new-category {
        background-color: #01FF70;
    }*/];
    static _lastUsedIndex = -1;

    getColor(index) {
        if(index !== undefined) {
            let i = Number(index)
            if(i && 0 <= i && i < ColorPalette._allColors.length) {
                return [index, ColorPalette._allColors[i]];
            } else {
                // todo: this is an error case; come up with a good error-return
                return [false, false];
            }
        }

        // get an unused index
        let newIndex = ++ColorPalette._lastUsedIndex;
        // check if there's a pre-defined color left
        if(newIndex < ColorPalette._allColors.length) {
            return [newIndex, ColorPalette._allColors[newIndex]];
        }
    }

    _createNewColor() {

    }
}
const rgbToHex = (r, g, b) => '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('')
// https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/ for hsv_to_rgb() & gen_html()
// HSV values in [0..1[
// returns [r, g, b] values from 0 to 255
function hsv_to_rgb(h, s, v) {
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
let golden_ratio_conjugate = 0.618033988749895
let h = Math.random() // use random start value
export function gen_html() {
  h += golden_ratio_conjugate
  h %= 1
  return hsv_to_rgb(h, 0.7, 0.95)
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
