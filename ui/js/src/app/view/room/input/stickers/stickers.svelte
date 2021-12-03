<script>
import { onMount } from 'svelte'
import { stickers } from '../../../../utils/icons.js'


export let room;

$: roomID = room?.room_id

let active = false;

function insert() {
    active = !active;
    let opts = {
        container: container,
        dispatch: true,
    }
}


let container;


import { onDestroy, tick } from 'svelte'
import { store } from '../../../../store/store.js'
import { makeid } from '../../../../utils/utils.js'
import { search, close, em0} from '../../../../utils/icons.js'


export let isActive;


$: matrix = $store.accounts.filter(account => account.user_id == $store.active_account)[0]?.matrix

function sendSticker(e) {
    //let key = e.target.getAttribute('alt')
    //let shortcode = e.target.getAttribute('title')


    let tempid = makeid(32)

    let content = {
        body: "mr. piggy",
        info: "mr. piggy",
        url: "/static/stickers/default/pig.png",
        msgtype: "commune.sticker",
    }

    matrix.sendEvent(roomID, "m.room.message", content, tempid, (err, res) => {
        console.log(res)
    });

    active = false

}


let eventID;

let dispatch = false;

let frequent = []


let userReactions;

onMount(() => {

    if(isActive) {
        active = true
        return
    }
})

async function focusSearchInput() {
    await tick()
    searchInput.focus()
}

$: if(active) {
    setTimeout(() => {
        document.addEventListener('click', killMe)
        stickerContainer.addEventListener('scroll', highlightNav)
        document.addEventListener('keydown', escape)
    }, 1)
    focusSearchInput()
} else {
    document.removeEventListener('click', killMe)
    if(stickerContainer) {
        stickerContainer.removeEventListener('scroll', highlightNav)
    }
    document.removeEventListener('keydown', escape)
}

onDestroy(() => {
    document.removeEventListener('click', killMe)
    if(stickerContainer) {
        stickerContainer.removeEventListener('scroll', highlightNav)
    }
    document.removeEventListener('keydown', escape)
})

let escape = (e) => {
    if(e.code == 'Escape' || e.key == 'Escape') {
        kill()
    }
}

let repeat = false;

let highlighted = 'em-frequently'

let highlightNav = (e) => {
    let el0 = document.getElementById('em-frequently')
    let el1 = document.getElementById('em-people')
    let el2 = document.getElementById('em-nature')
    let el3 = document.getElementById('em-food')
    let el4 = document.getElementById('em-travel')
    let el5 = document.getElementById('em-activities')
    let el6 = document.getElementById('em-objects')
    let el7 = document.getElementById('em-symbols')
    let el8 = document.getElementById('em-flags')
    if(stickerContainer.scrollTop >= el0.offsetTop) {
        highlighted = `em-frequently`
    }
    if(stickerContainer.scrollTop >= el1.offsetTop) {
        highlighted = `em-people`
    }
    if(stickerContainer.scrollTop >= el2.offsetTop) {
        highlighted = `em-nature`
    }
    if(stickerContainer.scrollTop >= el3.offsetTop) {
        highlighted = `em-food`
    }
    if(stickerContainer.scrollTop >= el4.offsetTop) {
        highlighted = `em-travel`
    }
    if(stickerContainer.scrollTop >= el5.offsetTop) {
        highlighted = `em-activities`
    }
    if(stickerContainer.scrollTop >= el6.offsetTop) {
        highlighted = `em-objects`
    }
    if(stickerContainer.scrollTop >= el7.offsetTop) {
        highlighted = `em-symbols`
    }
    if(stickerContainer.scrollTop >= el8.offsetTop) {
        highlighted = `em-flags`
    }
}


let killMe = (e) => {
    if(e.target !== picker && 
        !picker?.contains(e.target) && 
        e.target !== searchInput && 
        !e.target.classList.contains('tone-option') && 
        !e.target.classList.contains('sticker-nav-item') && 
        !e.target.classList.contains('search-icon') && 
        e.target !== searchContainer){
        if(!repeat) {
            kill()
        } else {
            repeat = false
        }
    }
}



function kill(e) {
    active = false
    query = ''
    searchInput.value = null
    stickerContainer.scrollTop = 0
    resetHovered()
    if(dispatch) {
        window.killSticker(true)
    }
    dispatch = false
}

let picker;

$: bounding = container?.getBoundingClientRect()

$: inside = (bounding?.top + picker?.getBoundingClientRect()?.height) < document.body.clientHeight 

$: at = document.body.clientHeight - picker?.getBoundingClientRect()?.height-80

$: topNormal = inside ? bounding?.top : at

$: calcTop = bounding?.top - 450

$: top = calcTop

$: rightOffset = dispatch ? -14 : 0
$: right = document.body.clientWidth - bounding?.right  - 50

let stickerContainer;

let searchContainer;
let searchInput;
let query = '';

$: filtering = query.length > 0

$: if(filtering) {
    highlighted = ``
}


function killFilter() {
    searchInput.value = null
    query = ''
    filterSticker()
    focusSearchInput()
}

function filterSticker() {
    if(searchInput.value.length == 0) {
        highlighted = 'em-frequently'
    }
    let el = document.getElementById(`em-frequently`)
    if(el) {
        el.scrollIntoView(true)
    }
}

let moved = false;

let placeholder = `Find the perfect sticker`
$: hoveredSticker = frequent[0]?.sticker
$: hoveredShortcode = frequent[0]?.shortcode

function resetHovered() {
    placeholder = `Find the perfect sticker`
    hoveredSticker = frequent[0]?.sticker
    hoveredShortcode = frequent[0]?.shortcode
}

function changePlaceholder(e) {
    if(e.target?.className?.includes?.('sticker-key')) {
        moved = true
        placeholder = e.target.title
        hoveredSticker = e.target.getAttribute('alt')
        hoveredShortcode = e.target.title
    }
}

function toggleSection(e) {
    let el = document.querySelector(`#${e}`)
    if(el) {
        if(el?.classList.contains('no-dis')) {
            el.classList.remove('no-dis')
        } else {
            el.classList.add('no-dis')
        }
    }
}

function cat1() {
    let el = document.getElementById(`em-frequently`)
    if(el) {
        el.scrollIntoView(true)
    }
}

$: filtered = []

</script>

<div class="unset flex ml3 ">
    <div class="sticker-b-c stick">
        <div class="sticker-icon" 
            class:active={active}
            bind:this={container}
            on:click={insert}>
            {@html stickers}
        </div>
    </div>
</div>

{#if active}
<div class="layer" class:inactive={!active}>

    <div class="sticker-container" 
    class:nd={!dispatch}
    class:fly={active && !dispatch}
    style="--top:{top};--right:{right}"
    bind:this={picker}>

        <div class="sticker-search pa2 flex" 
            bind:this={searchContainer}>

            <div class="gr-center flex-one relative">
                <input 
                class="pa2"
                bind:this={searchInput}
                bind:value={query}
                on:input={filterSticker}
                placeholder={placeholder} />
                {#if filtering}
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

        <div class="inner-container" >

            <div class="sticker-nav flex flex-column">
                <div class="sticker-nav-item gr-center-start gr-default"
                    class:highlight={highlighted=='em-frequently'}
                    on:click={cat1}>
                    <img src="/static/stickers/default/pig.png"/>
                </div>
                <div class="nav-sep mv2"></div>
            </div>

            <div class="stickers-container">

                {#if filtering}
                    <div class="stickers scrl" >
                        <div class="flex flex-column">
                            <div class="con">
                                {#each filtered as item (item.order)}
                                    <li class="sticker-item gr-default" >
                                        <div class="sticker-key"
                                        alt={item.unicode}
                                        title={item.shortcode}
                                        on:click={react}>
                                            {@html item.unicode}
                                        </div>
                                    </li>
                                {/each}
                            </div>
                        </div>
                    </div>
                {/if}


                <div class="stickers scrl" 
                    class:no-dis={filtering}
                    on:mouseover={changePlaceholder}
                    bind:this={stickerContainer}>


                        <div id="em-frequently" class="flex flex-column mb3">
                            <div class="sticker-title pl1 flex">
                                <div class="gr-center">
                                    default stickers
                                </div>
                            </div>
                            <div class="con">
                                <div class="sticker-item"
                                on:click={sendSticker}>
                <img src="/static/stickers/default/pig.png"/>
                                </div>
                            </div>
                        </div>


                </div>


            </div>


        </div>

    </div>

</div>
{/if}

<style>
.sticker-b-c {
    padding-top: 13px;
    padding-bottom: 12px;
    display: grid;
}

.sticker-icon {
    fill: var(--text);
    width: 22px;
    height: 22px;
    cursor: pointer;
}
.sticker-icon:hover {
    fill: var(--white);
}
.active {
    fill: var(--white);
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

.sticker-search {
    border-bottom: 1px solid var(--background-1);
}

.sticker-container {
    position: absolute;
    top: var(--top);
    right: var(--right);
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

.nd {
    margin-right: 1rem;
    transition: margin-right 0.05s;
}

.fly {
    margin-right: 0.25rem
}

.stickers {
    padding: 0 0.5rem;
    overflow: hidden auto;
    display: grid;
    grid-template-rows: repeat();
    width: 100%;
    height: 100%;
    scroll-behavior: auto;
    background-color: var(--background-2);
    position: relative;
}

.inner-container {
    display: grid;
    grid-template-columns: auto 1fr;
    overflow: hidden;
}

.stickers-container {
    display: grid;
    grid-template-rows: 1fr auto;
    overflow: hidden;
    background-color: var(--background-5);
}

.sticker-preview {
    height: 48px;
}

.sticker-nav {
    background-color: var(--background-1);
    width: 48px;
    padding: 0.5rem;
}

.sticker-nav-item {
    width: 32px;
    height: 32px;
    cursor: pointer;
    padding: 0.25rem;
}

.sticker-nav-item:hover {
    background-color: var(--background-2);
    border-radius: 3px;
}

.nav-sep {
    height: 1px;
    background-color: var(--background-3);
    width: 100%;
}

.con {
    display: grid;
    grid-template-columns: repeat(4, 85px);
}

.sticker-title {
    text-transform: uppercase;
    font-size: 0.72rem;
    letter-spacing: 1px;
    font-weight: bold;
    color: var(--text-light);
    position: sticky;
    position: -webkit-sticky;
    top: 0;
    background-color: var(--background-2);
    padding: 0.5rem 0;
    z-index: 25;
}

.hov-e {
    font-size: 28px;
}

.hov-s {
    font-weight: bold;
    font-size: 1.1rem;
}

.sticker-item {
    padding: 0.25rem;
    cursor: pointer;
    border-radius: 3px;
}

.sticker-item img{
    transition: 0.1s;
}

.sticker-item:hover img{
    transform: rotate(15deg);
}

.sticker-key {
    font-size: 32px;
    line-height: 32px;
    font-family: Twemoji;
    justify-self: center;
    align-self: center;
}

.high {
    background-color: var(--background-4);
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

input {
    font-size: 0.9rem;
    padding-left: 0.5rem;
    width: 100%;
    height: 100%;
}

.no-dis {
    display: none;
}
.dis {
    display: block;
}
.skintone {
    font-size: 24px;
    line-height: 24px;
    cursor: pointer;
    width: 36px;
}
.tone-options {
    position: absolute;
    top: -8px;
    left: 4px;
    z-index: 35;
    padding: 0.25rem;
    border: 1px solid transparent;
}
.toggled {
    border: 1px solid var(--background-1);
    background-color: var(--background-5);
    border-radius: 3px;
}
.highlight {
    background-color: var(--background-3);
    border-radius: 3px;
}
.search-icon {
    position:absolute;
    top: 6px;
    right: 8px;
    fill: var(--text-muted);
    height: 20;
    width: 20;
}
.sia {
    fill: var(--text);
}
.tin-c {
    width: 16px;
    height: 16px;
}
</style>
