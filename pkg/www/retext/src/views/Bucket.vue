<template>
    <div class="container">
        <div class="col-12">
            <div class="row">
                <div class="grid">
                    <div class="item" v-for="text in this.container.">
                        <div class="item-content">
                            <!-- Safe zone, enter your custom markup -->
                            This can be anything.
                            <!-- Safe zone ends -->
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>


export default {

    computed: {
        texts() {

        }
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