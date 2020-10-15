<template>
  <div class="list">
    <h1>{{ thelist.name }}</h1>
    <ul>
    <li v-for="item in thelist.items" :key="item.id" :class="classesFor(item)">
      <input type="checkbox" v-model="item.done"
             @change="persist(item)">
      <input v-if="existingItem && existingItem.id == item.id"
             type="text" v-model="existingItem.text"
             placeholder="What needs done?"
             ref="existing"
             @blur="persist(existingItem)"
             @change="persist(existingItem)">
      <span v-else @click="edit(item)">{{ item.text }}</span>
    </li>
    <li class="add" v-if="newItem">
      <input type="checkbox" :checked="newItem.done">
      <input type="text" v-model="newItem.text"
             placeholder="What needs done?"
             ref="new"
             @blur="persist(newItem)">
    </li>
    <li class="new">
      <button @click.prevent='create({text: "", done: false })'>+</button>
    </li>
    </ul>
  </div>
</template>

<script>
import http from '../http.js'
import Vue from 'vue'

export default {
  name: 'list',
  props: ['list'],
  data() {
    return {
      thelist: this.list,
      newItem: null,
      existingItem: null
    }
  },
  mounted() {
    http.GET(`/v1/todos`)
      .then(items => this.thelist.items = items)
  },
  methods: {
    classesFor(item) {
      return item.done ? 'done' : 'active'
    },
    edit(item) {
      this.existingItem = Object.assign({}, item)
      Vue.nextTick(() => this.$refs.existing[0].focus())
    },
    create(item) {
      this.newItem = item
      Vue.nextTick(() => this.$refs.new.focus())
    },
    persist(item) {
      if (item.id) {
        // existing
        if (item.text != "") {
          http.PUT(`/v1/todos/${item.id}`, item)
            .then(that => {
              this.thelist.items.forEach(item => {
                if (item.id == that.id) {
                  Object.assign(item, that)
                }
              })
            })
        } else {
          http.DELETE(`/v1/todos/${item.id}`)
            .then(() => this.thelist.items = this.thelist.items.filter(x => x.id != item.id))
        }
        this.existingItem = null

      } else {
        // new
        if (item.text != "") {
          http.POST(`/v1/todos`, item)
            .then(that => this.thelist.items.push(that))
        }
        this.newItem = null;
      }
    }
  }
}
</script>

<style scoped lang="scss">
.list {
  max-width: 480px;
  margin: 0 auto;

  h1 {
    border-bottom: 1px solid #ccc;
  }
  ul {
    list-style: none;
    text-align: left;

    li {
      span {
        padding-right: 1em;
        cursor: pointer;
      }

      &.done {
        opacity: 0.2;
        span, input[type=text] {
          text-decoration: line-through;
        }
      }

      span, input[type=text] {
        font-family: sans-serif;
        font-size: 12pt;
        padding: 2px;
      }
      span {
        display: inline-block;
        padding-bottom: 4px;
      }
      input[type=text] {
        border: 2px solid #ccc;
        border-width: 0 0 2px 0;
        outline: 0;
      }
    }

    li.new button {
      background-color: #5caddd;
      color: #fff;
      border: 0;
      border-radius: 4pt;
      font-size: 15pt;
      cursor: pointer;
      line-height: 1.1em;
      padding: 0.2em 0.6em;
      margin: 0.2em 0 0.2em 1em;
    }
  }
}
</style>
