<script>
import { onMount } from 'svelte'
import { audio as audioIcon, download, play, pause, refresh, volume, muted } from '../../utils/icons.js'
import { formatBytes } from '../../utils/utils.js'

export let url;
export let title;
export let size;

let audio;


let played = false;
let playing = false;

function togglePlay() {
    if(audio.paused) {
        audio.addEventListener('timeupdate', setTime);
        audio.addEventListener('play', () => {
            playing = true
        });
        audio.play()
    } else {
        audio.pause()
        playing = false
    }
}


let elapsed = `0:00`;
let percent = 0;

function setTime() {
    elapsed = calculateTime(audio?.currentTime)
    percent = (audio.currentTime / audio.duration) * 100
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

let getDuration = function (url, next) {
    var _player = new Audio(url);
    _player.addEventListener("durationchange", function (e) {
        if (this.duration!=Infinity) {
           var duration = this.duration
           _player.remove();
           next(duration);
        };
    }, false);      
    _player.load();
    _player.currentTime = 24*60*60; //fake big time
    _player.volume = 0;
    _player.play();
    //waiting...
};

const calculateTime = (secs) => {
  const minutes = Math.floor(secs / 60);
  const seconds = Math.floor(secs % 60);
  const returnedSeconds = seconds < 10 ? `0${seconds}` : `${seconds}`;
  return `${minutes}:${returnedSeconds}`;
}

let duration = `0:00`;
let realDuration;
onMount(() => {
    getDuration (url, function (d) {
        realDuration = d
        duration = calculateTime(d)
    });
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
    audio.currentTime = (el/w) * realDuration
    elapsed = calculateTime(audio?.currentTime)
}

let seeker;

function dragSeek(e) {
    const el = e.offsetX
    const w = seeker.getBoundingClientRect().width
    let per = (el/w) * 100
    percent = per
    audio.currentTime = (el/w) * realDuration
    elapsed = calculateTime(audio?.currentTime)
}


$: thumbLeft = seeker?.getBoundingClientRect()?.width * (percent / 100) - 4

let progress;

let isMuted = false;

function toggleMute() {
    isMuted = !isMuted
    if(isMuted) {
        audio.volume = 0
    } else {
        audio.volume = 1
    }
}

$: contentURL = `${homeServer}/_matrix/media/r0/download/${url?.substring(6)}`
$: formatSize = formatBytes(size)


function downloadFile() {
    fetch(url).then(res => {
        return res.blob().then(b =>{
            var a = document.createElement("a");
            a.href = URL.createObjectURL(b);
            a.setAttribute("download", title);
            a.click();
        });
    });
}

</script>


<div class="audio-item flex flex-column pa2" >
    <div class="flex">
        <div class="gr-default audio-icon">
            {@html audioIcon}
        </div>
        <div class="ml3 gr-center flex-one flex flex-column">
            <div class="bold clmp-2">
                <a href={contentURL}>{title}</a>
            </div>
            <div class="wh-s">
                {formatSize}
            </div>
        </div>
        <div class="download-icon pointer gr-default flex" on:click={downloadFile}>
            {@html download}
        </div>
    </div>
    <div class="audio-container flex pt2">



        <div class="audio-player">
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
        </div>
        <audio bind:this={audio} src={url} preload=”metadata” hidden></audio>

    </div>
</div>

<style>
.audio-item {
    min-height: 98px;
    width: 400px;
    background-color:var(--background-2);
    border: 1px solid var(--background-5);
    border-radius: 4px;
}
.audio-player {
    height: 32px;
    width: 100%;
    background-color: var(--dark-1);
    border-radius: 2px;
    display: grid;
    grid-template-columns: auto auto 1fr auto;
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

audio {
    display: none;
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
.audio-icon {
    height: 40px;
    width: 24px;
}

.audio-item {
    min-height: 98px;
    width: 400px;
    background-color:var(--background-2);
    border: 1px solid var(--background-5);
    border-radius: 4px;
}
.download-icon {
    height: 24px;
    width: 24px;
    fill: var(--text);
}
.download-icon:hover {
    fill: var(--white);
}
</style>
