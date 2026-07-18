import { writable } from 'svelte/store';

export const toast = writable({
  message: '',
  type: 'success',
  visible: false,
});

let timer;

export function notify(message, type = 'success') {
  toast.set({ message, type, visible: true });
  clearTimeout(timer);
  timer = setTimeout(() => {
    toast.update(current => ({ ...current, visible: false }));
  }, 2200);
}
