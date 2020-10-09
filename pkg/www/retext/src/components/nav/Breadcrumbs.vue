<template>
    <ul class="breadcrumb" v-if="crumbs.length" >
        <li v-for="(crumb, i) in crumbs" :key="i" class="breadcrumb-item" :class="{'active': isLast(i)}">
            <router-link v-if="!isLast(i)" :to="getTo(crumb)">{{ crumbName(crumb) }}</router-link>
            <span v-else class="last">{{ crumbName(crumb) }}</span>
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
            },
            getTo(crumb) {
                let to = {params: this.$route.params}
                if (crumb.name) {
                    to.name = crumb.name;
                } else if (crumb.path) {
                    to.path = crumb.path;
                }
                return to;
            },
            crumbName(crumb) {
                // `crumb.path` contains the full path and would have to be parsed to figure out
                //  which `params` to use as display text. If `crumb.path` contained only the child route's
                //  path, meta.name would likely be unnecessary...
                if (crumb.meta.name) {
                    return this.$route.params[crumb.meta.name];
                } else {
                    return crumb.name;
                }
            }
        }
    }
</script>

<style scoped>
li {
    font-weight: bold;
}
.breadcrumb {
    background-color: transparent;
}
</style>