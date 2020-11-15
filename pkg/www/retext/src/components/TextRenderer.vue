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

                            :class="{active: isHighlighted(word)}"
                            class="border-on-hover word non-selectable"
                        >
                            {{ word.text }}
                        </span>
                    </span>
                    </p>

<!--                    <div class="done" v-if="wordCoordTextMap[parIndex]">-->
<!--                        <div-->
<!--                            v-for="(containerId) in wordCoordTextMap[parIndex].cIds" :key="containerId"-->
<!--                            :style="paragraphStyle(containerId)"-->
<!--                            @click="toggleColor(containerId)"-->
<!--                        ></div>-->
<!--                    </div>-->
                </div>
            </div>
        </div>
<!--        <vue-context ref="deleteMenu" v-slot="{data: codedTexts}">-->
<!--            <template v-if="codedTexts">-->
<!--                <li><a @click="deleteTexts(codedTexts)">Delete</a></li>-->
<!--            </template>-->
<!--        </vue-context>-->
<!--        <vue-context ref="associateMenu" v-slot="{data}">-->
<!--            <template>-->
<!--            <li style="margin-left:10px"><h5>Associate</h5></li>-->
<!--            <li-->
<!--                v-for="(container) in containers"-->
<!--                :key="container.containerId"-->
<!--                :class="{'v-context__sub': container.subcodes.length}"-->
<!--            >-->
<!--                <a @click.prevent="associateSelectedText(container.main.id, data)">{{ container.main.name }}</a>-->
<!--                <ul class="v-context" v-if="container.subcodes.length">-->
<!--                    <li v-for="(subcode) in container.subcodes" :key="subcode.id">-->
<!--                        <a @click.prevent="associateSelectedText(container.main.id, data)">{{ subcode.name }}</a>-->
<!--                    </li>-->
<!--                </ul>-->
<!--            </li>-->
<!--            </template>-->
<!--        </vue-context>-->
    </div>
</template>

<script>
import {actions} from "@/store";
import {mapGetters} from "vuex";
import vue from 'vue';
// import VueContext from 'vue-context';
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

    const highlighted = "highlighted";

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
                console.log(this.$store.getters.source);
                return this.$store.getters.source;
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

            isHighlighted,
            // chooseMenu(e, p, s, w) {
            //     if(this.activeContainerId) {
            //         // check if the word is part of an excerpt in the active container
            //         if(this.activeContainerId && this.wordCoordTextMap[p].s[s].w[w].texts[this.activeContainerId]) {
            //             if(
            //                 p in this.wordCoordTextMap &&
            //                 s in this.wordCoordTextMap[p].s &&
            //                 w in this.wordCoordTextMap[p].s[s].w &&
            //                 "texts" in this.wordCoordTextMap[p].s[s].w[w] &&
            //                 this.activeContainerId in this.wordCoordTextMap[p].s[s].w[w].texts
            //             ) {
            //                 let codedTexts = this.wordCoordTextMap[p].s[s].w[w].texts[this.activeContainerId];
            //                 this.$refs.deleteMenu.open(e, codedTexts);
            //             }
            //         }
            //     } else if(this.text.Paragraphs[p].Sentences[s].Parts[w].Selected) {
            //         this.$refs.associateMenu.open(e, {data: {p, s, w}});
            //     }
            // },

            deleteTexts(codedTexts) {
                this.$store.dispatch(actions.DISASSOCIATE_TEXT, {codedTexts});
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