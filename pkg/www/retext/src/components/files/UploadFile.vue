<template>
    <div>
        <input type="file" ref="form" v-on:change="upload()" :accept="acceptedFiles">
        <button class="btn btn-primary" @click="clickFile()">
            Upload {{fileType}}
        </button>
    </div>
</template>

<script>
    export default {
        name: "UploadFile",
        props: ["fileType", "acceptedFiles"],
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
                let file = this.$refs.form.files[0];

                formData.append("file", file);

                this.axios.post("/file/upload",
                    formData,
                    {
                        headers: {
                            "Content-Type": 'multipart/form-data'
                        }
                    }
                )
                .then(
                    (received) => this.$emit("success", received.data)
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