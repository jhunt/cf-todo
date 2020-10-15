<template>
  <div id="app">
    <list :list="this.list"></list>
    <footer>Copyright &copy; 2020 <a href="https://jameshunt.us">James</a> <a href="https://huntprod.com">Hunt</a>.</footer>
  </div>
</template>

<script>
import List from './components/List.vue'

export default {
  name: 'App',
  components: {
    List
  },
  data() {
    return {
      loading: true,
      list: {
        name: "Things To Do",
        items: []
      }
    }
  }
}
</script>

<style lang="scss" scoped>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  position: relative;

  footer {
    font-size: 9pt;
    color: #888;
    position: fixed;
    bottom: 0px;
    left: 0px; right: 0px;
    padding: 1em;
    text-align: center;

    a {
      text-decoration: none;
      color: inherit;
    }
  }
}
</style>
