<template>
    <div>
        <div class="scroll-container flip">
            <div class="text-display flip" v-if="document !== null">
                <div v-for="(paragraph, parIndex) in document.source" v-bind:key="parIndex" draggable="false"
                     class="clearfix">
                    <p class="paragraph">
                    <span v-for="(sentence, senIndex) in paragraph" v-bind:key="senIndex">
                        <span
                            v-for="(word, wordIndex) in sentence" v-bind:key="wordIndex"

                            v-on:mousedown.left.stop="start(parIndex, senIndex, wordIndex, $event)"
                            v-on:mouseenter="dragged(parIndex, senIndex, wordIndex)"

                            :style="color(word)"
                            class="border-on-hover word non-selectable"
                        >
                            {{ word.text }}
                        </span>
                    </span>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import {actions} from "@/store";
import vue from 'vue';
// import VueContext from 'vue-context';
// the default styling relies on <li> elements and specific classes.
import 'vue-context/dist/css/vue-context.css';
import {blend} from "@/core/Colors.ts";
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

    const highlighted = "#98FB98";

    const updateToTrue = word => vue.set(word.attributes, highlighted, true);
    const updateToFalse = word => vue.set(word.attributes, highlighted, false);

    const isHighlighted = word => highlighted in word.attributes && word.attributes[highlighted];

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

    export default {
        name: "TextRenderer",
        props: ['documentId'],
        data: function() {
            return {
                path: [],
                dragging: false,
                highlighted,
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
        // components: {VueContext},
        computed: {
            document() {
                return this.$store.getters.source;
            },

            containers() {
                return this.$store.getters.containers;
            },
        },

        watch: {
            containers: {
                deep: true,
                handler(containers) {
                    if (this.document === null) return;

                    containers.forEach(container => {
                        const on = container.colorInfo.activeHover || container.colorInfo.activeClick;

                        [container.main, ...container.subcodes].forEach(code => {
                            code.texts.filter(t => t.document_id === this.documentId).forEach(text => {
                                this.document.walk(text.first_word, text.last_word, word => {
                                    vue.set(word.attributes, container.colorInfo.bg, on);
                                })
                            })
                        })
                    })
                },
            }
        },

        methods: {

            isHighlighted,

            deleteTexts(codedTexts) {
                this.$store.dispatch(actions.DISASSOCIATE_TEXT, {codedTexts});
            },

            color(word) {
                let backgroundColors = [];

                Object.keys(word.attributes).forEach(color =>  {
                    if (word.attributes[color]) {
                        backgroundColors.push(color)
                    }
                })

                if (backgroundColors.length > 0) {
                    const colors = blend(backgroundColors);
                    return {backgroundColor: `rgba(${colors[0]}, ${colors[1]}, ${colors[2]}, .7)`};
                } else {
                    return {};
                }
            },

            start: function(paragraph, sentence, word, e) {
                this.dragTool.shift = e.shiftKey;

                const coord = {paragraph, sentence, word};

                if (isHighlighted(this.document.word(coord))) {
                    this.pickupStart(coord, e);
                } else {
                    this.dragStart(coord);
                }
            },

            dragStart: function(coord) {
                this.dragging = true;
                this.dragTool.anchor = coord;
                this.dragTool.last = coord;

                document.addEventListener("mouseup", this.dragStop);
            },

            pickupStart: function(coord, e) {
                const [words, anchor, last] = this.document.search(coord, isHighlighted);

                const region = {
                    anchor, last,
                    text: words.map(w => w.text).join(" "),
                    documentId: this.documentId
                }

                let div = createDiv(e.clientX, e.clientY);
                div.innerText = `"${region.text}"`;

                document.body.appendChild(div);

                let move = (e) => {
                    div.style.left = (1 + e.clientX) + "px";
                    div.style.top = e.clientY + "px";
                }

                let remove = (e) => {
                    div.remove()
                    document.removeEventListener("mouseup", remove);
                    document.removeEventListener("mousemove", move);

                    let textDropEvent = new CustomEvent("text-drop", {
                        bubbles: true, cancelable: true,
                        detail: {
                            data: {words: region},
                            callback: () => {
                                console.log("callback called");
                            }
                        }
                    })

                    e.target.dispatchEvent(textDropEvent);
                }

                document.addEventListener("mousemove", move);
                document.addEventListener("mouseup", remove);
            },

            dragStop: function() {
                this.dragging = false;
                document.removeEventListener("mouseup", this.dragStop)
            },

            dragged: function(paragraph, sentence, word) {
                if (!this.dragging) {
                    return;
                }

                const newCoord = {paragraph, sentence, word};

                // this walk handles the case where a user drags back towards the anchor
                this.document.walk(this.dragTool.last, this.dragTool.anchor, updateToFalse);

                // update to correct highlight
                const update = this.dragTool.shift ? updateToFalse : updateToTrue;
                this.document.walk(this.dragTool.anchor, newCoord, update);

                this.dragTool.last = newCoord;
            },
        }
    }
</script>

<style scoped>
    .border-on-hover:hover {
        box-shadow: inset 0 0 10px indianred;
    }

    /*noinspection CssUnusedSymbol*/
    .active {
        background-color: #98fb98;
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