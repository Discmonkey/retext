<template>
    <div class="container">
        <div class="row">
            <div class="col-md-6 code-container" v-for="container in $store.getters.containers" v-bind:key="container.id">
                <ToBucket :code-container="container"></ToBucket>
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