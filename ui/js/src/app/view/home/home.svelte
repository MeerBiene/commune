<script>
import { useLocation } from 'svelte-navigator'

let text = `<div class="pb5">lol</div>`.repeat(20)

const location = useLocation()

$: showHome = $location.pathname == '/'

function newDM() {
    window.addNewDM()
}

function newServer() {
    window.addNewServer()
}

</script>

<div class="home" class:no-dis={!showHome}>
    <div class="chat w-100 scrl gr-default">
        <div class="gr-center flex flex-column">
            <div class=" center">
                This page isn't very useful right now. It ought to be a central
                feed.
            </div>
            <div class="pt3 center">
                <span class="it" on:click={newServer}>
                    Add/Join a Server
                </span>
                or
                <span class="it" on:click={newDM}>
                    start a DM
                </span>
            </div>
        </div>
    </div>
</div>

<style>
.home {
    height: 100%;
    overflow: hidden;
    display: grid;
    background-color: var(--background-3);
}
.chat {
    overflow: auto;
}

.no-dis {
     display: none;
}
.it {
    color: var(--text-link);
    cursor: pointer;
}
.it:hover {
    text-decoration: underline;
}
</style>
