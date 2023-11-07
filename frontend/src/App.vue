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
import EqWindow                           from "./components/eq-ui/EQWindow.vue";
import {AppInitializationCheck, CloseApp} from "../wailsjs/go/main/App.js";
import EqTabs                             from "./components/eq-ui/EQTabs.vue";
import EqTab                              from "./components/eq-ui/EQTab.vue";
import {WindowToggleMaximise}             from "../wailsjs/runtime/runtime.js";
import AppInitialize                      from "@/components/AppInitialize.vue";
import Patcher                            from "@/components/Patcher.vue";

export default {
  components: { Patcher, AppInitialize, EqTab, EqTabs, EqWindow },
  data() {
    return {
      isInitialized: false,
    }
  },
  mounted() {
    this.checkAppInitialized()
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