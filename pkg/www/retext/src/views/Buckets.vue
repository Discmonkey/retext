<template>
    <div class="bucket-container">

    </div>
</template>

<script>

import {actions} from "@/store";

export default {
    beforeCreate() {
        this.$store.dispatch(actions.INIT_CONTAINERS);
    },

    name: "Buckets"
}
</script>

<style scoped>

</style>