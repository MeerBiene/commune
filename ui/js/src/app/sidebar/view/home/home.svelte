<script>
import Items from './items/items.svelte'

export let sender;
</script>

<div class="dc-c">
    <Items  sender={sender} />
</div>

<style>
.dc-c {
    display: grid;
    overflow: hidden;
}
</style>
