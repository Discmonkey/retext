<template>
    <div class="container">
    <div class="row">
        <div class="col-12 text-center">
            <DemoDisplay :filename="$route.params.documentID"></DemoDisplay>
        </div>
    </div>
    </div>
</template>

<script>
import DemoDisplay from "../components/DemoDisplay";
export default {
    components: {DemoDisplay},

    name: "Demo"
}
</script>

<style scoped>

</style>