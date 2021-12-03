<script>
import { fade, fly } from 'svelte/transition'
import { onMount, createEventDispatcher } from 'svelte'
const dispatch = createEventDispatcher()

import Menu from './menu.svelte'
import Create from './create-new-server.svelte'
import Join from './join-server.svelte'

let page = 'menu'

$: menu = page == 'menu'
$: create = page == 'create'
$: join = page == 'join'

let active = false;

onMount(() => {
    active = true
})

function kill() {
    dispatch('kill', true)
}

function showCreate() {
    page = 'create'
}

function showMenu() {
    page = 'menu'
}

function showJoin() {
    page = 'join'
}


</script>

{#if active}
<div class="mask gr-default" 
    on:click|self={kill}
    transition:fade="{{duration: 100}}">
    <div class="modal gr-center"
        in:fly="{{ y: -200, duration: 100 }}">
        {#if menu}
            <Menu 
            on:create={showCreate}
            on:join={showJoin}
            on:kill={kill}/>
        {/if}
        {#if create}
            <Create 
            on:go-back={showMenu}
            on:kill={kill}/>
        {/if}
        {#if join}
            <Join 
            on:go-back={showMenu}
            on:kill={kill}/>
        {/if}
    </div>
</div>
{/if}

<style>
.mask {
    transition: 0.3s;
    position: fixed;
    z-index: 1000;
    width: 100%;
    height: 100%;
    background-color: var(--mask);
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
}

.modal {
    background-color: var(--background-1);
    width: 440px;
    border-radius: 7px;
    transition: 0.2s;
    box-shadow: 0 30px 60px rgba(0,0,0,.1);
}
</style>
