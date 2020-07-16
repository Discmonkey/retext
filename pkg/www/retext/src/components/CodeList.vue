<template>
    <div class="container-fluid">
        <div class="row">
            <button class="btn btn-primary offset-md-5 col-md-7"
                    @text-drop="textDrop(defaultNewCode, $event.detail, $event)"
                    @click="createCode(codes, null)"
            >Add a New Code
            </button>

        </div>
        <div v-for="code in codes" :key="code.main.id"
             @text-drop="textDrop(code, $event.detail, $event)"
             class="col-md-12 code-wrapper">
            <div class="row">
                <div class="col-sm-7"><b>
                    <code-drop-zone :code="code.main">{{code.main.name}}</code-drop-zone>
                </b></div>
                <div class="col-sm-5 row">
                    <div class="col-md-4">
                        <code-drop-zone
                                :code="defaultNewCode.main"
                        >
                            <div
                                @click="createCode(code.subcodes, code.main.id)"
                                class="btn-primary parent-code-option">+</div>
                        </code-drop-zone>
                    </div>
                    <div class="col-md-4">
                        <div class="btn-primary parent-code-option">D</div>
                    </div>
                    <div class="col-md-4">
                        <div class="btn-primary parent-code-option">{{getTextsLength(code)}}</div>
                    </div>
                </div>
            </div>
            <div style="margin-left: 7%" class="row">
                <code-drop-zone v-for="subCode in code.subcodes" :key="subCode.id"
                                    class="row rounded-series subcode" style="width:100%"
                                    :code="subCode"
                >
                    <span class="col-md-9">{{subCode.name}}</span>
                    <span class="col-md-3">{{subCode.texts.length}}</span>
                </code-drop-zone>
            </div>
        </div>
    </div>
</template>

<script>
    // todo: add back colors
    import CodeDropZone from "./CodeDropZone";

    // eslint-disable-next-line no-unused-vars
    let codesTypes = {
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
        name: 'CodeList',
        components: {CodeDropZone},
        data: () => {
            let newCode = {main: {name: "New", id: 0, texts: []}, subcodes: {}};
            return {
                codes: [],
                defaultNewCode: newCode,
            }
        },
        mounted() {
            this.axios.get("/code/list").then((res) => {
                let codes = res.data

                this.codes = [];
                for(let c of codes) {
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
                    if(typeof code === "boolean" || !code)
                        return;
                    let c = code;
                    if(parentCode.main.id === 0) {
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
                let newCatName = prompt("Name of new code?");
                if (newCatName === null) {
                    // prompt was cancelled
                    return Promise.reject(true).then(() => {}, () => {});
                }

                return this.axios.post("/code/create", {
                    code: newCatName,
                    parentCodeID: parentCodeID
                }).then(function (res) {
                    let newCode = res.data
                    if(!parentCodeID) {
                        newCode = prepareCode(newCode);
                    }
                    codes.push(newCode);
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
            }
        }
    }
</script>

<style scoped>
    .code-wrapper {
        border: 2px solid blue;
        border-radius: .25em;
        margin: 5px;
        /*padding-bottom: 5px;*/
    }

    .code-wrapper, .subcode {
        margin-bottom: 1%;
    }

    .parent-code-option {
        padding: 0;
        margin-top: 50%;
        line-height: 1em;
        width: 1em;
        text-align: center;
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
</style>