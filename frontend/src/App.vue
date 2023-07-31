<template>
  <eq-window
      style="margin: 10px auto 0; width: 100vh; height: 90vh;"
      :title-draggable="true"
      title="ProjectEQ Expansions Client Switcher Utility"
  >
    <div class="row">
      <div
          v-for="(expansion, expansionId) in expansions"
          class="col-12"
      >
        <div
            class="mt-1"
            :style="(isExpansionSelected(expansionId) ? 'opacity: 1' : 'opacity: .5')"
            @mouseover="selectedExpansiones[expansionId] = true"
            @mouseout="selectedExpansiones[expansionId] = false"
        >
          <img
              :style="'width: 56px;' + (isExpansionSelected(expansionId) ? 'border: 2px solid #dadada; border-radius: 7px;' : 'border: 2px solid rgb(218 218 218 / 30%); border-radius: 7px;')"
              :src="getExpansionImage(expansion)" style="width: 56px;"
          >
          ({{ expansionId }})
          {{ expansion.name }}
        </div>
      </div>
    </div>

  </eq-window>
</template>

<style>

</style>

<script>
import {EXPANSIONS_FULL} from "./expansions/eq-expansions.ts";
import EqWindow          from "./components/eq-ui/EQWindow.vue";

export default {
  components: { EqWindow },
  data() {
    return {
      expansions: EXPANSIONS_FULL,
      selectedExpansiones: {},
    }
  },
  mounted() {

    console.log("hello")
  },
  methods: {
    isExpansionSelected: function (expansionId) {
      return this.selectedExpansiones[expansionId]
    },
    getExpansionImage(expansion) {
      return new URL('/src/assets/expansions/' + expansion.icon, import.meta.url).href;
    },
  }
}
</script>