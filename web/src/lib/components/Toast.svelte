<script>
  import { toast } from '../stores/toast.js';
</script>

<div
  class="toast"
  class:show={$toast.visible}
  class:failure={$toast.type === 'failure'}
  class:success={$toast.type !== 'failure'}
  role="status"
  aria-live="polite"
>
  {$toast.message}
</div>
