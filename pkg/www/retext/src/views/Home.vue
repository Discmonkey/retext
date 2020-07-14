<template>
  <div class="home">
    <DocumentDisplay />
    <CategoryList />
  </div>
</template>

<script>
// @ is an alias to /src
import CategoryList from '@/components/CategoryList'
import DocumentDisplay from "@/components/DocumentDisplay";
export default {
  name: 'Home',
  components: {
    CategoryList,
    DocumentDisplay
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