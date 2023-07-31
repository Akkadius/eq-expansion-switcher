<template>
  <eq-window
      style="margin: 0px; margin-top:10px; width: 98vw; height: 96vh;"
      :title-draggable="true"
      title="ProjectEQ Expansions Client Switcher Utility"
  >
    <div class="row">
      <div class="col-4">
        <h3>ProjectEQ Expansions Client <br>Switcher Utility</h3>
        <p>
          This utility will allow you to switch between ProjectEQ expansions
          without having to manually copy files around.
        </p>

        <div class="row">
          <div
              v-for="(expansion, expansionId) in expansions"
              class="col-12"
          >
            <div
                class="mt-1"
                :style="(isExpansionSelected(expansionId) ? 'opacity: 1' : 'opacity: .5')"
                @mouseover="selectedExpansions[expansionId] = true"
                @mouseout="selectedExpansions[expansionId] = false"
                @click="selectExpansion(expansionId)"
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
      </div>
      <div
          class="col-8"
          style="overflow-y: scroll; height: 92vh"
          v-if="selectedExpansion >= 0 && expansions[selectedExpansion]"
      >
        <div
            class="row"
            v-if="filesToCopy && filesToCopy.length > 0"
        >
          <div class="col-12">
            <button
                class='eq-button'
                onclick="alert('click')"
                style="display: inline-block; margin: 0 0 10px;"
            >
              Patch
            </button>

            <h3>Total Files ({{ totalFilesCopyCount(filesToCopy)}})</h3>

            <div v-for="f in filesToCopy">
              <h4 class="eq-header mt-3 mb-1">
                <img
                    :src="getExpansionImage(expansions[f.expansion.id])"
                    style="width: 56px; "
                >
                {{ f.expansion.name }}

              </h4>

              <h3>Files ({{f.files.length}})</h3>

              <table
                  class="eq-table eq-highlight-rows mt-3"
                  style="display: table; font-size: 14px; overflow-x: scroll "
                  v-if="f.files && f.files.length > 0"
              >
                <thead>
                <tr>
                  <th>File</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="file in f.files">
                  <td>{{ file.split('/').slice(2).join('/') }}</td>
                </tr>
                </tbody>
              </table>
            </div>
          </div>

        </div>
        <div class="row" v-else>
          <div class="col-12">
            <h3>No files to copy for Expansion</h3>
            <div>
              <h2>
                <img
                    :src="getExpansionImage(expansions[selectedExpansion])"
                    style="width: 56px; border-radius: 7px;"
                >
                {{ expansions[selectedExpansion].name }}
              </h2>
            </div>
          </div>
        </div>
      </div>
    </div>


  </eq-window>
</template>

<style>

</style>

<script>
import {EXPANSIONS_FULL}   from "./expansions/eq-expansions.ts";
import EqWindow            from "./components/eq-ui/EQWindow.vue";
import {GetExpansionFiles} from "../wailsjs/go/main/App.js";

export default {
  components: { EqWindow },
  data() {
    return {
      expansions: EXPANSIONS_FULL,
      selectedExpansions: {},
      selectedExpansion: -1,
      filesToCopy: []
    }
  },
  mounted() {
  },
  methods: {
    isExpansionSelected: function (expansionId) {
      return this.selectedExpansion === expansionId
    },
    getExpansionImage(expansion) {
      return new URL('/src/assets/expansions/' + expansion.icon, import.meta.url).href;
    },
    async selectExpansion(expansionId) {
      this.selectedExpansion = expansionId
      const files            = await GetExpansionFiles(expansionId)
      console.log('files', files)
      this.filesToCopy = files ? files : []

      // alert('Selected expansion: ' + expansion.name)
    },
    totalFilesCopyCount(filesToCopy) {
      let total = 0
      for (const f of filesToCopy) {
        total += f.files.length
      }
      return total
    }
  }
}
</script>