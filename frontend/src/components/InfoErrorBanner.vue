<template>
  <div class="row" v-if="localNotification || localError">
    <div class="col-12">

      <!-- Notification -->
      <div
          class="mb-0 alert alert-warning"
          v-if="localNotification"
          role="alert"
      >
        <div class="row" @click="dismissNotification()">
          <div class="col-11">
            {{ localNotification }}
          </div>
          <div class="col-1 text-right">
            x
          </div>
        </div>
      </div>

      <!-- Error -->
      <div
          v-if="localError && !localNotification"
          class="eq-window-complex"
      >
        <div class="row" @click="dismissError()">
          <div class="col-11">
            Error: {{ localError }}
          </div>
          <div class="col-1 text-right">
            x
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  name: "InfoErrorBanner",
  props: {
    notification: {
      type: String,
      required: false
    },
    error: {
      type: String,
      required: false
    },
  },
  data() {
    return {
      localNotification: "",
      localError: "",

      notificationTimer: null,
      errorTimer: null,
    }
  },
  mounted() {
    if (this.notification && this.notification.length > 0) {
      this.sendNotification(this.notification, 5000)
    }
    if (this.error && this.error.length > 0) {
      this.localError = this.error
    }
  },
  watch: {
    notification: {
      handler(newVal) {
        console.log("[InfoErrorBanner] notification watcher [%s]", this.notification)
        this.sendNotification(this.notification)
        if (this.notification.length > 0) {
          this.dismissError()
        }
      },
    },
    error: {
      handler(newVal) {
        console.log("[InfoErrorBanner] error watcher [%s]", this.error)

        this.localError = this.error
        if (this.error.length > 0) {
          this.dismissNotification()
        }
      },
    },
  },
  methods: {
    dismissError() {
      this.localError = ''
      this.$emit("dismiss-error", true);
    },
    dismissNotification() {
      this.localNotification = ''
      this.$emit("dismiss-notification", true);
    },

    sendNotification(message) {
      this.localNotification = message

      if (this.notificationTimer) {
        clearTimeout(this.notificationTimer);
      }

      // dismiss in interval
      this.notificationTimer = setTimeout(() => {
        this.dismissNotification()
      }, 5000)
    },
  }
}
</script>

<style scoped>

</style>