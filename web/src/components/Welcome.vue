<template>
  <div class="centered">
    <br />
    <div>
      <b-img
        center
        :src="require('../assets/logo.png')"
        alt="Center image"
      ></b-img>
    </div>
    <br />

    <h1 class="centered">Welcome to Pull Requests status tool.</h1>

    <p>
      In order to communicate with Github API, a personal access token is
      required.<br />
      Please follow the
      <a href="https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token" target="_blank" rel="noopener"
        >documentation</a
      >
      to create one.
    </p>
  </div>
</template>

<script>
export default {
  name: "Welcome",
  props: {},
  data() {
    return {};
  },
  methods: {},
  created() {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
a {
  color: #42b983;
}

p {
  text-align: center;
}

h1 {
  text-align: center;
}

</style>