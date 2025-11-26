<template>
  <router-view />
</template>

<script setup>
</script>

<style>
:root {
  --ui-primary: #2b68ff;
  --ui-primary-2: #5f8bff;
  --ui-accent: #11cfa1;
  --ui-bg: #f7f8fb;
  --ui-card: #ffffff;
  --ui-border: #e5e8ef;
  --ui-text: #1f2d3d;
  --ui-subtext: #6b778c;
  --ui-radius: 12px;
  --ui-shadow: 0 10px 30px rgba(27, 55, 99, 0.08);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial,
    'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol',
    'Noto Color Emoji';
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background: radial-gradient(120% 120% at 20% 20%, rgba(91, 140, 255, 0.08), transparent 40%),
    radial-gradient(120% 120% at 80% 0%, rgba(17, 207, 161, 0.08), transparent 35%),
    var(--ui-bg);
}

#app {
  width: 100%;
  height: 100vh;
}

code {
  font-family: 'Courier New', Courier, monospace;
  background-color: #f0f2f7;
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.9em;
}
</style>
