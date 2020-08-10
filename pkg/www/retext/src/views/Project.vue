<template>
    <div class="container">
        <div class="row">
            <div class="col-5">
                <h3 class="text-center font-weight-bold mb-4">Existing Projects</h3>
                <p>Current projects</p>
            </div>

            <div class="col-7">
                <h3 class="text-center font-weight-bold mb-4">Start a new Project</h3>

                <input type="text" placeholder="Project Name" class="mb-4 w-100 d-inline-block">

                <input class="mb-4 w-48 d-inline-block mr-4" placeholder="Month">

                <input class="mb-4 w-48 d-inline-block" placeholder="Year">

                <textarea class="mb-4 w-100" placeholder="Project Description">

                </textarea>

                <textarea class="mb-4 w-100" placeholder="Initial Codes for Project (separated by ';')">

                </textarea>

                <div class="w-100 text-center">
                <button class="btn btn-primary d-inline-block"> Create Project </button>
                </div>




            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "Project"
}
</script>

<style scoped>
input {
    padding: 10px;
    border-radius: 5px;
}

textarea {
    padding: 10px;
    border-radius: 5px;
    border-width: 2px;
    border-style: inset;
}

.mr-4 {
    margin-right: 4%;
}
.w-48 {
    width: 48%;
}

button {
    font-weight: bolder;
    padding: 10px 20px 10px 20px;
}
</style>