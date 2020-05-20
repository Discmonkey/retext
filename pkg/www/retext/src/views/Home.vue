<template>
  <div class="home">
    <DocumentDisplay />
    <HelloWorld msg="Two component test" />
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from '@/components/HelloWorld.vue'
import DocumentDisplay from "@/components/DocumentDisplay";
export default {
  name: 'Home',
  components: {
    HelloWorld,
    DocumentDisplay
  }
}
</script>

<style>
  .home {
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: 65% 35%;
  }




</style>