<script>
import { store } from '../../../../store/store.js'
import { fly } from 'svelte/transition'
import { onMount, createEventDispatcher } from 'svelte'
import { mic, micOff, close } from '../../../../utils/icons.js'
import Modal from '../../../../components/modal/modal.svelte'

export let room;

const dispatch = createEventDispatcher()

export let page;


export let recording = false;

let started;

$: elapsed = new Date() - started

let permissions;

$: denied = permissions == "denied"
$: granted = permissions == "granted"

onMount(() => {
    navigator.permissions.query({name:'microphone'}).then(function(result) {
        permissions = result.state
    });
})

let showBlocked = false;

function record() {
    dispatch('toggleRecording', true)
}

$: if(recording) {
    init()
}

function init() {
    if(permissions == 'denied' || permissions == undefined) {
        showBlocked = true
        return
    }
    navigator.permissions.query({name:'microphone'}).then(function(result) {
        permissions = result.state
      if (result.state == 'granted') {
        startRecording()
      } else if (result.state == 'prompt') {
          navigator.mediaDevices.getUserMedia({ audio: true, video: false })
      } else if (result.state == 'denied') {
          showBlocked = true
      }
      result.onchange = function() {
        permissions = result.state
          if (result.state == 'denied') {
              showBlocked = true
          };
          if (result.state == 'granted') {
              startRecording()
          };
      };
    });
}

let downloadLink;
let mediaRecorder;

let blob;
let blobURL;

let stream;


function startRecording() {
    navigator.mediaDevices.getUserMedia({ audio: true, video: false })
        .then(stream => {
            console.log(stream)
            const options = {mimeType: 'audio/webm'};
            const recordedChunks = [];
            mediaRecorder = new MediaRecorder(stream, options);

            mediaRecorder.addEventListener('dataavailable', function(e) {
                if (e.data.size > 0) recordedChunks.push(e.data);
            });

            mediaRecorder.addEventListener('stop', function() {
                if(!recording) {
                    stream.getTracks().forEach( track => track.stop() );
                    return
                }
                blob = new Blob(recordedChunks, { 'type' : 'audio/webm; codecs=opus' });
                blobURL = URL.createObjectURL(blob)
                stream.getTracks().forEach( track => track.stop() );

                uploadRecording(blob).then(res => {
                    console.log(res)

                    if(res?.content_uri) {

                        let content = {
                            "body": "recording",
                            "msgtype": "m.recording",
                            "url": res?.content_uri,
                            "info" : {
                                "mimetype": "audio/webm",
                                "size": blob.size,
                            }
                        }

                        matrix.sendEvent(room.room_id, "m.room.message", content, "",(err, res) => {
                            if(res?.event_id) {
                                killRecord()
                            }
                        });
                    }
                })

            });

            mediaRecorder.start();
        });
    startTimer()
}

let timer;

function startTimer() {
    started = 0
    timer = setInterval(() => {
        started += 1
    }, 1000)
}

function yikes() {
    showBlocked = false
}

function killModal() {
    showBlocked = false
}

function killRecord() {
    dispatch('toggleRecording', true)
}

$: if(!recording) {
    stop()
}

function stop() {
    //mediaRecorder?.stop();
    stream = null
    mediaRecorder = null
    blob = null;
    clearInterval(timer)
    started = 0
    sending = false
}

let sending = false;

function send() {
    sending = true
    clearInterval(timer)
    mediaRecorder.stop();
}

$: account = $store.accounts.filter(account => account.user_id == $store.active_account)[0]

$: matrix = account?.matrix

async function uploadRecording(content) {
    //const key = createStorageKey(content)
    //const formData = createFormData(key, content)

    let endpoint = `${homeServer}/_matrix/media/r0/upload`

    let resp = await fetch(endpoint, {
    method: 'POST', // or 'PUT'
    body: content,
    headers:{
        'Authorization': `Bearer ${account.matrix_access_token}`,
        'Content-Type': "audio/webm"
    }
    })
    const ret = await resp.json()
    return Promise.resolve(ret)
}
function createStorageKey(file) {
    var date = new Date()
    var day = date.toISOString().slice(0,10)
    var name = date.getTime() + "-" + file.name
    return [ "tmp", day, name ].join("/")
}

function createFormData(key, file) {
    var data = new FormData()
    data.append("file", file)
    data.append("name", "test.webm")
    return data
}

$: time = formatTime(started)

function formatTime() {
    if(started < 10) {
        return `0:0${started}`
    }
    if(started < 60) {
        return `0:${started}`
    }
    if(started > 59) {
        return `1:${started - 60}`
    }
}
</script>


{#if !recording}
<div class="recording-container unset relative" >

        <div class="r-icon gr-center-start r-h stick"
        aria-label="Send a recording"
        data-microtip-position="top"
        data-microtip-size="fit"
        role="tooltip" 
            on:click={record}>
            {#if denied}
                {@html micOff}
            {:else}
                {@html mic}
            {/if}
        </div>

</div>

{:else}
    <div class="r-hol flex ph3" 
    in:fly="{{ y: -20, duration: 100 }}">
        <div class="gr-center r-icon red" >
            {@html mic}
        </div>
        <div class="gr-center timer ph2">
            {time}
        </div>
        <div class="gr-center">
            escape to <span class="link">cancel</span> • enter to <span class="link">send</span>
        </div>
        <div class="gr-center ph2">
            <button class="send" on:click={send}>{sending? 'Sending...' : 'Send'}</button>
        </div>
        <div class="gr-center r-icon" on:click={killRecord}>
            {@html close}
        </div>
    </div>
{/if}


{#if showBlocked}
    <Modal noStyle={true} on:killModal={killModal}>
        <div class="blocked ">
            <div class="gr-center flex flex-column pa3">
                <div class="gr-center tt">
                    microphone access is denied
                </div>
                <div class="gr-center pt3">
                    It looks like you've denied Commune access to your
                    microphone.
                </div>
            </div>
            <div class="da gr-center pa3">
                <button class="pv3" on:click={killModal}>Okay</button>
            </div>
         </div>
    </Modal>
{/if}

<style>
.recording-container {
}
.r-h{
    padding-top: 13px;
    padding-bottom: 12px;
}
.r-icon {
    height: 22px;
    width: 22px;
    cursor: pointer;
    fill: var(--text);
}
.r-icon:hover {
    fill: var(--white);
}
.red {
    fill: red;
    animation-name: blink;
    animation-duration: 0.5s;
    animation-iteration-count: infinite;
}
.red:hover {
    fill: red;
}
@keyframes blink {
  from {fill: red;}
  to {fill: var(--white);}
}
.red-alt {
    fill: red;
}

.no-dis {
  display: none;
}

.blocked {
    background-color: var(--background-3);
    border-radius: 7px;
    min-width: 440px;
    min-height: 210px;
    display: grid;
    grid-template-rows: 1fr auto;
}
.da {
    width: 100%;
    height: 100%;
    background-color: var(--background-2);
}
.send {
    width: 100%;
    border-radius: 500px;
    padding: 0.5rem 0.75rem;
}
.tt {
    font-weight: bold;
    text-transform: uppercase;
}
:root {
    --chat-input: #40444b;
}
.r-hol {
    background-color: var(--chat-input);
    border-radius: 8px;
    display: grid;
    grid-template-columns: auto auto 1fr auto auto;
    width: 100%;
}


.timer {
    width: 54px;
}
.mt {
    padding-top: 13px;
}
</style>
