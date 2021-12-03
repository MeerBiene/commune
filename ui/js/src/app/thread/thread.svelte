<script>
import { store } from '../store/store.js'
import { onMount }  from 'svelte'
import NewThread from './new-thread.svelte'
import ViewThread from './view-thread.svelte'

export let room;

$: creating = room?.thread?.creating
$: viewing = room?.thread?.viewing
$: active = room?.thread?.active

</script>

<div class="thread-container">

    {#if creating && active}

        <NewThread room={room}/>

    {/if}

    {#if viewing && active}

        <ViewThread room={room}/>

    {/if}

</div>

<style>
.thread-container {
    width: 400px;
    background-color: var(--background-3);
    border-radius: 8px 0 0 8px;
    display: grid;
    grid-template-rows: 48px 1fr auto;
}
@media screen and (max-width: 1150px) {
    .thread-container {
        position: fixed;
        top:0;
        right: 0;
        bottom: 0;
        -webkit-filter: drop-shadow(0 8px 40px rgba(0,0,0,.32));
        filter: drop-shadow(0 8px 40px hsla(0,0%,0%,.32));
    }
}

</style>
