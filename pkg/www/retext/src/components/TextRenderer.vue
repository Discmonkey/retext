<template>
    <div class="large">
        <p class="paragraph" v-for="(paragraph, parIndex) in text.Paragraphs" v-bind:key="parIndex">
            <span v-for="(sentence, senIndex) in paragraph.Sentences" v-bind:key="senIndex">
                <span
                        v-for="(word, wordIndex) in sentence.Parts" v-bind:key="wordIndex"
                        v-on:mousedown.stop="start(parIndex, senIndex, wordIndex, $event)"
                        v-on:mouseenter="dragged(parIndex, senIndex, wordIndex)"
                        v-bind:class="{active: word.Selected}"
                        class="border-on-hover word non-selectable"
                >
                    {{word.Text}}
                </span>
            </span>
        </p>
    </div>
</template>

<script>
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
    import Vue from "vue"

    let createDiv = (x, y) => {
        let div = document.createElement("div")
        div.style.maxWidth = "400px";
        // div.style.maxHeight = "200px";
        // div.style.borderRadius = "10px";
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

        return div;
    };

    export default {
        name: "TextRenderer",
        props: {text: TextType, channel: Object},
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
        methods: {
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
                console.log(e);

                let selectedWords = [];
                let currentWord = this.words(paragraph, sentence)[word];
                let indices = [paragraph, sentence, word];

                while (currentWord.Selected) {
                    selectedWords.push(currentWord.Text);
                    indices = this.previous(indices[0], indices[1], indices[2]);

                    if (indices[3] === 0) {
                        break;
                    }
                    currentWord = this.words(indices[0], indices[1])[indices[2]];
                }

                selectedWords.reverse();

                indices = this.next(paragraph, sentence, word);
                currentWord = this.words(indices[0], indices[1])[indices[2]];

                while (currentWord.Selected) {
                    selectedWords.push(currentWord.Text);
                    indices = this.next(indices[0], indices[1], indices[2]);

                    if (indices[3] === 0) {
                        break;
                    }
                    currentWord = this.words(indices[0], indices[1])[indices[2]];
                }

                let div = createDiv(e.clientX, e.clientY);
                div.innerText = "\"" +  selectedWords.reduce((acc, val) => {return acc + " " + val;}) + "\"";

                document.body.appendChild(div);

                let move = (e) => {
                    div.style.left = (1 + e.clientX) + "px";
                    div.style.top = (e.clientY) + "px";
                }

                let remove = () => {
                    div.remove()
                    document.removeEventListener("mouseup", remove);
                    document.removeEventListener("mousemove", move);

                    let words = JSON.parse(JSON.stringify(this.dragTool));
                    words.documentID = this.$parent.selected;
                    words.text = selectedWords.join(" ");
                    // this.$emit("dragStop", words);

                    this.channel.two(words);
                }

                document.addEventListener("mousemove", move);
                document.addEventListener("mouseup", remove);

            },

            // updates the view
            updateWord: function(paragraph, sentence, word) {
                Vue.set(this.words(paragraph, sentence), word, this.words(paragraph, sentence)[word]);
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

                // if were not in the first sentence of a paragraph, we just need to grab the last word of the
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

    .active {
        background-color: palegreen;
    }

    .non-selectable {
        -moz-user-select: none;
        -webkit-user-select: none;
        -ms-user-select: none;
    }

    .word {
        padding-top: .25em;
        padding-bottom: .25em;
    }
</style>