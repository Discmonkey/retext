<template>


    <div class="container-fluid">

        <div class="row">
            <div class="col-md-3 border-right">
                <div class="document-select row" v-for="uploadedFile in uploadedFiles" :key="uploadedFile">
                    <button v-bind:class="{ active: uploadedFile === selected}"
                            class="btn white-background" @click="loadDocument(uploadedFile)">{{uploadedFile}} </button>
                </div>
            </div>
            <div class="col-md-9">
                <TextRenderer :text="currentText" :document-i-d="selected" :channel="channel"></TextRenderer>
            </div>
        </div>
    </div>
</template>

<script>
    import TextRenderer from "./TextRenderer";
    export default {
        name: "DocumentDisplay",

        components: {
            TextRenderer
        },

        props: ["channel"],

        data: () => {
            return {
                files: [],
                uploadedFiles: [],
                currentText: "",
                selected: "",
            }
        },

        mounted() {
            this.axios.get("/file/list").then((res) => {
                for (let f of res.data.Files) {
                    this.uploadedFiles.push(f);
                }
            })
        },

        methods: {
            loadDocument: function(documentName) {
                this.axios.get(`/file/load?key=${documentName}`).then(res => {
                    this.currentText = res.data;
                    this.selected = documentName;
                })
            },

            openUpload: function() {
                this.$modal.show("upload-modal");
            }
        },

        watch: {
            // easy way to inform the user that new stuff got uploaded
            files: function (newFiles) {
                for (let obj of newFiles) {
                    if (typeof(obj.response) === "string") {

                        let file = obj.response;

                        if (!this.uploadedFiles.includes(file)) {
                            this.uploadedFiles.push(file);
                        }
                    }
                }
            }
        },


    }
</script>

<style scoped>
    .btn {
        border: 1px solid black;
        margin-right: 8px;
        margin-bottom: 5px;
    }

    .white-background {
        background-color: white;
    }

    .active {
        -webkit-box-shadow: inset 0 0 9px 0 rgba(80,191,93,1);
        -moz-box-shadow: inset 0 0 9px 0 rgba(80,191,93,1);
        box-shadow: inset 0 0 9px 0 rgba(80,191,93,1);
    }

    .space-bottom {
        margin-bottom: 20px;
    }

    .border-right {
        border-right: 3px blue dashed;
    }

    .col-md-9 {
        padding-right: 50px;
    }

    .row {
        padding-top: 10px;
    }
</style>