<script>
import { useLocation } from 'svelte-navigator'

const location = useLocation()

$: showMe = $location.pathname.substring(1,2) == `@`

$: username = $location.pathname.slice(1)

$: if(showMe) {
    document.title = username
}

</script>

<div class="home" class:no-dis={!showMe}>
    <div class="user-profile w-100 scrl">
        User Profile
    </div>
</div>

<style>
.home {
    height: 100%;
    overflow: hidden;
    display: grid;
}
.user-profile {
    background-color: var(--background-3);
    overflow: auto;
}

.no-dis {
     display: none;
}
</style>
