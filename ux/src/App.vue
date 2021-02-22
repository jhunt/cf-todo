<template>
  <div id="app">
    <div v-if="loading && failures == 0" class="dialog info">
      <h1>Loading...</h1>
      <p>Please wait while we get our bearings...</p>
    </div>
    <div v-else-if="!database" class="dialog error">
      <h1>{{ failureMessage }}</h1>
      <p>The database that stores todo list items has not been configured yet.</p>
      <p>Perhaps you could try the following:</p>
      <pre><code>$ cf create-service mariadb simple db1
$ watch cf service db1  # ... wait for "create succeeded"
$ cf bind-service todo db1
$ watch cf service db1  # ... wait for binding to finish
$ cf restart todo</code></pre>
      <button v-if="loading" style="background-color: #777; cursor: default;">Re-checking...</button>
      <button v-else         @click.prevent="sync()">All Fixed?  Try Again!</button>
    </div>
    <list v-else :list="this.list"></list>

    <footer>Copyright &copy; 2020 <a href="https://jameshunt.us">James</a> <a href="https://huntprod.com">Hunt</a>.</footer>
  </div>
</template>

<script>
import List from './components/List.vue'
import http from './http.js'

export default {
  name: 'App',
  components: {
    List
  },
  data() {
    return {
      loading: true,
      failures: 0,
      database: false,
      list: {
        name: "Things To Do",
        items: []
      },
      oops: [
        'Missing Database (still)',
        'Database Not Configured',
        'Oops, Not Fixed Yet!',
        'Database Still Missing',
      ]
    }
  },
  computed: {
    failureMessage() {
      return this.oops[this.failures % this.oops.length]
    }
  },
  methods: {
    sync() {
      this.loading = true
      window.setTimeout(() =>
        http.GET('/v1/ping')
          .then(ping => {
            this.database = ping.db
            this.loading = false
            if (!this.database) {
              this.failures += 1
            }
          })
          .catch(() => {
            this.loading = false
            this.failures += 1
          }), 800)
    }
  },
  mounted() {
    this.sync()
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
.dialog {
  max-width: 700px;
  margin: 4em auto;
  box-shadow: 0 0 24px #ccc;
  border: 1px solid #e0e0e0;
  padding: 6em;
  box-sizing: border-box;

  pre {
    text-align: left;
    background-color: #224;
    color: yellow;
    padding: 1em;
    font-size: 110%;
    line-height: 1.4em;
  }
  button {
    display: block;
    margin: 2em auto;
    font-size: 14pt;
    padding: 1em;
    background-color: #1bca78;
    color: #fff;
    font-weight: bold;
    border-radius: 6px;
    border: none;
    box-shadow: 0 0 8px #ccc;
    cursor: pointer;
  }
}
</style>
