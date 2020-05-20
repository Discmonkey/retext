<template>
    <div>
        <div>
            This Display a Document
        </div>

        <div>
            <div class="example-btn">

                <div class="container">
                    <div class="row">
                        <div class="col-md-2"></div>
                        <div class="col-md-10"></div>
                    </div>
                </div>

                <file-upload
                        class="btn btn-primary"
                        post-action="/file/upload"
                        extensions="txt"
                        :multiple="true"
                        :size="1024 * 1024 * 10"
                        v-model="files"
                        ref="upload">
                    <i class="fa fa-plus"></i>
                    Select files
                </file-upload>

                <button type="button" class="btn btn-success" v-if="!$refs.upload || !$refs.upload.active" @click.prevent="$refs.upload.active = true">
                    <i class="fa fa-arrow-up" aria-hidden="true"></i>
                    Start Upload
                </button>

                <button type="button" class="btn btn-danger"  v-else @click.prevent="$refs.upload.active = false">
                    <i class="fa fa-stop" aria-hidden="true"></i>
                    Stop Upload
                </button>
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
                files: []
            }
        }
    }
</script>

<style scoped>

</style>