<script>
import { tick, onMount, createEventDispatcher } from 'svelte'
import { fly } from 'svelte/transition'
import {gif, close, search, back} from '../../../../utils/icons.js'
import { debounce } from '../../../../utils/utils.js'

const dispatch = createEventDispatcher()

let top;
let left;

let active = false;
function activate() {

    top = icon.getBoundingClientRect().top
    left = icon.getBoundingClientRect().left - 307
    active = !active
}

let container;


onMount(() => {
    getCategories()
})

$: if(active) {
    focusSearchInput()
    setTimeout(() => {
        document.addEventListener('click', killMe)
        document.addEventListener('keydown', escape)
    }, 10)
} else {
    document.removeEventListener('click', killMe)
    document.removeEventListener('keydown', escape)
}

let repeat = false;

let killMe = (e) => {
    if(e.target !== container && 
        !container?.contains(e.target)) {
        kill()
    }
}
let escape = (e) => {
    if(e.code == 'Escape' || e.key == 'Escape') {
        kill()
    }
}

function kill(e) {
    searchResults = null
    selectedCategory = null
    searchInput.value = null
    query = ''
    active = false
}

async function focusSearchInput() {
    await tick()
    searchInput.focus()
}


let searchInput;
let query = '';

$: searching = query?.length > 0


function killFilter() {
    resetScroll()
    searchResults = null
    selectedCategory = null
    searchInput.value = null
    query = ''
    focusSearchInput()
}

async function resetScroll() {
    await tick();
    setTimeout(() => {
        catCon.scrollTop = 0
        itemsCon.scrollTop = 0
    }, 10)
}

async function fetchCategories(opts) {
    if(tenorKey == null || tenorKey == undefined) {
        return
    }
    let endpoint =
        `https://g.tenor.com/v1/categories?key=${tenorKey}&media_filter=basic`
    if(opts?.search){
        endpoint =
            `https://g.tenor.com/v1/search?q=${opts.query}&key=${tenorKey}&media_filter=minimal&limit=100`
    }
    let resp = await fetch(endpoint, {
        method: 'GET',
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}

let categories = []

let fetched = false;

function getCategories() {
    if(fetched || window.tenorCategories) {
        categories = window.tenorCategories
        categories = categories
        fetched = true
        return
    }
    fetchCategories().then(res => {
        console.log(res)
        if(res?.tags && res?.tags?.length > 0) {
            categories = res?.tags
            categories = categories
            window.tenorCategories = res?.tags
        }
    }).then(() => {
        fetched = true
    })
}

let selectedCategory;

$: categorySelected = selectedCategory != null

let fetchedCategories = {};

function selectCategory(c) {
    resetScroll()
    selectedCategory = c
    searchInput.value = c.searchterm
    query = c.searchterm
    if(!(c.searchterm in fetchedCategories)) {
        fetched = false
        let opts = {search: true, query: c.searchterm}
        fetchCategories(opts).then(res => {
            console.log(res)
            fetchedCategories[c.searchterm] = res?.results
            fetchedCategories = fetchedCategories
        }).then(() => {
            fetched = true
        })
    }
}

$: gifItems = fetchedCategories?.[selectedCategory?.searchterm] ?
    fetchedCategories?.[selectedCategory?.searchterm] : searching ? queryResults: []

$: queryResults = searchResults ? searchResults : []

let searchResults;

function filterGIFs() {
    categorySelected = null
    fetched = false
    debounce(() => {
        let opts = {search: true, query: query}
        fetchCategories(opts).then(res => {
            searchResults = res?.results
        }).then(() => {
            fetched = true
        })
    }, 500)
}

function findHeight(item) {
    let w = item.media[0].tinygif.dims[0]
    let h = item.media[0].tinygif.dims[1]
    let x = 194 / w
    return h * x
}

function selectGIF(item) {
    let i = {
        info: {
            w: item.media[0].mp4.dims[0],
            h: item.media[0].mp4.dims[1],
            size: item.media[0].mp4.size,
            preview: item.media[0].mp4.preview,
            mimetype: "image/gif",
            thumbnail_info: {
                w: item.media[0].tinygif.dims[0],
                h: item.media[0].tinygif.dims[1],
                size: item.media[0].tinygif.size,
                mimetype: "image/gif",
            },
            thumbnail_url: item.media[0].tinygif.url
        },
        url: item.media[0].mp4.url
    }
    dispatch('selected', i)
    kill()
}

let catCon;
let itemsCon;

let icon;

let animate = false;

function toggleAnimate() {
    animate = !animate
}

</script>

<div class="unset flex ml3">
    <div class="gif-c stick ">
        <div class="pointer" 
            bind:this={icon} 
            class:active={active} 
            class:animate={animate} 
            class:gif={!animate} 
            on:mouseover={toggleAnimate}
            on:mouseout={toggleAnimate}
            on:click={activate}>
            {@html gif}
        </div>
    </div>
</div>

{#if active}

    <div class="layer" class:inactive={!active}>

        <div class="gif-container" 
        style={`--top: ${top}px;--left: ${left}px;`}
        bind:this={container}>


            <div class="search-container pa2 relative flex">

                <div class="discard-search relative ph2 pointer gr-center"
                    on:click={killFilter}
                    class:no-dis={!categorySelected && !searching}>
                    {@html back}
                </div>


                <div class="flex-one">
                    <input class="pa2" 
                    bind:this={searchInput}
                    bind:value={query}
                    on:input={filterGIFs}
                    placeholder="Search Tenor" />
                    {#if searching}
                        <div class="search-icon pointer sia" 
                        on:click|stopPropagation={killFilter}>
                            {@html close}
                        </div>
                    {:else}
                        <div class="search-icon">
                            {@html search}
                        </div>
                    {/if}
                </div>
            </div>

            <div class="gifs">
                    <div class="loading" class:no-dis={fetched}>
                        <div class="spinner-s gr-center">
                        </div>
                    </div>

                    <div class="gif-categories" class:no-dis={!fetched}>

                        <div class="d-con scrl pa2"
                            bind:this={catCon}
                            class:no-dis={categorySelected || searching}>
                            {#each categories as category (category.name)}
                                <div class="gif-category bg-img" 
                                style="background-image: url({category.image});">
                                    <div class="category-mask gr-default"
                                        on:click="{e => selectCategory(category)}">
                                        <div class="category-name gr-center">
                                            {category.searchterm}
                                        </div>
                                    </div>
                                </div>
                            {/each}
                        </div>


                        <div class="d-con-a scrl pa2"
                            bind:this={itemsCon}
                            class:no-dis={!categorySelected && !searching}>

                            <div class="dc">
                                {#each gifItems as item, i (item.id)}
                                    {#if i % 2 == 0}
                                        <div class="gif-item" 
                                        in:fly="{{ y: 20, duration: 200 }}"
                                        on:click={e => selectGIF(item)}>
                                        <img width="194" height={findHeight(item)} loading="lazy" src={item.media[0].tinygif.url} />
                                    </div>
                                        {/if}
                                {/each}
                            </div>

                            <div class="dc">
                                {#each gifItems as item, i (item.id)}
                                    {#if i % 2 != 0}
                                    <div class="gif-item" 
                                        in:fly="{{ y: 20, duration: 200 }}"
                                        on:click={e => selectGIF(item)}>
                                        <img width="194" height={findHeight(item)} loading="lazy" src={item.media[0].tinygif.url} />
                                    </div>
                                        {/if}
                                {/each}
                            </div>

                        </div>


                    </div>

            </div>

        </div>
    </div>

{/if}

<style>
.gif-c {
    padding-top: 12px;
    padding-bottom: 12px;
}
.gif {
    fill:var(--text);
}
.active {
    fill:var(--white);
}
.gif:hover {
    fill:var(--white);
}


.animate {
    fill: var(--primary);
    transition: 0.1s;
    animation: dance 0.2s 2 alternate;
    animation-timing-function: ease-in;
}

@keyframes dance {
  20% {
    transform: rotate(25deg);
  }
  40% {
    transform: rotate(0deg);
  }
  60% {
    transform: rotate(-25deg);
  }
  80% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(25deg);
  }
}

.inactive {
    visibility: hidden;
    pointer-events: none!important;
}
.layer {
    position: fixed;
    pointer-events: none!important;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: none!important;
    z-index: 1002;
}

.gif-container {
    position: absolute;
    top: calc(var(--top) - 446px);
    left: var(--left);
    pointer-events: auto;
    overflow: hidden;
    height: 420px;
    width: 424px;
    border-radius: 8px;
    background-color: var(--background-2);
    color: var(--text);
    font-size: 14px;
    line-height: 1.4;
    outline: 0;
    transition-property: transform, visibility, opacity;
    box-shadow: 0 10px 20px rgba(0,0,0,.18);
    z-index: 10000;
    overflow: hidden;
    display: grid;
    grid-template-rows: auto 1fr;
}
.gifs {
    overflow: hidden;
    display: grid;
}

input {
    font-size: 0.9rem;
    padding-left: 0.5rem;
    width: 100%;
    height: 100%;
}

.search-container {
    border-bottom: 1px solid var(--background-1);
}
.search-icon {
    position:absolute;
    top: 15px;
    right: 15px;
    fill: var(--text-muted);
    height: 20;
    width: 20;
}
.sia {
    fill: var(--text);
}
.loading {
     position: absolute;
     top: 0;
     bottom: 0;
     left: 0;
     right: 0;
     width: 100%;
     height: 100%;
     display: grid;
}

.gif-categories {
    display: grid;
    overflow: hidden;
}

.d-con {
    display: grid;
    grid-template-columns: repeat(2, auto);
    grid-column-gap: 0.5rem;
    grid-row-gap: 0.5rem;
    overflow: hidden auto;
}

.d-con-a {
    display: grid;
    grid-template-columns: auto auto;
    grid-column-gap: 0.5rem;
    overflow: hidden auto;
}

.dc {
    display: grid;
    grid-template-columns: 100%;
    grid-row-gap: 0.5rem;
    grid-template-rows: repeat();
}

.gif-category {
    height: 100px;
    width: 100%;
    background-color: var(--background-3);
    display: grid;
    border-radius: 8px;
}

.category-mask {
    border-radius: 8px;
    background-color: rgba(0,0,0,.4);
    width: 100%;
    height: 100%;
    font-weight: bold;
    border: 2px solid transparent;
    cursor: pointer;
    transition: 0.1s;
    color: var(--white);
    font-size: 1.1rem;
}
.category-mask:hover {
    border: 2px solid var(--green);
    background-color: rgba(0,0,0,.6);
}

.gif-item{
    transition: 0.1s;
    border-radius: 8px;
}

.gif-item:hover img{
    border: 2px solid var(--green);
    cursor: pointer;
}

.gif-item img {
    vertical-align: middle;
    border-radius: 8px;
    border: 2px solid transparent;
    transition: 0.1s;
    background-color: var(--background-3);
}

.scrl  {
    overflow-y: auto;
    scrollbar-width: thin;
    scrollbar-color: var(--background-1) transparent;
    scroll-behavior: auto;
}

.scrl::-webkit-scrollbar {
  width: 4px;
    border-radius: 1px;
}
.scrl::-webkit-scrollbar-track {
    background: transparent;
}
.scrl::-webkit-scrollbar-thumb {
    background-color: var(--background-1);
}

.discard-search {
    fill: var(--text);
}
.discard-search:hover {
    fill: var(--white);
}
.no-dis{
    display: none;
}
.mt {
    margin-top: 13px;
}
</style>
