<template>
    <div>
        <input type="file" :multiple="multiple" ref="form" v-on:change="upload()" :accept="acceptedFiles">
        <button class="btn btn-primary font-weight-bold" @click="clickFile()">
            <slot>{{fileType}}</slot>
        </button>
    </div>
</template>

<script>
    import {API} from "@/core/API.ts";

    export default {
        name: "UploadFile",
        props: {
            fileType: {type: String},
            tooltip: {type: String, default: ""},
            acceptedFiles: {type: String},
            multiple: {
                type: Boolean,
                default: false
            },
            projectId: {
                type: Number,
                default: -1
            }
        },
        data: function() {
            return {
                file: null,
                uploading: false
            }
        },

        methods: {
            upload() {
                if (this.$refs.form.files.length === 0) {
                    return;
                }

                let formData = new FormData();
                this.$refs.form.files.forEach((file) => {
                    formData.append("files", file);
                });

                let items;
                if (this.fileType === "KSOURCE") {
                    items = API.source.post(this.projectId, formData);
                } else {
                    items = API.demo.post(this.projectId, formData);
                }

                items.then(
                    received => this.$emit("success", received)
                )
                .catch(
                    (error) => console.error(error)
                )
            },

            clickFile() {
                this.$refs.form.click();
            }
        },


    }
</script>

<style scoped>
    input {
        position: absolute;
        top: -1000px;
    }
</style>