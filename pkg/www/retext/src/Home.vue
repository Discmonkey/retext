<template><router-view v-if="projectLoaded"/></template>

<script>
import {actions} from "@/store";

export default {
    name: "Home",
    mounted() {
        this.$store.dispatch(actions.project.SELECT_PROJECT, parseInt(this.$route.params.projectId))
    },

    computed: {
        projectLoaded() {
            return this.$store.getters.currentProject !== null &&
                this.$store.getters.currentProject.id === parseInt(this.$route.params.projectId);
        }
    }
}

</script>
