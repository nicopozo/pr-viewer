<template>
  <div>
    <b-navbar toggleable="lg" type="light" variant="light">
      <b-navbar-brand href="#">Rull Requests</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <b-navbar-nav>
          <b-nav-item
            active-class="active"
            class="nav-link"
            key="welcome"
            :to="{ name: 'Welcome' }"
            exact
            >Home</b-nav-item
          >
        </b-navbar-nav>

        <b-navbar-nav>
          <b-nav-item
            active-class="active"
            class="nav-link"
            key="listPullRequests"
            :to="{ name: 'ListPullRequests' }"
            exact
            >List</b-nav-item
          >
        </b-navbar-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown right>
            <!-- Using 'button-content' slot -->
            <template #button-content>
              <em>{{ username }}</em>
            </template>
            <b-dropdown-item
              v-if="token == null || token == ''"
              v-b-modal.modal-prevent-closing
              href="#"
              >Sign In</b-dropdown-item
            >
            <b-dropdown-item
              v-if="token != null && token != ''"
              @click="logout()"
              href="#"
              >Sign Out</b-dropdown-item
            >
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>

    <b-modal
      id="modal-prevent-closing"
      ref="modal"
      title="Login"
      @show="resetModal"
      @ok="handleOk"
    >
      <form ref="form" @submit.stop.prevent="handleSubmit">
        <b-form-group
          label="Token"
          label-for="token-input"
          invalid-feedback="Token is required"
          :state="tokenState"
        >
          <b-form-input
            id="token-input"
            v-model="token"
            :state="tokenState"
            required
          ></b-form-input>
        </b-form-group>
      </form>
    </b-modal>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "NavBar",
  data() {
    return {
      tokenState: null,
      token: "",
      username: "Login",
    };
  },
  methods: {
    logout() {
      this.token = "";
      this.username = "Login";
      this.$cookies.remove("pr-token");
    },
    login() {
      var params = {
        token: this.token,
      };

      axios
        .get("http://localhost:8081/pr-viewer/username", {
          params: params,
        })
        .then((res) => {
          if (res.data.username != "") {
            this.username = res.data.username;
            this.$store.commit("SET_TOKEN", this.token);
            this.$store.commit("SET_USERNAME", this.username);
            this.$cookies.set("pr-token", this.token, 60 * 60 * 24 * 365);
          } else {
            this.showErrorModal("Login Error!", "Verify if token is valid.");
          }
        })
        .catch((err) => {
          this.showErrorModal("Login Error!", err.Message);
        });
    },
    checkFormValidity() {
      if (this.token == "") {
        this.tokenState = false;

        return false;
      }

      return true;
    },
    resetModal() {
      this.token = "";
      this.tokenState = null;
    },
    handleOk(bvModalEvt) {
      // Prevent modal from closing
      bvModalEvt.preventDefault();
      // Trigger submit handler
      this.handleSubmit();
    },
    handleSubmit() {
      // Exit when the form isn't valid
      if (!this.checkFormValidity()) {
        return;
      }
      this.login();
      // Hide the modal manually
      this.$nextTick(() => {
        this.$bvModal.hide("modal-prevent-closing");
      });
    },
    showErrorModal(title, msg) {
      this.$bvModal.msgBoxOk(msg, {
        title: title,
        okVariant: "danger",
        centered: true,
      });
    },
  },
  created() {
    this.token = this.$cookies.get("pr-token");

    if (this.token != null && this.token != "") {
      this.login();
    }
  },
};
</script>