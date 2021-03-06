<template>
    <div>
        <div class="scroll-container flip">
            <div class="text-display flip" v-if="document !== null">
                    <p class="paragraph" v-for="(paragraph, parIndex) in document.source" v-bind:key="parIndex" draggable="false">
                    <span v-for="(sentence, senIndex) in paragraph" v-bind:key="senIndex">
                        <span
                            v-for="(word, wordIndex) in sentence" v-bind:key="wordIndex"

                            v-on:mousedown.left.stop="start(parIndex, senIndex, wordIndex, $event)"
                            v-on:mouseenter="dragged(parIndex, senIndex, wordIndex)"
                            @contextmenu.prevent="chooseMenu($event, parIndex, senIndex, wordIndex)"

                            :style="color(word)"
                            class="border-on-hover word non-selectable"
                        >
                            {{ word.text }}
                        </span>
                    </span>
                    </p>
            </div>
        </div>

        <vue-context ref="deleteMenu" v-slot="{data: activeTexts}">
            <template v-if="activeTexts && activeTexts.length" class="delete-menu">
                <li><h5 class="header">Delete</h5></li>
                <li v-for="t of activeTexts" :key="t.id" class="pointer">
                    <a @click="deleteText(t.id)"
                       :title="t.text"
                       class="delete-menu-active-text"
                       :style="getContainerColor(t.code_id)"
                    >
                        {{ t.text }}
                    </a>
                </li>
            </template>
        </vue-context>
    </div>
</template>

<script>
import {actions} from "@/store";
import vue from 'vue';
import VueContext from 'vue-context'
// the default styling relies on <li> elements and specific classes.
import 'vue-context/dist/css/vue-context.css';
import {blend} from "@/core/Colors.ts";

const highlighted = "#98FB98";

    const updateToTrue = word => vue.set(word.attributes.colors, highlighted, true);
    const updateToFalse = word => vue.set(word.attributes.colors, highlighted, false);

    const isHighlighted = word => highlighted in word.attributes.colors && word.attributes.colors[highlighted];

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
        components: {VueContext},
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
                                    vue.set(word.attributes.colors, container.colorInfo.bg, on);
                                    vue.set(word.attributes.textIds, text.id, on);
                                })
                            })
                        })
                    })
                },
            }
        },

        methods: {

            isHighlighted,

            chooseMenu(e, paragraph, sentence, word) {
                let allIds = this.document.word({paragraph, sentence, word}).attributes.textIds;
                let activeTexts = Object.keys(allIds).filter(id => allIds[id]).map(id => {
                    return this.$store.getters.idToText[id];
                });

                if(activeTexts.length) {
                    this.$refs.deleteMenu.open(e, activeTexts);
                }
            },

            getContainerColor(codeId) {
                const code = this.$store.getters.idToCode[codeId];
                const container = this.$store.getters.idToContainer[code.container];
                return {borderColor: container.colorInfo.bg};
            },

            async deleteText(textId) {
                const text = this.$store.getters.idToText[textId];


                await this.$store.dispatch(actions.code.DELETE_TEXT, {textId, codeId: text.code_id});

                const code = this.$store.getters.idToCode[text.code_id];
                const container = this.$store.getters.idToContainer[code.container];

                this.document.walk(text.first_word, text.last_word, word => {
                    vue.set(word.attributes.colors, container.colorInfo.bg, false);
                    vue.set(word.attributes.textIds, text.id, false);
                })
            },

            color(word) {
                let backgroundColors = [];

                Object.keys(word.attributes.colors).forEach(color =>  {
                    if (word.attributes.colors[color]) {
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
                const [words, first_word, last_word] = this.document.search(coord, isHighlighted);

                const documentText = {
                    first_word, last_word,
                    text: words.map(w => w.text).join(" "),
                    document_id: this.documentId
                }

                let div = createDiv(e.clientX, e.clientY);
                div.innerText = `"${documentText.text}"`;

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
                            data: {words: documentText},
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

    .pointer {
        cursor: pointer;
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

    .v-context {
        box-shadow: 0 2px 2px 0 rgba(0, 0, 0, 0.65),0 3px 1px -2px rgba(0, 0, 0, 0.65),0 1px 5px 0 rgba(0, 0, 0, 0.65);
    }

    .v-context .header {
        text-align: center;
    }

    .delete-menu-active-text {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 150px;
        border: 1px solid;
        border-radius: .25em;
        margin: 5px;
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