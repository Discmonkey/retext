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
export default {
  name: 'Home',
  components: {
    CategoryList,
    DocumentDisplay
  },
  data: () => {
    return {
      channel: {
        obj: false
      }
    }
  },
  mounted() {
    // don't know what to call this thing but channel doesn't really fit...
    let channel = this.channel;
    // eslint-disable-next-line no-unused-vars
    // todo: make the order that one() and two() are called not matter
    channel.one = (threeCb) => {
      channel.obj = {threeCb: threeCb};
    }
    // eslint-disable-next-line no-unused-vars
    channel.two = (obj, cb) => {
      // todo: require three() to return a promise and call cb on the completion of that promise?
      let x = channel.obj;
      if(!x) // drag-drop didn't land on a category
        return;

      channel.obj = false;

      x.threeCb(obj, cb);
      return x;
    }
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