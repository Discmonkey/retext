<template>
  <div class="home">
    <DocumentDisplay />
    <div>
      Area for categories
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
// import HelloWorld from '@/components/HelloWorld.vue'
import DocumentDisplay from "@/components/DocumentDisplay";
export default {
  name: 'Home',
  components: {
    DocumentDisplay
  }
}
</script>

<style>
  .home {
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: 60% 40%;
  }




</style>