<template>
  <div class='eq-window-simple'>
    <div
        class='eq-window-title-bar'
        v-if="title"
        :style="titleDraggable ? '--wails-draggable:drag' : false"
    >{{ title }}
    </div>
    <div>
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EqWindow',
  props: {
    title: {
      type: String,
      required: false
    },
    titleDraggable: {
      type: Boolean,
      required: false,
      default: false
    }
  }
}
</script>

<style scoped>

</style>
