<template><router-view/></template>

<script>
import {ProjectActions, ProjectMutations} from "@/store/modules/project";

export default {
    name: "LoadProject",

    mounted() {
        this.$store.dispatch(ProjectActions.LOAD_PROJECTS).then(
            () => {
                this.$store.commit(ProjectMutations.SELECT_PROJECT, this.$route.params.id)
            }
        );
    },
}
</script>
