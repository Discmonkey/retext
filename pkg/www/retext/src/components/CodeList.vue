<template>
    <div class="container-fluid pad-20 limited">
        <div class="row">
            <div class="col-12 text-right">
                <button class="btn btn-primary bold add-height" id="add-button"
                        @click="createCode(false)">
                    Add a New Code
                </button>
            </div>
        </div>

        <draggable v-model="containers" group="people" @start="drag=true" @end="drag=false">
            <div v-for="container in containers"
                 :key="container.containerId"
                 @text-drop="textDrop(container.main, $event.detail, $event)"
                 @click="toggleColor(container)"
            >
                <div class="spacer">
                    <div :style="style(container)">
                        <code-drop-zone :code="container.main">
                            <div class="top-container margined">
                                <h5 class="code-title">{{container.main.name}}</h5>

                                <div class="btn btn-primary float-right no-events just-number self-right">
                                    {{getContainerTextsLength(container)}}
                                </div>

                                <button class="btn btn-primary float-right" @click="createCode(container.containerId)">
                                    <i class="fa fa-plus"></i>
                                </button>
                            </div>
                        </code-drop-zone>

                        <div v-for="subCode in container.subcodes" v-bind:key="subCode.id" class="subcode margin-top">
                            <code-drop-zone :code="subCode">
                                <div class="row item">
                                    <div class="col-10 center-text pad-3 rborder">
                                        {{subCode.name}}
                                    </div>

                                    <div class="col-2 center-text pad-3">
                                        {{subCode.texts == null ? 0 : subCode.texts.length}}
                                    </div>
                                </div>
                            </code-drop-zone>
                        </div>

                    </div>
                </div>
            </div>
        </draggable>
    </div>
</template>

<script>
import Draggable from 'vuedraggable';
import CodeDropZone from "@/components/CodeDropZone";
import {actions, getters, mutations} from "@/store"
import {mapGetters} from "vuex";
// eslint-disable-next-line no-unused-vars
    let codeTypes = {
        codes: [{
            id: 0,
            name: "",
            texts: [{
                documentId: "",
                text: ""
            }]
        }]
    }
    export default {
        name: 'codeList',
        components: {Draggable, CodeDropZone},
        computed: {
            [getters.CONTAINERS]: {
                get() {
                    return this.$store.getters[getters.CONTAINERS];
                },
                set(c) {
                    this.$store.commit(mutations.SET_CONTAINERS, c)
                },
            },
            ...mapGetters([getters.ID_TO_CONTAINER])
        },
        mounted() {
            this.$store.dispatch(actions.INIT_CONTAINERS)
        },
        methods: {
            toggleColor(container) {
                this.$store.dispatch(actions.SET_COLOR_ACTIVE, container);
            },
            textDrop: async function (code, packet, e) {
                e.stopPropagation(); // stop the event

                // unless an error happens, this function will get called
                let associate = (code) => {
                    if (typeof code === "boolean" || !code)
                        return;

                    this._actuallyAssociate(e.detail.data.code, packet.data.words, packet.callback);
                };

                associate(code);
            },

            createCode: async function (containerId) {
                let name = prompt("Name of new code?");
                if (name === null) {
                    // prompt was cancelled
                    return Promise.reject(true).then(() => {
                    }, () => {
                    });
                }

                if (!containerId) {
                    this.$store.dispatch(actions.CREATE_CONTAINER, {name}).then(() => {});
                } else {
                    this.$store.dispatch(actions.CREATE_CODE, {containerId, name}).then(() => {});
                }
            },

            _actuallyAssociate: function (code, words, callback) {
                this.$store.dispatch(actions.ASSOCIATE_TEXT, {codeId: code.id, words}).then(() => {
                    callback();
                    // todo: "success" toast or something
                }, () => {
                    // todo: "an error occurred" toast or something
                });
            },
            getCode(codeId) {
                return this.$store.getters[getters.GET_CODE](codeId)
            },
            getContainerTextsLength(container) {
                return this.$store.getters[getters.GET_TEXTS_LENGTH](container.containerId)
            },


            style(container) {
                return `border-radius: 10px; border: 3px solid ${container.colorInfo.bg}; padding: 20px;`;
            }
        }

    }
</script>

<style scoped>

    .bold {
        font-weight: bold;
    }

    #add-button {
        height: 50px;
    }

    .add-height {
        width: 200px
    }


    .pad-20 {
        padding: 20px;
    }

    .pad-3 {
        padding: 3px;
    }

    .rounded-series * {
        border: 1px solid grey;
        height: 100%;
        text-align: center;
    }

    .rounded-series > :first-child {
        text-align: left;
        border-bottom-left-radius: .25em;
        border-top-left-radius: .25em;
    }

    .rounded-series > :last-child {
        border-bottom-right-radius: .25em;
        border-top-right-radius: .25em;
    }

    .spacer {
        margin-top: 15px;
        margin-bottom: 5px;
    }

    .text-right {
        text-align: right;
    }


    .center-text {
        text-align: center;
    }
    .code-title {
        text-transform: capitalize;
        font-weight: bolder;
        margin: 0;
    }

    .just-number {
        background-color: white;
        color: black;
        border: none;
        font-weight: bolder;
    }

    .text-right {
        text-align: right;
    }

    .margin-top {
        margin-top: 20px;
    }

    .limited {
        height: 80vh;
        overflow-y: scroll;
        -ms-overflow-style: none;
        scrollbar-width: none;
    }

    .limited::-webkit-scrollbar {
        display: none;
    }


    .no-events {
        pointer-events:none;
    }

    .rborder {
        border-right: 1px solid gray;
    }

    .subcode {
        border-radius: 3px;
        color: gray;
        margin-bottom: 5px;
        border: 1px solid gray;
    }

    .top-container {
        display: flex;
        align-content: center;
        align-items: center;
    }

    .self-right {
        margin-left: auto;
    }
</style>