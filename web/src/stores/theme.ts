import { defineStore } from "pinia";
import { ref } from "vue";

const useThemeStore = defineStore('themeStore', () => {
  const theme = ref<'light' | 'dark'>('light');

  return {
    theme,
  }
});
