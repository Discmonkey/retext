<template>
    <div>
        <input type="file" :multiple="multiple" ref="form" v-on:change="upload()" :accept="acceptedFiles">
        <button class="btn btn-primary" @click="clickFile()"
            v-b-tooltip.bottom="(tooltip ? tooltip + ' | ' : '') +  ' Max total upload size: 2MB'">
            <slot>{{fileType}}</slot> <i class="fa fa-question-circle"></i>
        </button>
    </div>
</template>

<script>
    export default {
        name: "UploadFile",
        props: {
            fileType: {type: String},
            tooltip: {type: String, default: ""},
            acceptedFiles: {type: String},
            multiple: {
                type: Boolean,
                default: false
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
                formData.append("fileType", this.fileType);
                this.$refs.form.files.forEach((file) => {
                    formData.append("file", file);
                });

                this.axios.post("/file/upload",
                    formData,
                    {
                        headers: {
                            "Content-Type": 'multipart/form-data'
                        }
                    }
                )
                .then(
                    received => this.$emit("success", received.data)
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