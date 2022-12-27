<script setup>
import logo from "@/assets/images/logo.png";
</script>
<template>
  <div class="flex flex-col justify-center w-fit items-center">
    <div class="flip-horizontal-bottom">
      <img class="w-12" :src="logo" />
    </div>
    <span class="font-inter">Doing science...</span>
  </div>
</template>

<script>
export default {
  name: "LoadingLogo",
};
</script>

<style>
.flip-horizontal-bottom {
  -webkit-animation: flip-horizontal-bottom 0.4s cubic-bezier(0.455, 0.030, 0.515, 0.955) alternate both infinite;
  animation: flip-horizontal-bottom 0.4s cubic-bezier(0.455, 0.030, 0.515, 0.955) alternate both infinite;
}
@-webkit-keyframes flip-horizontal-bottom {
  0% {
    -webkit-transform: rotateX(0);
    transform: rotateX(0);
  }
  100% {
    -webkit-transform: rotateX(-180deg);
    transform: rotateX(-180deg);
  }
}
@keyframes flip-horizontal-bottom {
  0% {
    -webkit-transform: rotateX(0);
    transform: rotateX(0);
  }
  100% {
    -webkit-transform: rotateX(-180deg);
    transform: rotateX(-180deg);
  }
}
</style>