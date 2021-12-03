<script>
import Youtube from './youtube.svelte'
export let link;
export let first;
export let last;

$: isYoutube = link?.metadata?.youtube_id?.length > 0 
</script>

{#if isYoutube}
    <Youtube first={first} last={last} link={link} />
{:else}
<div class="link-item-container pa2 flex" 
    class:mt2={first}
    class:mb2={!last}>
    <div class="link-item flex flex-column">
        <div class="bold">
            <a href={link.href}>
                <span class="clmp-2">{link.metadata.title}</span>
            </a>
        </div>
        {#if link?.metadata?.description?.length > 0}
            <div class="desc clmp-2 pt1">
                {@html link.metadata.description}
            </div>
        {/if}
    </div>
</div>
{/if}


<style>
.link-item-container {
    border-left: 3px solid var(--background-1);
    border-radius: 3px;
    background-color: var(--background-2);
    max-width: 428px;
}

.desc {
    font-size: 0.9rem;
}
</style>
