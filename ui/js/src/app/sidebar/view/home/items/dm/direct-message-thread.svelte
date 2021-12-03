<script>
import {onMount} from 'svelte'
import { store } from '../../../../../store/store.js'
export let thread;
export let single;
export let last;

import { navigate, useLocation } from 'svelte-navigator'

const location = useLocation()

$: pathname = $location.pathname

$: indicate = pathname == `/messages/${thread.alias}`

$: if(indicate) {
    document.title = thread.name
}

let goToThread =() => {
    navigate(`/messages/${thread.alias}`)
    store.updateActiveHomePage(`/messages/${thread.alias}`)
    store.updateActiveDirectMessages(`/messages/${thread.alias}`)
}

onMount(() => {
    if(indicate) {
        store.updateActiveHomePage(`/messages/${thread.alias}`)
        store.updateActiveDirectMessages(`/messages/${thread.alias}`)
    }
})
</script>

<div class="room-item-thread flex mr1 pl2">
    <div class="relative">
        {#if last}
            {#if single}
                <div class="spine-single"></div>
            {:else}
                <div class="spine-last"></div>
            {/if}
        {:else}
            <div class="spine"></div>
        {/if}
    </div>
    <div class="flex-one flex mb1"
        class:thread-name-ac={indicate}
        class:thread-name={!indicate}
        on:click={goToThread}>

        <div class="flex-one">
            {thread.name}
        </div>
        {#if thread.owner}
            <div class="icon">
        </div>
        {/if}
    </div>
</div>

<style>
.spine {
    margin-left: 15px;
    width: 0.5rem;
    height: 100%;
    border-left: 2px solid var(--text-muted);
    position: relative;
}

.spine:after {
    background: var(--text-muted);
    content: "";
    height: 2px;
    position: absolute;
    top: 16px;
    width: 10px;
}

.spine-last {
    margin-top: -15px;
    margin-left: 15px;
    width: 0.5rem;
    height: 100%;
    border-radius: 0 0 0 4px;
}

.spine-last:after {
    top: 0;
    margin-top: 0;
    margin-left: 0;
    width: 10px;
    height: 18px;
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 4px;
    content: "";
    position: absolute;
}

.spine-single {
    margin-left: 15px;
    width: 0.5rem;
    height: 20px;
    border-left: 2px solid var(--text-muted);
    border-bottom: 2px solid var(--text-muted);
    border-radius: 0 0 0 4px;
}
.thread-name {
    border-radius: 5px;
    cursor: pointer;
    color: var(--text-light);
    padding: 0.5rem 0.25rem;
    margin-left: 0.5rem;
}

.thread-name:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.thread-name-ac {
    border-radius: 5px;
    cursor: pointer;
    color: var(--white);
    padding: 0.5rem 0.25rem;
    margin-left: 0.5rem;
    background-color: var(--background-6);
}

.thread-name-ac:hover {
    background-color: var(--background-6);
    color: var(--white);
}

.room-item-threads {
    padding: 0 0 0 0.5rem;
}

.mem {
    font-size: 0.9rem;
}

.log {
    fill: var(--text);
    width: 16px;
    height: 16px;
}
</style>
