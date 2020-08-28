<template>
  <div class="images-header">
    <div>
      <h1 class="title is-1">Orchestrus</h1>
    </div>
    <div>
      <button>
        <span class="icon">
          <i class="fas fa-redo-alt"></i>
        </span>
      </button>
    </div>
  </div>
  <Workers/>
</template>

<script>
import Workers from './Workers.vue'

export default {
  name: 'Menu',
  components: {
    Workers
  }
}
</script>
