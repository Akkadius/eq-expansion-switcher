<template>
  <div
      class="row justify-content-center"
  >
    <div style="top: 30%; position: absolute">
      <eq-window>
        <div>
          By using this program you agree that you own a legal copy of EverQuest and will not use this program to
          distribute files to others.
        </div>

        <div class="mt-3" style="left: 45%; position: relative" v-if="!installing">
          <button
              class='eq-button'
              @click="initialize()"
              style="display: inline-block; margin: 0 0 10px;"
          >
            Install
          </button>
        </div>

        <loader-fake-progress
            class="mt-3"
            v-if="installing"
            :progress="0"
            :progress-text="'Installing...'"
        />

        <info-error-banner
            :notification="notification"
            :error="error"
            @dismiss-error="error = ''"
            @dismiss-notification="notification = ''"
        />

      </eq-window>
    </div>
  </div>
</template>

<script>
import EqWindow            from "@/components/eq-ui/EQWindow.vue";
import {AppInitialization} from "../../wailsjs/go/main/App.js";
import InfoErrorBanner     from "@/components/InfoErrorBanner.vue";
import LoaderFakeProgress  from "@/components/FakeProgressLoader.vue";

export default {
  components: { LoaderFakeProgress, InfoErrorBanner, EqWindow },
  data() {
    return {
      installing: false,

      // notification / errors
      notification: "",
      error: "",
    }
  },
  methods: {
    async initialize() {
      this.installing = true
      try {
        await AppInitialization();
        this.installing = false
        this.$emit("initialized")
      } catch (e) {
        this.error      = e
        this.installing = false
      }
    }
  }
}
</script>

<style scoped>

</style>