<template>
<div class="project-container">
    <h4></h4>
</div>
</template>

<script>
export default {
    computed: {
        project() {
            return this.$store.getters.currentProject;
        }
    },
name: "ProjectHub"

}
</script>

<style scoped>
.project-container {
    width: 100%;
}
</style>