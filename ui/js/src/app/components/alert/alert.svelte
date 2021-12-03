<script>
import { onMount } from 'svelte'
import { store } from '../../store/store.js'

onMount(() => {
    
  window.newAlert = (alert) => {
    store.newAlert(alert)
  }

})

$: alertsExist = $store.alerts?.length > 0

</script>

{#if alertsExist}
    <div class="alert-bar">
    </div>
{/if}

<style>
.alert-bar {
    height: 40px;
    background-color: var(--alert);
}
</style>
