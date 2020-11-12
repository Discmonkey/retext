<template>
    <div class="container">
        <div class="row">
            <div class="col-md-6 code-container" v-for="code in $store.getters.containers" v-bind:key="code.id">
                <ToBucket code-container="code"></ToBucket>
            </div>
        </div>
    </div>
</template>

<script>

import {actions} from "@/store";
import ToBucket from "@/components/nav/ToBucket";

export default {
    components: {ToBucket},
    beforeCreate() {
        this.$store.dispatch(actions.INIT_CONTAINERS);
    },

    name: "Buckets"
}
</script>

<style scoped>
    .code-container {
        padding: 20px;
    }

</style>