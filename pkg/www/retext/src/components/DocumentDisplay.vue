<template>


    <div class="container-fluid">

        <div class="row">
            <div class="col-md-12 document-display">
                <TextRenderer :text="currentText" :document-i-d="selected"></TextRenderer>
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
        data: () => {
            return {
                files: [],
                uploadedFiles: [],
                currentText: "",
                selected: "",
            }
        },

        mounted() {
            this.loadDocument(this.$route.params.documentID);
        },

        methods: {
            loadDocument: function(documentName) {
                this.axios.get(`/file/load?key=${documentName}`).then(res => {
                    this.currentText = res.data;
                    this.selected = documentName;
                })
            },
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
    .document-display {
        max-height: 90%;
        overflow-y:hidden;
    }

    .row {
        padding-top: 10px;
    }
</style>