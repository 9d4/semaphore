<script setup>
import { RouterView } from "vue-router";
document.getElementsByTagName("html")[0].setAttribute("data-theme", "dark");
</script>

<template>
  <RouterView />
</template>

<style>
@import url("./assets/app.css");
</style>
