<template>
    <ul class="breadcrumb" v-if="crumbs.length" >
        <li v-for="(crumb, i) in crumbs" :key="i" class="breadcrumb-item" :class="{'active': isLast(i)}">
            <router-link v-if="!isLast(i)" :to="crumb">{{ crumb.name }}</router-link>
            <span v-else class="last">{{ crumb.name }}</span>
        </li>
    </ul>
</template>

<script>
    /*=============================================================================
      started with this link... doesn't much resemble it anymore. consider it a url-parsing-alternative...
      https://github.com/NxtChg/pieces/blob/master/js/vue/vs-crumbs/vs-crumbs.js
    =============================================================================*/
    export default {
        name: "Breadcrumbs",
        computed: {
            crumbs() {
                return this.$route.matched;
            },
        },
        methods: {
            isLast(i) {
                return i === this.crumbs.length - 1;
            }
        }
    }
</script>

<style scoped>
.breadcrumb {
    background-color: transparent;
}
</style>