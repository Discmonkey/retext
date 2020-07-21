// a list of colors that are pleasing to the eye (taken from some website)
let allColors = ["#FFDC00", "#F012BE", "#aaa", "#FF4136", "#2ECC40", "#111", "#7FDBFF", "#001F3F", "#FF851B", "#0074D9", "#85144B", "#3D9970", "#39CCCC", "#01FF70", "#ddd", "#B10DC9", "#fff"];

/**
 * Get a unique color based on the index you pass in. The color returned for a
 *  particular index will be the same, even if the page is reloaded.
 *
 * Uses an array of pre-defined colors taken from colors.js for the first several
 *  colors and then starts generating colors for anything past the pre-defined set.
 *
 * @param index the color returned
 * @returns {string} a hex color code (eg #007fbb)
 */
export function getColor(index) {
    index = parseInt(index);
    while(index >= allColors.length) {
        allColors.push(generateColor(allColors.length));
    }

    return allColors[index];
}

/**
 * takes 3 numbers and turns them into a hex color code
 * from https://stackoverflow.com/questions/5623838/rgb-to-hex-and-hex-to-rgb
 *
 * @param {Number} r
 * @param {Number} g
 * @param {Number} b
 * @returns {string}
 */
const rgbToHex = (r, g, b) => '#' + [r, g, b].map(x => x.toString(16).padStart(2, '0')).join('')

/**
 * Creates a color based on the values passed in
 * h, s, v values must be in [0..1[
 * from https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/
 *
 * @param h hue
 * @param s saturation
 * @param v value
 * @returns {string} a hex color code
 */
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

/**
 * generates a pleasing color using the "golden ratio conjugate". (see from link).
 * Passing in the same number for h will always return the same color.
 *
 * Note that the original code stored h as a global variable so, if the colors are weird,
 *  definitely check the "from" link(and in the wikipedia articles linked by the link)
 *  there's some crazy math going on that may have been broken by the current set up.
 *
 * from https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/
 *
 * @param h
 * @returns {string} returns a hex color code
 */
let goldenRatioConjugate = 0.618033988749895;
function generateColor(h) {
    h = goldenRatioConjugate * h + goldenRatioConjugate;
    h %= 1;
    return hsvToRgb(h, 0.7, 0.95);
}


/**
 * Calculates an "opposite color". For example, pass in a background color,
 *  get back a foreground color that should be easily visible on the passed-in bg.
 *
 * from https://stackoverflow.com/questions/35969656/how-can-i-generate-the-opposite-color-according-to-current-color/35970186#35970186
 *
 * @param hex string containing a hex color code (with or without #). handles 3 and 6 digit codes
 * @param bw bool, true if you want only either black or white as the returned foreground color
 * @returns {string} hex color string (always includes #)
 */
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
