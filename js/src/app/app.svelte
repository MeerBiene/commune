<script>
import { onMount } from 'svelte';
import { store } from '../store/store.js'
import { fade } from 'svelte/transition'
import Switcher from './switcher/switcher.svelte'
import Sidebar from './sidebar/sidebar.svelte'
import View from './view/view.svelte'

onMount(() => {
    store.activate()
})

</script>


{#if $store.active}
<div class="foundation" transition:fade="{{duration: 100}}">
    <div class="container">

        <Switcher />

        <Sidebar />

        <View />

    </div>
</div>
{:else}
    <div class="spinner gr-center"></div>
{/if}
