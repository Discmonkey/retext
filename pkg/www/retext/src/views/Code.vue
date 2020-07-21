<template>
  <div class="home">
    <DocumentDisplay />
    <CodeList />
  </div>
</template>

<script>
// @ is an alias to /src
import CodeList from '@/components/CodeList'
import DocumentDisplay from "@/components/DocumentDisplay";
export default {
  name: 'Code',
  components: {
    CodeList,
    DocumentDisplay
  },
}
</script>

<style>
  .home {
    display: grid;
    grid-template-rows: 100%;
    grid-template-columns: 75% 25%;
  }
</style>