<script>
import { onMount } from 'svelte'
import { play, pause, refresh, volume, muted, fullscreen } from '../../utils/icons.js'
import { formatBytes } from '../../utils/utils'

export let title;
export let size;
export let url;


export let height;
export let width;

let video;


let active = false;

let played = false;
let playing = false;

let paused;

function togglePlay() {
    if(video.paused) {
        video.addEventListener('timeupdate', setTime);
        video.addEventListener('play', () => {
            realDuration = video.duration
            duration = calculateTime(video.duration)
            playing = true
            active = true
        });
        paused= false
        video.play()
    } else {
        video.pause()
        paused = true
        playing = false
    }
}


let elapsed = `0:00`;
let percent = 0;

function setTime() {
    elapsed = calculateTime(video?.currentTime)
    percent = (video.currentTime / video.duration) * 100
}

$: if(percent == 100) {
    killPlay()
    played = true
}

function killPlay() {
    elapsed = `0:00`;
    percent = 0
    playing = false
}


const calculateTime = (secs) => {
  const minutes = Math.floor(secs / 60);
  const seconds = Math.floor(secs % 60);
  const returnedSeconds = seconds < 10 ? `0${seconds}` : `${seconds}`;
  return `${minutes}:${returnedSeconds}`;
}

let duration = `0:00`;
let realDuration = 0;
onMount(() => {
})

let seekTime = `0:00`
let hoveredWidth = 0;
function track(e) {
    hoveredWidth = e.offsetX
    //console.log(e.target.getBoundingClientRect().width)
    const el = e.offsetX
    const w = seeker.getBoundingClientRect().width
    let p = (el/w) * realDuration
    seekTime = calculateTime(p)
}

function seekTo(e) {
    const el = e.offsetX
    const w = seeker.getBoundingClientRect().width
    let per = (el/w) * 100
    percent = per
    console.log(video.currentTime)
    video.currentTime = (el/w) * realDuration
    elapsed = calculateTime(video?.currentTime)
}

let seeker;

function dragSeek(e) {
    const el = e.offsetX
    const w = seeker.getBoundingClientRect().width
    let per = (el/w) * 100
    percent = per
    video.currentTime = (el/w) * realDuration
    elapsed = calculateTime(video?.currentTime)
}


$: thumbLeft = seeker?.getBoundingClientRect()?.width * (percent / 100) - 4

let progress;

let isMuted = false;

function toggleMute() {
    isMuted = !isMuted
    if(isMuted) {
        video.volume = 0
    } else {
        video.volume = 1
    }
}

$: formatSize = formatBytes(size)

let container;

let isFullscreen = false;

function toggleFullscreen() {
    isFullscreen = !isFullscreen
  if (isFullscreen) {
      container.requestFullscreen();
  } else {
      document.exitFullscreen();
  }
}

$: base = width < 400 ? width : 400

$: w = width > base ? base : width
$: h = (base/width) * height

</script>

<div class="video-item flex flex-column">

<div class="video-player relative" 
    style={`--width: ${w}px;--height: ${h}px;`}
    bind:this={container}
    class:active={active}>
    {#if !active}
        <div class="play-button" on:click={togglePlay}>
            {@html play}
        </div>
    {/if}
    <div class="video-title pa2 flex" class:op-0={active && !paused}>
        <div class="gr-center flex-one flex flex-column">
            <div class="clmp-1">
                {title}
            </div>
            <div class="size">
                {formatSize}
            </div>
        </div>
        <div class="gr-center">
        </div>
    </div>
    <div class="" on:click={togglePlay}>
        <video bind:this={video} src={url} preload=”metadata”></video>
    </div>
    <div class="video-controls" class:op-1={paused}>
        <div class="play gr-default ml1" on:click={togglePlay}>
            {#if (!playing && !played) || (!playing && percent > 0)}
                {@html play}
            {:else if playing}
                {@html pause}
            {:else if played}
                {@html refresh}
            {/if}
        </div>
        <div class="time gr-center ml1">
            <span class="">{elapsed}</span>
            <span class="sep">/</span>
            <span class="">{duration}</span>
        </div>
        <div class="seek-container w-100  ph2 gr-center">
            <div class="seek-wrapper pv2 relative" 
                bind:this={seeker}
                on:click={seekTo}
                on:dragstart={dragSeek}
                on:mousemove={track}>
                <div class="seek">
                    <div class="seek-inner"
                    style={`--hovered: ${hoveredWidth}px;`}>
                    </div>
                </div>
                <div class="seek-progress relative"
                    bind:this={progress}
                style={`--percent: ${percent}%;`}>
                </div>
                    <div class="seek-thumb"
                    style={`--left: ${thumbLeft}px;`}>
                    </div>
                <div class="seek-time gr-default"
                    style={`--left: ${hoveredWidth}px;`}>
                    <div class="gr-center">
                    {seekTime}
                    </div>
                </div>
            </div>
        </div>
        <div class="volume relative gr-default mr1">
            <div class="volume-icon gr-center" on:click={toggleMute}>
                {#if isMuted}
                    {@html muted}
                {:else}
                    {@html volume}
                {/if}
            </div>
        </div>
        <div class="fullscreen gr-center mr1" on:click={toggleFullscreen}>
            {@html fullscreen}
        </div>
    </div>
</div>
</div>


<style>
.video-player {
    width: var(--width);
    height: var(--height);
    background-color: var(--black);
}
.active:hover .video-controls{
    opacity: 1;
}

.active:hover .video-title{
    opacity: 1;
}

.video-title {
    position: absolute;
    top: 0;
    transition: 0.2s;
}

.video-controls {
    position: absolute;
    bottom: 0;
    height: 32px;
    width: 100%;
    background-color: #131416b0;
    border-radius: 2px;
    display: grid;
    grid-template-columns: auto auto 1fr auto auto;
    opacity: 0;
    transition: 0.2s;
}

.op-1 {
    opacity: 1;
}

.op-0 {
    opacity: 0;
}

.play {
    cursor: pointer;
    width: 24px;
    height: 100%;
    fill: var(--text);
}
.play:hover {
    fill: var(--white);
}

.play-button {
    cursor: pointer;
    width: 34px;
    fill: var(--white);
    position: absolute;
    left: calc(50% - 17px);
    top: calc(50% - 17px);
    box-shadow: 0 10px 20px rgba(0,0,0,.1);
}

.sep {
    margin: 0 2px;
}

.time {
    font-size: 0.8rem;
    color: var(--white);
    font-family: Consolas,Andale Mono WT,Andale Mono,Lucida Console,Lucida Sans Typewriter,DejaVu Sans Mono,Bitstream Vera Sans Mono,Liberation Mono,Nimbus Mono L,Monaco,Courier New,Courier,monospace;
}

.seek-wrapper {
    cursor: pointer;
}

.seek {
    width: 100%;
    height: 6px;
    border-radius: 50px;
    background-color: var(--background-3);
}

.seek-inner {
    width: var(--hovered);
    height: 100%;
    border-radius: 50px;
    background-color: var(--background-4);
    opacity: 0;
    transition: 0.01s;
}

.seek-wrapper:hover .seek-inner {
    opacity: 1;
}

.seek-thumb {
    position: absolute;
    left: var(--left);
    border-radius: 50px;
    background-color: var(--green);
    height: 6px;
    width: 6px;
    z-index: 100;
    top: 8px;
    transition: 0.1s;
}

.seek-wrapper:hover .seek-thumb {
    height: 10px;
    width: 10px;
    top: 6px;
}

.seek-progress {
    height: 6px;
    width: var(--percent);
    border-radius: 50px;
    background-color: var(--green);
    position: absolute;
    right: 0;
    left: 0;
    top: 8px;
    bottom: 0;
    transition: width 0.05s;
    z-index: 10;
}
.seek-time {
    position: absolute;
    left: calc(var(--left) - 30px);
    bottom: 24px;
    border-radius: 3px;
    background-color: var(--black);
    font-size: 0.8rem;
    padding: 0.25rem 0.5rem;
    opacity: 0;
    min-width: 50px;
}
.seek-wrapper:hover .seek-time {
    opacity: 1;
}

.volume-icon {
    cursor: pointer;
    width: 22px;
    height: 22px;
    fill: var(--text);
}
.volume-icon:hover {
    fill: var(--white);
}

.volume-seeker {
    position: absolute;
    bottom: 34px;
    left: 4px;
    width: 12px;
    height: 90px;
    background-color: var(--black);
    border-radius: 500px;
}

.volume-seeker input {
    margin: 0;
    -webkit-appearance: slider-vertical;
    width: 100%;
    height: 100%;
}
.size {
    font-size: 0.7rem;
    color: var(--text-muted);
}

.fullscreen {
    cursor: pointer;
    width: 22px;
    height: 22px;
    fill: var(--text);
}
.fullscreen:hover {
    fill: var(--white);
}

::-webkit-media-controls {
  display:none !important;
}

</style>
