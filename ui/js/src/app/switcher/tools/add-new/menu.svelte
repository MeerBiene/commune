<script>
import {createEventDispatcher} from 'svelte'
const dispatch = createEventDispatcher()

function kill() {
    dispatch('kill', true)
}
function create() {
    dispatch('create', true)
}
function join() {
    dispatch('join', true)
}
</script>

<div class="flex flex-column pa3">
    <div class="flex pb3">
        <div class="flex-one">
        </div>
        <div class="icon pointer" on:click={kill}>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M5.72 5.72a.75.75 0 011.06 0L12 10.94l5.22-5.22a.75.75 0 111.06 1.06L13.06 12l5.22 5.22a.75.75 0 11-1.06 1.06L12 13.06l-5.22 5.22a.75.75 0 01-1.06-1.06L10.94 12 5.72 6.78a.75.75 0 010-1.06z"></path></svg>
        </div>
    </div>
    <div class="menu-item flex" on:click={create}>
        <div class="gr-center">
        </div>
        <div class="gr-center flex-one flex flex-column">
            <div class="f4 bold">
                Create Your Own Server
            </div>
            <div class="lh-copy pt2">
                Your server is where you and your friends hang out. Create yours
                and start talking.
            </div>
        </div>
        <div class="gr-center">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M8.72 18.78a.75.75 0 001.06 0l6.25-6.25a.75.75 0 000-1.06L9.78 5.22a.75.75 0 00-1.06 1.06L14.44 12l-5.72 5.72a.75.75 0 000 1.06z"></path></svg>
        </div>
    </div>
    <div class="menu-item mt3 flex" on:click={join}>
        <div class="gr-center">
        </div>
        <div class="gr-center flex-one flex flex-column">
            <div class="f4 bold">
                Join An Existing Server
            </div>
            <div class="lh-copy pt2">
                If you have an invite link, join an existing server.
            </div>
        </div>
        <div class="gr-center">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24"><path fill-rule="evenodd" d="M8.72 18.78a.75.75 0 001.06 0l6.25-6.25a.75.75 0 000-1.06L9.78 5.22a.75.75 0 00-1.06 1.06L14.44 12l-5.72 5.72a.75.75 0 000 1.06z"></path></svg>
        </div>
    </div>
</div>

<style>
.menu-item {
    padding: 0.5rem 1rem;
    border: 1px solid var(--background-2);
    border-radius: 9px;
    cursor: pointer;
    fill: var(--text-muted);
}

.menu-item:hover {
    background-color: var(--background-2);
    fill: var(--white);
}
</style>
