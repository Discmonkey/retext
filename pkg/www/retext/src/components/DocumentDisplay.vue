<template>
    <div class="container">
        <div class="row">
            <div class="col-md-2">
                <div class="row">
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
                <div class="row">
                    <div class="document-select" v-for="uploadedFile in uploadedFiles" :key="uploadedFile">
                        <p>{{uploadedFile}} </p>
                    </div>
                </div>
            </div>
            <div class="col-md-10">

            </div>
        </div>
    </div>
</template>

<script>
    import FileUpload from "vue-upload-component"
    export default {
        name: "DocumentDisplay",

        components: {
            FileUpload
        },

        data: () => {
            return {
                files: [],
                uploadedFiles: [],
            }
        },

        mounted() {
            this.axios.get("/file/list").then((res) => {
                for (let f of res.data.Files) {
                    this.uploadedFiles.push(f);
                }
            })
        }


    }
</script>

<style scoped>
    .btn {
        border: 1px solid black;
    }

    .container {
        padding: 10px;
    }

    .document-select {

    }
</style>