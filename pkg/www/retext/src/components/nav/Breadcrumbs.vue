<template>
    <ul class="breadcrumb" v-if="crumbs.length" >
        <li v-if="nohome" class="breadcrumb-item">
            <router-link to="/"> Home </router-link>
        </li>
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
    import {HubName, ProjectName} from "@/router";

    export default {
        name: "Breadcrumbs",
        computed: {
            crumbs() {
                const crumbs = [];

                this.$route.matched.forEach(c => {
                    if (c.name !== HubName) {
                        crumbs.push(c);
                    }
                })

                return crumbs;
            },

            nohome() {
                return this.$route.matched[0].name !== "Projects"
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
                }

                if (crumb.name === ProjectName) {
                    to.name = HubName;
                    to.path = `/project/${this.$route.params.projectId}/`
                } else if (crumb.path) {
                    to.path = crumb.path;
                }


                return to;
            },

            crumbName(crumb) {
                // rewrite the name to the actual project name,
                // TODO refactor to method  if more rewrite need to take place
                if (crumb.name === ProjectName && this.$store.getters.currentProject !== null) {
                    return this.$store.getters.currentProject.name;
                }

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