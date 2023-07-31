<template>
  <div v-show="isActive">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'EqTab',
  props: {
    name: { required: true },
    selected: { default: false }
  },

  data() {
    return {
      isActive: false
    };
  },
  computed: {
    href() {
      return '#' + this.name.toLowerCase().replace(/ /g, '-');
    }
  },

  beforeDestroy() {
    const index = this.$parent.tabs.indexOf(this);
    this.$parent.tabs.splice(index, 1);
  },
  mounted() {
    this.isActive = this.selected;
  },
  created() {
    this.$parent.tabs.push(this);
  }
}
</script>

<style scoped>

</style>
