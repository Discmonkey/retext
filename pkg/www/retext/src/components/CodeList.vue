<template>
    <div class="container-fluid pad-20 limited">
        <div class="row">
            <div class="col-12 text-right">
                <button class="btn btn-primary bold add-height" id="add-button"
                        @text-drop="textDrop(defaultNewCode, $event.detail, $event)"
                        @click="createCode(codes, null)">
                    Add a New Code
                </button>
            </div>
        </div>

        <draggable v-model="codes" group="people" @start="drag=true" @end="drag=false">
            <div v-for="code in codes" :key="code.main.id" @text-drop="textDrop(code, $event.detail, $event)">
                <div class="spacer">
                    <div :style="style(code.main.id)">
                        <code-drop-zone :code="code.main">
                            <div class="top-container margined">
                                <h5 class="code-title">{{code.main.name}}</h5>

                                <div class="btn btn-primary float-right no-events just-number self-right">
                                    {{getTextsLength(code)}}
                                </div>

                                <button class="btn btn-primary float-right" @click="createCode(code.subcodes, code.main.id)">
                                    <i class="fa fa-plus"></i>
                                </button>
                            </div>
                        </code-drop-zone>

                        <div v-for="subCode in code.subcodes" v-bind:key="subCode.id" class="subcode margin-top">
                            <code-drop-zone :code="subCode">
                                <div class="row item">
                                    <div class="col-10 center-text pad-">
                                        {{subCode.name}}
                                    </div>

                                    <div class="col-2 center-text">
                                        {{subCode.texts.length}}
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
    import {getColor} from "@/core/Colors";
    // eslint-disable-next-line no-unused-vars
    let codeTypes = {
        codes: [{
            id: 0,
            name: "",
            texts: [{
                documentID: "",
                text: ""
            }]
        }]
    }
    function prepareCode(mainCode) {
        for(let i = 0; i <= mainCode.subcodes.length; ++i) {
            if(mainCode.subcodes[i].id === mainCode.main) {
                mainCode.main = mainCode.subcodes.splice(i, 1)[0];
                break;
            }
        }
        return mainCode;
    }
    export default {
        name: 'codeList',
        components: {Draggable, CodeDropZone},
        data: () => {
            let newCode = {main: {name: "New", id: 0, texts: []}, subcodes: {}};
            return {
                codes: [],
                defaultNewCode: newCode,
            }
        },
        mounted() {
            this.axios.get("/code/list").then((res) => {
                let categories = res.data

                for(let c of categories) {
                    this.codes.push(prepareCode(c));
                }
            });
        },
        methods: {
            textDrop: function (parentCode, packet, e) {
                e.stopPropagation(); // stop the even

                let code = packet.data.code;

                // unless an error happens, this function will get called
                let associate = (code) => {
                    if (typeof code === "boolean" || !code)
                        return;
                    let c = code;
                    if (parentCode.main.id === 0) {
                        c = code.main;
                    }
                    this._actuallyAssociate(c, packet.data.words, packet.callback);
                };

                if (parentCode.main.id === 0) {
                    this.createCode(this.codes, null).then(
                        associate
                    )
                } else if (!code) {
                    // dropped on the code-wrapper but not in a designated drop-zone.
                    // TODO: make the whole code-wrapper a drop zone? (see template above)
                    return false;
                } else if (code.id === 0) {
                    this.createCode(parentCode.subcodes, parentCode.main.id).then(
                        associate
                    )
                } else {
                    associate(code);
                }
            },

            createCode: function (codes, parentCodeID) {
                let name = prompt("Name of new code?");
                if (name === null) {
                    // prompt was cancelled
                    return Promise.reject(true).then(() => {
                    }, () => {
                    });
                }

                return this.axios.post("/code/create", {
                    code: name,
                    parentCodeID: parentCodeID
                }).then(function (res) {
                    let newCode = res.data
                    if (!parentCodeID) {
                        newCode = prepareCode(newCode);
                    }
                    codes.splice(0, 0, newCode);
                    return newCode;
                }, function () {
                    // todo: alert the user of failure
                    return false;
                });
            },
            _actuallyAssociate: function (code, words, callback) {
                this.axios.post("/code/associate", {
                    key: words.documentID,
                    codeID: code.id,
                    text: words.text
                }).then(() => {
                    code.texts.push(words);
                    callback();
                    // todo: "success" toast or something
                }, () => {
                    // todo: "an error occurred" toast or something
                });
            },
            getTextsLength(code) {
                let length = code.main.texts.length;

                if (code.subcodes) {
                    for (let subCode of code.subcodes) {
                        length += subCode.texts.length;
                    }
                }

                return length;
            },


            style(idx) {
                return `border-radius: 10px; border: 3px solid ${getColor(idx)}; padding: 20px;`;
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
        padding: 3px;
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