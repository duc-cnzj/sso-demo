<template>
  <div class="dashboard-container">
    <panel-group />
    <todo-list />
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup.vue'
import TodoList from './components/TodoList'
import { mapGetters } from 'vuex'

export default {
  name: 'Dashboard',
  components: {
    PanelGroup, TodoList
  },
  computed: {
    ...mapGetters([
      'name'
    ])
  }
}
</script>

<style lang="scss" scoped>

.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    text-transform: uppercase;
  }
}
</style>
