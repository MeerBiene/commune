<script>
import {onMount} from 'svelte'
import { play } from '../../../../../utils/icons.js'

export let link;
export let first;
export let last;

$: id = link?.metadata?.youtube_id
$: title = link?.metadata?.title
$: description = link?.metadata?.description
$: href = link?.href

export let reply;

let active = false;

let iframe;

$: src = `https://www.youtube.com/embed/${id}?autoplay=1&enablejsapi=1`


function activate() {
    active = true
}

</script>


<div class="link-item-container pa2 flex " 
    class:mt2={first}
    class:mb2={!last}
    on:click={activate} >

    <div class="link-item flex flex-column">
        <div class="bold">
            <a href={link.href}>
                <span class="clmp-2">{link.metadata.title}</span>
            </a>
        </div>
        {#if link?.metadata?.description?.length > 0}
            <div class="desc clmp-2 pt1">
                {link.metadata.description}
            </div>
        {/if}


        {#if active}

            <div class="frame-con mt3">
                <iframe 
                    title={title}
                    bind:this={iframe}
                    src="https://www.youtube.com/embed/{id}?autoplay=1&enablejsapi=1"
                    width="100%" 
                    height="100%" 
                    frameborder="0" 
                    allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" 
                    allowfullscreen>
                </iframe>
            </div>
        {:else}
            <div class="vp-i mt3 gr-default bg-img pointer"
            style="background-image: url(https://img.youtube.com/vi/{id}/mqdefault.jpg);">
                <div class="play-icon-container gr-center gr-default pa3">
                    <div class="play-icon gr-center">
                        {@html play}
                    </div>
                </div>
            </div>
        {/if}


    </div>

</div>


<style>
.frame-con {
  overflow: hidden;
  padding-top: 56.25%;
  position: relative;
  border-radius: 2px;
  background-color: var(--primary-lightest-gray);
}
 
.frame-con iframe {
   border: 0;
   height: 100%;
   left: 0;
   position: absolute;
   top: 0;
   width: 100%;
}

.link-item-container {
    border-left: 3px solid var(--background-1);
    border-radius: 2px;
    background-color: var(--background-2);
    max-width: 428px;
}

.desc {
    font-size: 0.9rem;
}
.vp-i {
    height: 225px;
    width: 400px;
}
.play-icon-container {
    background-color: rgba(0,0,0,.6);
    border-radius: 500px;
}

.play-icon {
    width: 24px;
    height: 24px;
    fill: var(--text);
}
.play-icon:hover {
    fill: var(--white);
}
</style>

