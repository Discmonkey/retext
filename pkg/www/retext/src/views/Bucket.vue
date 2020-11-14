<template>
    <div class="container">
        <div class="col-12">
            <div class="row">
                <div class="grid">

                </div>
            </div>
        </div>
    </div>
</template>

<script>


export default {

    computed: {

    },
    name: "Bucket",
    data() {
        return {
            container: this.$store.getters.idToContainer[this.$route.params.codeId]
        }
    },

    mounted() {

    }
}
</script>

<style scoped>
.grid {
    position: relative;
}
.item {
    display: block;
    position: absolute;
    width: 100px;
    height: 100px;
    margin: 5px;
    z-index: 1;
    background: #000;
    color: #fff;
}
.item.muuri-item-dragging {
    z-index: 3;
}
.item.muuri-item-releasing {
    z-index: 2;
}
.item.muuri-item-hidden {
    z-index: 0;
}
.item-content {
    position: relative;
    width: 100%;
    height: 100%;
}
</style>