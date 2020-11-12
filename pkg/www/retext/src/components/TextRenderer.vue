<template>
    <div>
        <div class="scroll-container flip">
            <div class="text-display flip">
                <div v-for="(paragraph, parIndex) in text.Paragraphs" v-bind:key="parIndex" draggable="false"
                     class="clearfix">
                    <p class="paragraph">
                    <span v-for="(sentence, senIndex) in paragraph.Sentences" v-bind:key="senIndex">
                        <span
                            v-for="(word, wordIndex) in sentence.Parts" v-bind:key="wordIndex"
                            v-on:mousedown.left.stop="start(parIndex, senIndex, wordIndex, $event)"
                            v-on:mouseenter="dragged(parIndex, senIndex, wordIndex)"
                            :class="{active: word.Selected}"
                            :style="wordStyle(parIndex, senIndex, wordIndex)"
                            @contextmenu.prevent="chooseMenu($event, parIndex, senIndex, wordIndex)"
                            class="border-on-hover word non-selectable"
                        >
                            {{ word.Text }}
                        </span>
                    </span>
                    </p>

                    <div class="done" v-if="wordCoordTextMap[parIndex]">
                        <div
                            v-for="(containerId) in wordCoordTextMap[parIndex].cIds" :key="containerId"
                            :style="paragraphStyle(containerId)"
                            @click="toggleColor(containerId)"
                        ></div>
                    </div>
                </div>
            </div>
        </div>
        <vue-context ref="deleteMenu" v-slot="{data: codedTexts}">
            <template v-if="codedTexts">
                <li><a @click="deleteTexts(codedTexts)">Delete</a></li>
            </template>
        </vue-context>
        <vue-context ref="associateMenu" v-slot="{data}">
            <template>
            <li style="margin-left:10px"><h5>Associate</h5></li>
            <li
                v-for="(container) in containers"
                :key="container.containerId"
                :class="{'v-context__sub': container.subcodes.length}"
            >
                <a @click.prevent="associateSelectedText(container.main.id, data)">{{ container.main.name }}</a>
                <ul class="v-context" v-if="container.subcodes.length">
                    <li v-for="(subcode) in container.subcodes" :key="subcode.id">
                        <a @click.prevent="associateSelectedText(container.main.id, data)">{{ subcode.name }}</a>
                    </li>
                </ul>
            </li>
            </template>
        </vue-context>
    </div>
</template>

<script>
import {actions} from "@/store";
import {mapGetters} from "vuex";
import VueContext from 'vue-context';
// the default styling relies on <li> elements and specific classes.
import 'vue-context/dist/css/vue-context.css';
// eslint-disable-next-line no-unused-vars
    let TextType = {
        Paragraphs: [{
            Sentences: {
                Parts: [{
                    Text: "",
                    Selected: false
                }]
            }
        }]
    };

    let createDiv = (x, y) => {
        let div = document.createElement("div")
        div.style.maxWidth = "400px";
        div.style.zIndex = "9999";

        // to give floating effect
        div.style.boxShadow = "0px 0px 43px -8px rgba(52,179,128,1)"

        // otherwise its see through
        div.style.backgroundColor = "rgba(255, 255, 255, 1)"

        // #styling
        div.style.padding = "10px"

        // if its a lot of text, let's just cut it off
        div.style.overflowY = "hidden"

        // so that it follows the cursor
        div.style.position = "absolute";
        div.style.left = x + "px";
        div.style.top = y + "px"
        div.style.pointerEvents = "none";

        return div;
    };

    function mapCodifiedTexts(v, cc, code, container) {
        if(!code || !code.texts) {
            return;
        }

        let cid = container.containerId
        for (let text of code.texts) {

            let [paragraph, sentence, word] = [text.first_word.paragraph, text.first_word.sentence, text.first_word.word];
            // eslint-disable-next-line no-constant-condition
            while(true) {
                // check if p, s, w are keys; if not, add
                if(!(paragraph in cc)) {
                    cc[paragraph] = {
                        cIds: new Set(),
                        s: {}
                    };
                }
                if(!(sentence in cc[paragraph].s)) {
                    cc[paragraph].s[sentence] = {
                        w: {}
                    };
                }
                if(!(word in cc[paragraph].s[sentence].w)) {
                    cc[paragraph].s[sentence].w[word] = {
                        cIds: new Set(),
                        texts: {},
                    };
                }
                // check if con is already in [p].cons; if not, add
                cc[paragraph].cIds.add(cid);
                // check if con is already in [p,s,w].cons; if not, add
                cc[paragraph].s[sentence].w[word].cIds.add(cid);
                // add text to [p, s, w].texts
                if(!(cid in cc[paragraph].s[sentence].w[word].texts)) {
                    cc[paragraph].s[sentence].w[word].texts[cid] = {};
                }
                // add a "code layer" so we can use it if we disassociate text
                if(!(code.id in cc[paragraph].s[sentence].w[word].texts[cid])) {
                    cc[paragraph].s[sentence].w[word].texts[cid][code.id] = [];
                }
                cc[paragraph].s[sentence].w[word].texts[cid][code.id].push(text.id);
                // check if we're at the lastCoord; if yes, break. if not, move to next coord
                if(paragraph === text.last_word.paragraph
                    && sentence === text.last_word.sentence
                    && word === text.last_word.word) {
                    break;
                }
                [paragraph, sentence, word] = v.next(paragraph, sentence, word);
            }
        }
    }

/**
 * Given p, s, and w, finds all adjacent, selected words and returns their coordinates
 *  of the first and last words, and a string containing all the text between (inclusive)
 *
 * @param tr a TextRenderer instance
 * @param p  paragraph index
 * @param s  sentence index
 * @param w  word index
 * @returns {{anchor, last, text}}
 */
function getSelectedRegion(tr, p, s, w) {
        let regionInfo = {};

        let selectedWords = [];
        let currentWord = tr.words(p, s)[w];
        let indices = [p, s, w];
        let currentIndices = indices;

        while (currentWord.Selected) {
            selectedWords.push(currentWord.Text);
            currentIndices = indices;
            indices = tr.previous(indices[0], indices[1], indices[2]);

            if (indices[3] === 0) {
                break;
            }
            currentWord = tr.words(indices[0], indices[1])[indices[2]];
        }

        selectedWords.reverse();
        regionInfo.anchor = {
            paragraph: currentIndices[0],
            sentence: currentIndices[1],
            word: currentIndices[2],
        }

        currentWord = tr.words(indices[0], indices[1])[indices[2]];
        indices = tr.next(p, s, w);
        currentIndices = indices;

        while (currentWord.Selected) {
            selectedWords.push(currentWord.Text);
            currentIndices = indices;
            indices = tr.next(indices[0], indices[1], indices[2]);

            if (indices[3] === 0) {
                break;
            }
            currentWord = tr.words(indices[0], indices[1])[indices[2]];
        }

        regionInfo.last = {
            paragraph: currentIndices[0],
            sentence: currentIndices[1],
            word: currentIndices[2],
        }
        regionInfo.text = selectedWords.join(" ");

        return regionInfo;
    }

    export default {
        name: "TextRenderer",
        props: ["text", "documentId"],
        data: function() {
            return {
                path: [],
                dragging: false,

                dragTool: {
                    anchor: {
                        paragraph: -1,
                        sentence: -1,
                        word: -1,
                    },

                    last: {
                        paragraph: -1,
                        sentence: -1,
                        word: -1
                    },

                    shift: false,
                }
            }
        },
        components: {VueContext},
        computed: {
            wordCoordTextMap() {
                let containers = this.containers;
                let coordTextMap = {};

                for(let container of containers) {
                    mapCodifiedTexts(this, coordTextMap, container.main, container);

                    for(let subcode of container.subcodes) {
                        mapCodifiedTexts(this, coordTextMap, subcode, container);
                    }
                }

                return coordTextMap;
            },
            activeContainerId() {
                let containers = this.containers;

                for(let c of containers) {
                    if(c.colorInfo.active) {
                        return c.containerId
                    }
                }

                return false;
            },
            ...mapGetters(["containers", "idToContainer", "idToCode"]),
        },
        methods: {
            chooseMenu(e, p, s, w) {
                if(this.activeContainerId) {
                    // check if the word is part of an excerpt in the active container
                    if(this.activeContainerId && this.wordCoordTextMap[p].s[s].w[w].texts[this.activeContainerId]) {
                        if(
                            p in this.wordCoordTextMap &&
                            s in this.wordCoordTextMap[p].s &&
                            w in this.wordCoordTextMap[p].s[s].w &&
                            "texts" in this.wordCoordTextMap[p].s[s].w[w] &&
                            this.activeContainerId in this.wordCoordTextMap[p].s[s].w[w].texts
                        ) {
                            let codedTexts = this.wordCoordTextMap[p].s[s].w[w].texts[this.activeContainerId];
                            this.$refs.deleteMenu.open(e, codedTexts);
                        }
                    }
                } else if(this.text.Paragraphs[p].Sentences[s].Parts[w].Selected) {
                    this.$refs.associateMenu.open(e, {data: {p, s, w}});
                }
            },

            deleteTexts(codedTexts) {
                this.$store.dispatch(actions.DISASSOCIATE_TEXT, {codedTexts});
            },

            associateSelectedText(codeId, {p, s, w}) {
                let words = getSelectedRegion(this, p, s, w);

                this.$store.dispatch(actions.ASSOCIATE_TEXT, {codeId, words}).then(() => {
                    // todo: "success" toast or something
                }, () => {
                    // todo: "an error occurred" toast or something
                });
            },
            wordStyle(p, s, w) {
                if (!this.activeContainerId) {
                    return;
                }

                if(
                    p in this.wordCoordTextMap &&
                    s in this.wordCoordTextMap[p].s &&
                    w in this.wordCoordTextMap[p].s[s].w &&
                    "texts" in this.wordCoordTextMap[p].s[s].w[w]
                ) {
                    let texts = this.wordCoordTextMap[p].s[s].w[w].texts;
                    if (this.activeContainerId in texts) {
                        let ci = this.idToContainer[this.activeContainerId].colorInfo;
                        return {
                            backgroundColor: ci.bg + " !important",
                            color: ci.fg,
                        };
                    }
                }
            },

            paragraphStyle(containerId) {
                let c = this.idToContainer[containerId];
                return {
                    backgroundColor: c.colorInfo.bg,
                    borderRadius: "50% !important",
                    padding: ".5em !important",
                    marginLeft: "50% !important",
                    marginRight: "50% !important",
                };
            },

            toggleColor(containerId) {
                this.$store.dispatch(actions.SET_COLOR_ACTIVE, {containerId});
            },

            start: function(paragraph, sentence, word, e) {
                this.dragTool.shift = e.shiftKey;

                if (!this.words(paragraph, sentence)[word].Selected || this.dragTool.shift) {
                    this.dragStart(paragraph, sentence, word);
                } else {
                    this.pickupStart(paragraph, sentence, word, e);
                }
            },

            dragStart: function(paragraph, sentence, word) {
                this.dragging = true;
                this.dragTool.anchor.paragraph = paragraph;
                this.dragTool.anchor.sentence = sentence;
                this.dragTool.anchor.word = word;

                this.dragTool.last.paragraph = paragraph;
                this.dragTool.last.sentence = sentence;
                this.dragTool.last.word = word;

                this.words(paragraph, sentence)[word].Selected = !this.dragTool.shift;

                this.updateWord(paragraph, sentence, word);

                document.addEventListener("mouseup", this.dragStop);
            },

            pickupStart: function(paragraph, sentence, word, e) {
                let words = getSelectedRegion(this, paragraph, sentence, word);

                let div = createDiv(e.clientX, e.clientY);
                div.innerText = `"${words.text}"`;

                document.body.appendChild(div);

                let move = (e) => {
                    div.style.left = (1 + e.clientX) + "px";
                    div.style.top = e.clientY + "px";
                }
                const documentId = this.documentId;
                let remove = (e) => {
                    div.remove()
                    document.removeEventListener("mouseup", remove);
                    document.removeEventListener("mousemove", move);

                    words.documentId = documentId;

                    let textDropEvent = new CustomEvent("text-drop", {
                        bubbles: true, cancelable: true,
                        detail: {
                            data: {words},
                            callback: () => {
                                // todo: add code-specific color-class
                                console.log(`sample: ${JSON.stringify(words)}`);
                            }
                        }
                    })

                    e.target.dispatchEvent(textDropEvent);
                }

                document.addEventListener("mousemove", move);
                document.addEventListener("mouseup", remove);
            },

            // updates the view
            updateWord: function(paragraph, sentence, word) {
                this.$set(this.words(paragraph, sentence), word, this.words(paragraph, sentence)[word]);
            },

            dragStop: function() {
                this.dragging = false;
                document.removeEventListener("mouseup", this.dragStop)
            },


            // returns true if point set 2 is greater than equal to point set 1
            greaterEqual: function(i2, j2, k2, i, j, k) {
                if (i2 > i) {
                    return true;
                }

                if (i2 < i) {
                    return false;
                }

                if (j2 > j) {
                    return true;
                }

                if (j2 < j) {
                    return false;
                }

                return k2 >= k;
            },

            inside: function(i, j, k, iStart, jStart, kStart, iEnd, jEnd, kEnd) {
                // quick swap just in case
                if (this.greaterEqual(iStart, jStart, kStart, iEnd, jEnd, kEnd)) {
                    let a = iStart;
                    let b = jStart;
                    let c = kStart;

                    iStart = iEnd;
                    jStart = jEnd;
                    kStart = kEnd;

                    iEnd = a;
                    jEnd = b;
                    kEnd = c;
                }

                return this.greaterEqual(i, j, k, iStart, jStart, kStart) && this.greaterEqual(iEnd, jEnd, kEnd, i, j, k);
            },

            dragged: function(paragraph, sentence, word) {
                if (!this.dragging) {
                    return;
                }

                // so first off we want to have everything between our anchor and the current point highlighted
                // we want to have everything outside of this range, but in the range of lastPoint -> currentPoint
                // un-highlighted

                let lasti = this.dragTool.last.paragraph;
                let lastj = this.dragTool.last.sentence;
                let lastk = this.dragTool.last.word;

                let anchori = this.dragTool.anchor.paragraph;
                let anchorj = this.dragTool.anchor.sentence;
                let anchork = this.dragTool.anchor.word;

                this.dragTool.last.paragraph = paragraph;
                this.dragTool.last.sentence = sentence;
                this.dragTool.last.word = word;
                // two ranges
                // so we have some overlap issues
                // there are always a lot of edge cases. but let's do the easy thing.

                let updateToFalse = (i, j, k) => {
                    this.words(i, j)[k].Selected = false;
                    this.updateWord(i, j, k);
                }
                if (this.greaterEqual(paragraph, sentence, word,
                    lasti, lastj, lastk)) {
                    this.iterate(lasti, lastj, lastk, paragraph, sentence, word, updateToFalse);
                } else {
                    this.iterate(paragraph, sentence, word, lasti, lastj, lastk, updateToFalse);
                }

                let updateAgainstShift = (i, j, k) => {
                    this.words(i, j)[k].Selected = !this.dragTool.shift;
                    this.updateWord(i, j, k);
                }

                if (this.greaterEqual(paragraph, sentence, word, anchori, anchorj, anchork)) {
                    this.iterate(anchori, anchorj, anchork, paragraph, sentence, word, updateAgainstShift)
                } else {
                    this.iterate(paragraph, sentence, word, anchori, anchorj, anchork, updateAgainstShift)
                }
            },

            // returns the next i, j, k (paragraph, sentence, word) tuple
            next: function(i, j, k) {
                // last word in sentence condition
                if (this.words(i, j).length - 1 !== k) {
                    return [i, j, k + 1, 1];
                }

                // now only the last word in the sentence, but also the last sentence in the paragraph
                if (this.sentences(i).length - 1 !== j) {

                    // return the first word of the next sentence
                    return [i, j + 1, 0, 1];
                }

                // well now we're really unlucky,
                if (this.paragraphs().length -1 !== i) {
                    return [i + 1, 0, 0, 1]
                }

                // if we're in the last paragraph, last sentence, last word state, let's just return i, j, k
                return [i, j, k, 0]
            },

            paragraphs: function() {
                return this.text.Paragraphs;
            },

            sentences: function(i) {
                return this.paragraphs()[i].Sentences;
            },

            words: function(i, j) {
                return this.sentences(i)[j].Parts;
            },

            previous: function(i, j, k) {
                if (k !== 0) {
                    return [i, j, k - 1, 1];
                }

                // if we're not in the first sentence of a paragraph, we just need to grab the last word of the
                // previous sentence
                if (j !== 0) {

                    // return the first word of the next sentence
                    return [i, j - 1, this.words(i, j - 1).length - 1, 1];
                }

                // well now we're really unlucky, but so long as we also weren't in the first paragraph we should be okay
                if (i !== 0 ) {
                    let lastSentence = this.sentences(i - 1).length - 1;
                    return [i - 1,  lastSentence, this.words(i - 1, lastSentence).length - 1, 1]
                }

                return [i, j, k, 0]
            },

            iterate: function(i1, j1, k1, i2, j2, k2, callback) {

                // for each paragraph
                for (let i = i1; i <= i2; i++) {
                    let startSentenceIndex = 0;
                    let endSentenceIndex = this.sentences(i).length - 1;

                    if (i === i1) {
                        startSentenceIndex = j1;
                    }

                    if (i === i2) {
                        endSentenceIndex = j2;
                    }

                    for (let j = startSentenceIndex; j <= endSentenceIndex; j++) {
                        let startWordIndex = 0;
                        let endWordIndex = this.words(i, j).length - 1;

                        if (i === i1 && j === j1) {
                            startWordIndex = k1;
                        }

                        if (i === i2 && j === j2) {
                            endWordIndex = k2;
                        }

                        for (let k = startWordIndex; k <= endWordIndex; k++) {
                            callback(i, j, k);
                        }
                    }
                }
            }
        }
    }
</script>

<style scoped>
    .border-on-hover:hover {
        box-shadow: inset 0 0 10px indianred;
    }

    /*noinspection CssUnusedSymbol*/
    .active {
        background-color: palegreen;
    }

    .flip {
        transform: scaleX(-1);
    }

    .non-selectable {
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
    }

    .paragraph {
        width: 95%;
        float: left;
    }

    .done {
        width: 5%;
        float: right;
    }

    .text-display {
        padding: 20px;

    }

    .scroll-container {
        height: 80vh;
        overflow: auto;
    }

    .word {
        padding-top: .25em;
        padding-bottom: .25em;
    }

    /*
    used to give the parent div of .paragraph+.done some height.
    without the height, the "list" of container-color-drops starts overlapping(try it out)
    from: https://stackoverflow.com/questions/12540436/
    */
    .clearfix:after {
        content: ".";
        display: block;
        clear: both;
        visibility: hidden;
        line-height: 0;
        height: 0;
    }
    .clearfix {
        display: inline-block;
    }
</style>