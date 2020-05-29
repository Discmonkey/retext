<template>
  <div class="home">
    <DocumentDisplay :channel="channel"  />
    <CategoryList :channel="channel" />
  </div>
</template>

<script>
// @ is an alias to /src
import CategoryList from '@/components/CategoryList'
import DocumentDisplay from "@/components/DocumentDisplay";
import {Channel} from "@/components/Channel"
export default {
  name: 'Home',
  components: {
    CategoryList,
    DocumentDisplay
  },
  data: () => {
    return {
      channel: new Channel()
    }
  },
}
</script>

<style>
  .home {
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: 60% 40%;
  }
</style>