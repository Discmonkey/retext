<template>
    <div class="container">
        <div class="row">
            <div class="col-md-3 border-right">
                <div class="row space-bottom">
                    <file-upload
                            class="btn"
                            post-action="/file/upload"
                            extensions="txt"
                            :multiple="true"
                            :size="1024 * 1024 * 10"
                            v-model="files"
                            ref="upload">
                        +
                    </file-upload>

                    <button type="button" class="btn"
                            v-if="!$refs.upload || !$refs.upload.active" @click.prevent="$refs.upload.active = true">
                        <i class="fa fa-arrow-up" aria-hidden="true"></i>
                        ->
                    </button>

                    <button type="button" class="btn btn-danger"  v-else @click.prevent="$refs.upload.active = false">
                        <i class="fa fa-stop" aria-hidden="true"></i>
                        x
                    </button>
                </div>
                <div class="document-select row" v-for="uploadedFile in uploadedFiles" :key="uploadedFile">
                    <button v-bind:class="{ active: uploadedFile === selected}"
                            class="btn" @click="loadDocument(uploadedFile)">{{uploadedFile}} </button>
                </div>
            </div>
            <div class="col-md-9">
                <TextRenderer :text="currentText"></TextRenderer>
            </div>
        </div>
    </div>
</template>

<script>
    import FileUpload from "vue-upload-component"
    import TextRenderer from "./TextRenderer";
    export default {
        name: "DocumentDisplay",

        components: {
            FileUpload,
            TextRenderer
        },

        data: () => {
            return {
                files: [],
                uploadedFiles: [],
                currentText: "",
                selected: "",
            }
        },

        mounted() {
            this.axios.get("http://localhost:3000/file/list").then((res) => {
                for (let f of res.data.Files) {
                    this.uploadedFiles.push(f);
                }
            })
        },

        methods: {
            loadDocument: function(documentName) {
                this.axios.get(`http://localhost:3000/file/load?key=${documentName}`).then(res => {
                    this.currentText = res.data;
                })
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
        background-color: white;
        margin-right: 8px;
        margin-bottom: 5px;
    }

    .container {
        padding: 10px;
    }

    .active {
        -webkit-box-shadow: inset 0px 0px 9px 0px rgba(80,191,93,1);
        -moz-box-shadow: inset 0px 0px 9px 0px rgba(80,191,93,1);
        box-shadow: inset 0px 0px 9px 0px rgba(80,191,93,1);
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
</style>