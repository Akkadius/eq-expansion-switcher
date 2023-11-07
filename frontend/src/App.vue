<template>
  <eq-window
      style="margin: 10px 0 0; width: 100%; height: 99%;"
      :title-draggable="true"
      :title="'ProjectEQ Expansions Client Switcher Utility (v' + version + ')'"
      class="main-window pb-0"
      id="main-window"
  >
    <div
        style="position: absolute; right: 50px; top: -8px; z-index: 999999; cursor: pointer; "
        @click="maximizeApp()"
    >
      <i class="gg-maximize hover-highlight" style="--ggs: .7;"></i>
    </div>

    <div
        style="position: absolute; right: 20px; top: -12px; z-index: 999999; cursor: pointer; "
        @click="closeApp()"
    >
      <i class="gg-close hover-highlight" style="--ggs: .7;"></i>
    </div>

    <app-initialize
        @initialized="checkAppInitialized()"
        v-if="!isInitialized"
    />

    <patcher v-if="isInitialized"/>

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
import EqWindow                                           from "./components/eq-ui/EQWindow.vue";
import {AppInitializationCheck, CheckForUpdate, CloseApp, GetEnv} from "../wailsjs/go/main/App.js";
import EqTabs                                             from "./components/eq-ui/EQTabs.vue";
import EqTab                              from "./components/eq-ui/EQTab.vue";
import {WindowToggleMaximise}             from "../wailsjs/runtime/runtime.js";
import AppInitialize                      from "@/components/AppInitialize.vue";
import Patcher                            from "@/components/Patcher.vue";

export default {
  components: { Patcher, AppInitialize, EqTab, EqTabs, EqWindow },
  data() {
    return {
      isInitialized: false,
      version: "0.0.0",
    }
  },
  async mounted() {
    await CheckForUpdate()

    await this.checkAppInitialized()

    const r = await GetEnv()
    if (r) {
      this.version = r.version
    }


  },
  methods: {
    async checkAppInitialized() {
      const r = await AppInitializationCheck()
      if (r) {
        this.isInitialized = r.is_initialized
      }
    },
    closeApp() {
      CloseApp()
    },
    maximizeApp() {
      WindowToggleMaximise()
    },
  }
}
</script>