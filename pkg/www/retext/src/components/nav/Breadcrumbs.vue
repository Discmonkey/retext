<template>
    <ul class="breadcrumb" v-if="$route.matched.length" >
        <li v-for="(crumb, i) in $route.matched" :key="i" class="breadcrumb-item">
            <router-link v-if="i < $route.matched.length-1" :to="crumb">{{ crumb.name }}</router-link>
            <span v-else class="last">{{ crumb.name }}</span>
        </li>
    </ul>
</template>

<script>
/*=============================================================================
  started with this... doesn't much resemble it anymore. consider it a url-parsing-alternative...
  https://github.com/NxtChg/pieces/blob/master/js/vue/vs-crumbs/vs-crumbs.js
=============================================================================*/

export default {
    name: "Breadcrumbs",
}
</script>
