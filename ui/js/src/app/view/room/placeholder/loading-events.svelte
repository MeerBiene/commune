<script>
import { onMount } from 'svelte'
import LoadingEventItem from './loading-event-item.svelte'

export let topics;

let count = 50;

let container;

onMount(() => {
    if(!topics) {
        container.scrollTop = container.scrollHeight;
    }
})

</script>

<div class="pe-c" class:fht={topics}>
    {#if topics}
        <div class="tp-c scrl">
            <div class="topics mt5">
                {#each Array(count) as _, i}
                    <LoadingEventItem topics={topics} />
                {/each}
            </div>
        </div>
    {:else}
        <div class="p-events scrl pa3" bind:this={container}> 
            {#each Array(count) as _, i}
                <LoadingEventItem />
            {/each}
        </div>
    {/if}
</div>

<style>
.pe-c {
    display: grid;
    height: 100%;
    overflow: hidden;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    height: calc(100% - 68px);
}

.fht {
    height: 100%;
}

.p-events {
    overflow: hidden auto;
}
.tp-c {
    display: grid;
    overflow: hidden auto;
}
.topics {
    max-width: 960px;
    width: 100%;
    justify-self: center;
}

</style>
