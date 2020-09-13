<template>
    <div class="container">
        <div class="row">
            <div class="col-5">
                <h3 class="text-center font-weight-bold mb-4">Existing Projects</h3>

                <ToProject v-for="project " :project-date="new Date()" project-description="test description" project-name="test name"></ToProject>
            </div>

            <div class="col-7">
                <div class="form-group">
                    <h3 class="text-center font-weight-bold mb-4">Start a new Project</h3>

                    <input type="text" placeholder="Project Name" class="mb-4 w-100 d-inline-block form-control">

                    <input type="month" class="mb-4 d-inline-block form-control">

                    <textarea class="mb-4 w-100 form-control" placeholder="Project Description">

                    </textarea>

                    <textarea class="mb-4 w-100 form-control" placeholder="Initial Codes for Project (separated by ;)">

                    </textarea>

                    <div class="w-100 text-center">
                    <button class="btn btn-primary d-inline-block"> Create Project </button>
                    </div>
                </div>

            </div>
        </div>
    </div>
</template>

<script>
import ToProject from "@/components/nav/ToProject";
export default {
    name: "Project",
    data() {

    },
    components: {ToProject}
}
</script>

<style scoped>
input {
    padding: 10px;
    border-radius: 5px;
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