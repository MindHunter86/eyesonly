<script>
  export let id = undefined;
  export let name = id;
  export let value;
  export let options = [];
  export let onChange = null;

  function handleChange(event) {
    value = event.currentTarget.value;
    if (onChange) onChange(value);
  }
</script>

<span class="select-control">
  <select {id} {name} {value} on:change={handleChange}>
    {#each options as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
</span>
