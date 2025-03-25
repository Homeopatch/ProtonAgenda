<template>
  <div>
    <h1>Vue + Go App</h1>
    <p>Fffuuuu</p>
    <p>{{ message }}</p>
  </div>
</template>

<script>
import { ref, onMounted } from "vue";

export default {
  setup() {
    const message = ref("");

    onMounted(async () => {
      const res = await fetch("http://localhost:8080/api/hello");
      const data = await res.json();
      message.value = data.message;
    });

    return { message };
  },
};
</script>
