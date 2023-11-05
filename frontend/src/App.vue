<template>
  <eq-window
      style="margin: 10px 0 0; width: 100%; height: 99%;"
      :title-draggable="true"
      title="ProjectEQ Expansions Client Switcher Utility"
      class="main-window pb-0"
      id="main-window"
  >
    <div
        class="hover-highlight"
        style="position: absolute; right: 50px; top: -14px; z-index: 999999; font-size: 18px; cursor: pointer;"
        @click="maximizeApp()"
    >
      🗖
    </div>

    <div
        class="hover-highlight"
        style="position: absolute; right: 30px; top: -18px; z-index: 999999; font-size: 20px; cursor: pointer;"
        @click="closeApp()"
    >
      x
    </div>

    <div class="row">
      <div class="col-4">
        <h6>ProjectEQ Expansions Client <br>Switcher Utility</h6>
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
                class="mt-1 text-left"
                :style="(isExpansionSelected(expansionId) ? 'opacity: 1' : 'opacity: .5')"
                @mouseover="selectedExpansions[expansionId] = true"
                @mouseout="selectedExpansions[expansionId] = false"
                @click="selectExpansion(expansionId)"
            >
              <img
                  :style="'width: 56px;' + (isExpansionSelected(expansionId) ? 'border: 2px solid #dadada; border-radius: 7px;' : 'border: 2px solid rgb(218 218 218 / 30%); border-radius: 7px;')"
                  :src="getExpansionImage(expansion.icon)" style="width: 56px;"
              >
              ({{ expansionId }})
              {{ expansion.name }}
            </div>
          </div>
        </div>
      </div>
      <div
          class="col-8"
          style="overflow-y: scroll; height: 96vh"
      >
        <div class="eq-window-simple">
          <button
              class='eq-button'
              @click="findClientWindow()"
              style="display: inline-block; margin: 0 0 10px;"
          >
            Find Client
          </button>

          <input
              class="form-control mt-3"
              type="text"
              v-model="clientLocation"
              placeholder="Select client location above"
              disabled
          >
        </div>
        <div
            class="eq-window-simple"
            v-if="selectedExpansion >= 0 && expansions[selectedExpansion]"
        >
          <div
              class="row"
              v-if="filesToCopy && filesToCopy.length > 0"
          >
            <div class="col-12">
              <button
                  class='eq-button'
                  @click="patchFiles()"
                  style="display: inline-block; margin: 0 0 10px;"
              >
                Patch - Total Files ({{ totalFilesCopyCount(filesToCopy) }})
              </button>

              <span class="text-muted"> Files get patched in order of expansion</span>

              <button
                  class='eq-button ml-3'
                  @click="dumpPatchFiles()"
                  style="display: inline-block; margin: 0 0 10px;"
              >
                Dump Patch Files
              </button>

              <eq-tabs
                  :selected="filesToCopy[0].name"
                  :key="now"
                  v-if="now"
              >
                <eq-tab
                    v-for="(f, i) in filesToCopy"
                    :key="f.expansion.id + '-' + selectedExpansion"
                    :name="f.expansion.name + ' (' + f.files.length + ')'"
                    :selected="i === 0"
                >
                  <img
                      :src="getExpansionImage(expansions[f.expansion.id].icon)"
                      style="width: 56px; "
                      alt=""
                  >
                  {{ f.expansion.name }}

                  <table
                      class="eq-table eq-highlight-rows mt-3"
                      style="display: table; font-size: 14px; overflow-x: scroll "
                      v-if="f.files && f.files.length > 0"
                  >
                    <tbody>
                    <tr v-for="file in f.files">
                      <td>{{ file.split('/').slice(2).join('/') }}</td>
                    </tr>
                    </tbody>
                  </table>

                </eq-tab>
              </eq-tabs>
            </div>

          </div>
          <div class="row" v-else>
            <div class="col-12">
              <h6>No files to copy for Expansion</h6>
              <div>
                <h6 class="eq-header">
                  <img
                      :src="getExpansionImage(expansions[selectedExpansion].icon)"
                      style="width: 56px; border-radius: 7px;"
                  >
                  {{ expansions[selectedExpansion].name }}
                </h6>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>


  </eq-window>
</template>

<style>
#main-window::before {
  background-size: cover;
  opacity: .06;
  background: url("@/assets/images/firiona.png") center;
}
</style>

<script>
import {EXPANSIONS_FULL}      from "./expansions/eq-expansions.ts";
import EqWindow               from "./components/eq-ui/EQWindow.vue";
import {
  CloseApp,
  DumpPatchFilesForExpansion,
  GetConfig,
  GetExpansionFiles,
  OpenFileDialogueEqDir,
  PatchFilesForExpansion
}                             from "../wailsjs/go/main/App.js";
import EqTabs                 from "./components/eq-ui/EQTabs.vue";
import EqTab                  from "./components/eq-ui/EQTab.vue";
import useAssets              from "./assets/assets.js";
import {WindowToggleMaximise} from "../wailsjs/runtime/runtime.js";

export default {
  components: { EqTab, EqTabs, EqWindow },
  data() {
    return {
      expansions: EXPANSIONS_FULL,
      selectedExpansions: {},
      selectedExpansion: -1,
      filesToCopy: [],
      clientLocation: "",
      now: Date.now()
    }
  },
  mounted() {
    this.getConfig()
  },
  methods: {
    closeApp() {
      CloseApp()
    },
    maximizeApp() {
      WindowToggleMaximise()
    },
    async patchFiles() {
      if (confirm('Are you sure you want to patch these files?')) {
        await PatchFilesForExpansion(parseInt(this.selectedExpansion))
      }
    },
    async dumpPatchFiles() {
      if (confirm('Are you sure you want to generate a dump of these patch files?')) {
        await DumpPatchFilesForExpansion(parseInt(this.selectedExpansion))
        // alert("Files dumped successfully")
      }
    },
    async getConfig() {
      const config = await GetConfig()
      if (!config) {
        return
      }
      await this.selectExpansion(config.current_expansion.toString())
      this.clientLocation = config.eq_dir
    },
    async findClientWindow() {
      this.clientLocation = await OpenFileDialogueEqDir()
    },
    setEqFolderLocation(e) {
      const files = e.target.files || e.dataTransfer.files;
      if (!files.length) {
        return;
      }
      const eqFolderlocation = files[0]

      console.log('location', eqFolderlocation)
      this.clientLocation = eqFolderlocation
    },
    isExpansionSelected: function (expansionId) {
      return this.selectedExpansion === expansionId
    },
    getExpansionImage(icon) {
      return useAssets(`/src/assets/expansions/${icon}`)
    },
    async selectExpansion(expansionId) {
      this.selectedExpansion = expansionId
      const files            = await GetExpansionFiles(expansionId)
      this.filesToCopy       = files ? files : []

      this.now = Date.now()
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