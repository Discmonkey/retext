<template>
    <div class="container">
        <div class="row">
            <div class="col-12 text-center">
                <h3>Kidney Cancer Study</h3>
                <p>Ten interviews with patients and cargivers for concept testing</p>
            </div>
        </div>

        <div class="row text-center space-bottom">
            <div class="col-4">
                <button class="btn btn-primary">
                    Upload Directory
                </button>
            </div>
            <div class="col-4">
                <UploadFile/>
            </div>

            <div class="col-4">
                <button class="btn btn-primary">
                    Upload Demographics
                </button>
            </div>

        </div>

        <div class="row">
            <div class="col-12 text-center">
                <h3>Uploaded So Far</h3>
            </div>
        </div>
        <div class="row">

            <div class="col-6">

            </div>

            <div class="col-6">

            </div>
        </div>
    </div>


</template>

<script>
    import UploadFile from "../components/files/UploadFile";
    export default {
        components: {UploadFile},
        component: {
            UploadFile
        },
        name: "Upload",
    }
</script>

<style scoped>
    .row {
        margin-top: 20px;
    }

    .space-bottom {
        margin-bottom: 2em;
    }
</style>